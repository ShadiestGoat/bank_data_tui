package api

import (
	"net/url"
	"strconv"
	"time"
)

type Transaction struct {
	SettledAt          time.Time `json:"settledAt"`
	AuthedAt           time.Time `json:"authedAt"`
	Desc               string    `json:"description"`
	Amount             float64   `json:"amount"`
	ResolvedName       *string   `json:"resolvedName,omitempty"`
	ResolvedCategoryID *string   `json:"resolvedCategoryId,omitempty"`
}

type TransactionFields string

const (
	TOR_AUTH     TransactionFields = "authed_at"
	TOR_SETTLE   TransactionFields = "settled_at"
	TOR_AMOUNT   TransactionFields = "amount"
	TOR_CATEGORY TransactionFields = "category"
)

func (c *APIClient) TransactionsFetch(orderBy TransactionFields, page int, asc bool) (*RespPages[[]*Transaction], error) {
	q := url.Values{}
	q.Set("page", strconv.Itoa(page))
	q.Set("order", string(orderBy))
	if asc {
		q.Set("asc", "true")
	} else {
		q.Set("asc", "false")
	}

	return easyFetch[RespPages[[]*Transaction]](c, `GET`, `/transactions?`+q.Encode(), nil)
}
