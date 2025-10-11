package editor

import (
	"errors"

	"github.com/charmbracelet/bubbles/textinput"
)

var ErrRequired = errors.New("Required")

type FieldsMod func(fields []textinput.Model)

func ModifyField(i int, mod func(f textinput.Model) textinput.Model) FieldsMod {
	return func(fields []textinput.Model) {
		fields[i] = mod(fields[i])
	}
}

func RequireFields(is ...int) FieldsMod {
	return func(fields []textinput.Model) {
		for i := range is {
			AddFieldValidator(i, func(s string) error {
				if s == "" {
					return ErrRequired
				}

				return nil
			})(fields)
		}
	}
}

func AddFieldValidator(i int, validate func(s string) error) FieldsMod {
	return func(fields []textinput.Model) {
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
