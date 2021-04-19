package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

var tests = []struct {
	in          interface{}
	expectedErr error
}{
	{
		in: User{
			ID:     "1234567890", // wrong len
			Name:   "username",
			Age:    18,
			Email:  "user@email.ru",
			Role:   "admin",
			Phones: []string{"12345678901"},
			meta:   []byte("bytes"),
		},
		expectedErr: ErrWrongFieldLen,
	},
	{
		in: User{
			ID:     "123456789012345678901234567890123456",
			Name:   "username",
			Age:    17, // 17 less than 18
			Email:  "user@email.ru",
			Role:   "admin",
			Phones: []string{"12345678901"},
			meta:   []byte("bytes"),
		},
		expectedErr: ErrValueExceededMinimum,
	},
	{
		in: User{
			ID:     "123456789012345678901234567890123456",
			Name:   "username",
			Age:    51, // 51 greater than 50
			Email:  "user@email.ru",
			Role:   "admin",
			Phones: []string{"12345678901"},
			meta:   []byte("bytes"),
		},
		expectedErr: ErrValueExceededMaximum,
	},
	{
		in: User{
			ID:     "123456789012345678901234567890123456",
			Name:   "username",
			Age:    18,
			Email:  "user@email.ru!", // irregular email address
			Role:   "admin",
			Phones: []string{"12345678901"},
			meta:   []byte("bytes"),
		},
		expectedErr: ErrRegexpNotMatch,
	},
	{
		in: User{
			ID:     "123456789012345678901234567890123456",
			Name:   "username",
			Age:    18,
			Email:  "user@email.ru",
			Role:   "wrong_role", // wrong role
			Phones: []string{"12345678901"},
			meta:   []byte("bytes"),
		},
		expectedErr: ErrValueNotExists,
	},
	{
		in: User{
			ID:     "123456789012345678901234567890123456",
			Name:   "username",
			Age:    18,
			Email:  "user@email.ru",
			Role:   "admin",
			Phones: []string{"123456789012"}, // wrong phone length
			meta:   []byte("bytes"),
		},
		expectedErr: ErrWrongFieldLen,
	},
	{
		in: App{
			Version: "123456", // wrong string length
		},
		expectedErr: ErrWrongFieldLen,
	},
	{
		in: Response{
			Code: 400, // wrong response
			Body: "string",
		},
		expectedErr: ErrValueNotExists,
	},
	{
		in:          "i'm not a struct", // not a struct
		expectedErr: ErrNotAStruct,
	},
	{
		in: struct {
			ID int `validate:"len:36"`
		}{
			ID: 123, // wrong length
		},
		expectedErr: ErrWrongFieldType,
	},
	{
		in: struct {
			ID int `validate:"in:1,2,nan"` // wrong validator value
		}{
			ID: 1,
		},
		expectedErr: ErrWrongValidatorValue,
	},
	{
		in: struct {
			ID int `validate:"not_implemented_validator:value"` // not implemented calidator
		}{
			ID: 1,
		},
		expectedErr: ErrNotImplementedValidator,
	},
	{
		in: struct {
			ID int `validate:"min:1"`
		}{
			ID: 0, // less than 1
		},
		expectedErr: ErrValueExceededMinimum,
	},
	{
		in: struct {
			ID int `validate:"max:10"`
		}{
			ID: 11, // greater than 11
		},
		expectedErr: ErrValueExceededMaximum,
	},
	{
		in: struct {
			Email string `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		}{
			Email: "wrong@email.ru!", // wrong email
		},
		expectedErr: ErrRegexpNotMatch,
	},
	{
		in: struct {
			Role string `validate:"in:admin,stuff"`
		}{
			Role: "user", // not existence role
		},
		expectedErr: ErrValueNotExists,
	},
	{
		in: struct {
			ID int `validate:"in:1,2"`
		}{
			ID: 3, // not existence value
		},
		expectedErr: ErrValueNotExists,
	},
	{
		in: struct {
			Name string `validate:"len:5"`
		}{
			Name: "123", // wrong length
		},
		expectedErr: ErrWrongFieldLen,
	},
}

func TestValidate(t *testing.T) {
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt

			err := Validate(tt.in)

			va := ValidationErrors{}
			if errors.As(err, &va) {
				if len(va) == 0 {
					require.Truef(t, false, "expected error: %q, nothing given", tt.expectedErr)
				} else {
					for _, e := range va {
						require.Truef(t, errors.Is(e.Err, tt.expectedErr), "actual error %q", e.Err)
						break
					}
				}
			} else {
				require.Truef(t, errors.Is(err, tt.expectedErr), "actual error %q", err)
			}
		})
	}
}
