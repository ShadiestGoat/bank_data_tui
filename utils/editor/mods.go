package editor

import (
	"fmt"
	"log"
	"slices"
	"strconv"

	"github.com/bank_data_tui/utils"
	"github.com/charmbracelet/bubbles/textinput"
)

type ErrRequired struct {
	FieldType string
}
func (e ErrRequired) Error() string {
	if e.FieldType == "" {
		return "Required"
	}
	return "A " + e.FieldType + " is required"
}

func (e ErrRequired) Is(err error) bool {
	_, ok := err.(*ErrRequired)
	if ok {
		return true
	}
	_, ok = err.(ErrRequired)
	return ok
}

type APIErr string

func (a APIErr) Error() string {
	return string(a)
}
func (a APIErr) Is(err error) bool {
	_, ok := err.(APIErr)
	return ok
}

type FieldsMod func(fields []*textinput.Model)

func ModifyField(i int, mod func(f *textinput.Model) *textinput.Model) FieldsMod {
	return func(fields []*textinput.Model) {
		fields[i] = mod(fields[i])
	}
}

func RequireFields(is ...int) FieldsMod {
	return func(fields []*textinput.Model) {
		for i := range is {
			AddFieldValidator(i, func(s string) error {
				if s == "" {
					return ErrRequired{}
				}

				return nil
			})(fields)
		}
	}
}

func AddFieldValidator(i int, validate func(s string) error) FieldsMod {
	return func(fields []*textinput.Model) {
		og := fields[i].Validate
		fields[i].Validate = func(s string) error {
			if og != nil {
				if err := og(s); err != nil {
					return err
				}
			}

			return validate(s)
		}
	}
}

func fieldValues(is []int, fields []*textinput.Model) []string {
	res := make([]string, len(is))
	for i, fi := range is {
		res[i] = fields[fi].Value()
	}

	log.Println(res)

	return res
}

func AddOneOfRequirement(fieldType string, is ...int) FieldsMod {
	return AddMultiFieldValidator(func(s []string) error {
		if utils.All(slices.Values(s), func(v string) bool {return v == ""}) {
			return ErrRequired{FieldType: fieldType}
		}

		return nil
	}, is...)
}

func AddMultiFieldValidator(validate func (s []string) error, is ...int) FieldsMod {
	return func(fields []*textinput.Model) {
		for _, i := range is {
			og := fields[i].Validate
			fields[i].Validate = func(s string) error {
				if og != nil {
					if err := og(s); err != nil {
						return err
					}
				}

				return validate(fieldValues(is, fields))
			}
		}
	}
}

func AddIntValidator(is ...int) FieldsMod {
	return func(fields []*textinput.Model) {
		for _, i := range is {
			AddFieldValidator(i, func(s string) error {
				if s != "" {
					_, err := strconv.Atoi(s)
					if err != nil {
						return fmt.Errorf("Must be int!")
					}
				}

				return nil
			})(fields)
		}
	}
}

func AddFloatValidator(is ...int) FieldsMod {
	return func(fields []*textinput.Model) {
		for _, i := range is {
			AddFieldValidator(i, func(s string) error {
				if s != "" {
					_, err := strconv.ParseFloat(s, 64)
					if err != nil {
						return fmt.Errorf("Must decimal!")
					}
				}

				return nil
			})(fields)
		}
	}
}
