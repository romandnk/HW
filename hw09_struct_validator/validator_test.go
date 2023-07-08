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
		ID     string `json:"id" validate:"len:4"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:nolintlint
	}

	App struct {
		Version string `validate:"len:wrong"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	RegExp struct {
		Smth string `validate:"regexp:[a-z"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "uuid",
				Name:   "TestOneName",
				Age:    19,
				Email:  "test_email@yandex.com",
				Role:   "stuff",
				Phones: []string{"88005553535"},
				meta:   []byte{},
			},
			expectedErr: ValidationErrors{},
		},
		{
			in: User{
				ID: "uuidd",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   ErrInvalidLen,
				},
			},
		},
		{
			in: User{
				ID:  "uuid",
				Age: 14,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Age",
					Err:   ErrTooSmallValue,
				},
			},
		},
		{
			in: User{
				ID:  "uuid",
				Age: 51,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Age",
					Err:   ErrTooBigValue,
				},
			},
		},
		{
			in: User{
				ID:     "uuid",
				Age:    20,
				Email:  "test@test.com",
				Role:   "stuff",
				Phones: []string{"123456789"},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Phones",
					Err:   ErrInvalidLen,
				},
			},
		},
		{
			in: User{
				ID:    "uuid",
				Age:   40,
				Email: "sldkf",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Email",
					Err:   ErrStrNotMatch,
				},
			},
		},
		{
			in:          App{},
			expectedErr: ErrWrongLenCond,
		},
		{
			in:          Token{},
			expectedErr: ValidationErrors{},
		},
		{
			in: Response{
				Code: 300,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Code",
					Err:   ErrNumberNotInSet,
				},
			},
		},
		{
			in: RegExp{
				Smth: "sldfk",
			},
			expectedErr: ErrCompileRegExp,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			var validationErrors, expectedErrors ValidationErrors

			if errors.As(err, &validationErrors) && errors.As(tt.expectedErr, &expectedErrors) {
				for j := range expectedErrors {
					require.Equal(t, validationErrors[j], expectedErrors[j])
				}
			} else {
				require.ErrorIs(t, err, tt.expectedErr)
			}
		})
	}
}
