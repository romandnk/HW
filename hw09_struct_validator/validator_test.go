package hw09structvalidator

import (
	"encoding/json"
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
		Version string `validate:"len:5"`
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
			expectedErr: nil,
		},
		//{
		//	in: User{
		//		ID:     "uuidd",
		//		Name:   "TestOneName",
		//		Age:    19,
		//		Email:  "test_email@yandex.com",
		//		Role:   "stuff",
		//		Phones: []string{"88005553535"},
		//		meta:   []byte{},
		//	},
		//	expectedErr: ValidationErrors{
		//		ValidationError{
		//			Field: "ID",
		//			Err:   fmt.Errorf("len of string does not match with condition, value: uuidd"),
		//		},
		//	},
		// },
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			var validationErrors ValidationErrors

			require.ErrorAs(t, err, &validationErrors)

			for _, e := range validationErrors {
				require.Equal(t, tt.expectedErr, e.Err.Error())
			}
		})
	}
}
