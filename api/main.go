package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const API_BASE_URL = "http://localhost:3000"

type ReqLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RespLogin struct {
	Token string `json:"token"`
}

type APIErr struct {
	Status   int
	RespBody []byte
}

func (a APIErr) Error() string {
	return fmt.Sprintf("%d: %s", a.Status, a.RespBody)
}

type StdAPIError struct {
	Status int
	Message string `json:"error"`
	Details []string `json:"details"`
}

func (s StdAPIError) Error() string {
	return fmt.Sprintf("%d: %s", s.Status, s.Message)
}

func (v StdAPIError) Is(err error) bool {
	_, ok := err.(*APIErr)
	return ok
}

type ValidationErr struct {
	Status int
	// List of [field id, error]
	Details [][2]string
}

func (v ValidationErr) Is(err error) bool {
	_, ok := err.(*APIErr)
	return ok
}

func (ValidationErr) Error() string {
	return "Resource validation didn't go over well :("
}

func fetch[T any](method, path string, body any, authHeader string) (*T, error) {
	inpBuf := &bytes.Buffer{}
	err := json.NewEncoder(inpBuf).Encode(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, API_BASE_URL+path, inpBuf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		d, _ := io.ReadAll(resp.Body)
		std := &StdAPIError{Status: resp.StatusCode}
		if err := json.Unmarshal(d, std); err != nil {
			return nil, &APIErr{resp.StatusCode, d}
		}

		if resp.StatusCode == 400 && len(std.Details) != 0 {
			verr := &ValidationErr{
				Details: make([][2]string, len(std.Details)),
			}
			good := true

			for i, v := range std.Details {
				spl := strings.SplitN(v, ": ", 2)
				if len(spl) != 2 {
					good = false
					break
				}

				verr.Details[i] = [2]string{spl[0], spl[1]}
			}

			if good {
				return nil, verr
			}
		}

		return nil, std
	}

	var data T
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func loginReq(user string, pass string) (string, error) {
	l, err := fetch[RespLogin](`POST`, `/login`, &ReqLogin{user, pass}, "")
	if err != nil {
		return "", err
	}

	return l.Token, nil
}

func loginIntoClient(c *APIClient, ovr [2]string) error {
	if ovr[0] == "" {
		ovr = c.userPass
	}

	tok, err := loginReq(ovr[0], ovr[1])
	if err != nil {
		return err
	}
	parsed, _, err := jwt.NewParser().ParseUnverified(tok, jwt.MapClaims{})
	if err != nil {
		return err
	}

	c.userPass = ovr
	c.jwt = parsed

	return nil
}

func easyFetch[T any](c *APIClient, method, path string, body any) (*T, error) {
	if d, err := c.jwt.Claims.GetExpirationTime(); err != nil || d.Before(time.Now()) {
		if err := loginIntoClient(c, [2]string{}); err != nil {
			return nil, err
		}
	}

	t, err := fetch[T](method, path, body, c.jwt.Raw)
	if err != nil {
		if e, ok := err.(*APIErr); ok && e.Status == 401 {
			if err := loginIntoClient(c, [2]string{}); err != nil {
				return nil, err
			}

			t, err = fetch[T](method, path, body, c.jwt.Raw)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return t, nil
}

func deArray[T any](v *[]T, err error) ([]T, error) {
	if v != nil {
		return *v, err
	}
	return nil, err
}

type APIClient struct {
	userPass [2]string
	jwt      *jwt.Token
}

func (a *APIClient) Login(userAndPass [2]string) error {
	return loginIntoClient(a, userAndPass)
}
