package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	ErrNotStruct                   = errors.New("input data is not a struct")
	ErrInvalidValidationParameters = errors.New("invalid validation parameters")
	ErrValueNotNumber              = errors.New("the parameter value must be a number")
	ErrTooSmallValue               = errors.New("value must be bigger")
	ErrTooBigValue                 = errors.New("value must be smaller")
	ErrNumberNotInSet              = errors.New("number is not in specified set")
	ErrMultipleCondition           = errors.New("multiple conditions")
	ErrInvalidLen                  = errors.New("len of string does not match with condition")
	ErrCompileRegExp               = errors.New("error compiling regular expression")
	ErrStrNotMatch                 = errors.New("string does not match with regular expression")
	ErrWrongType                   = errors.New("can validate only int, string, []int, []string")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var resultString strings.Builder
	for _, err := range v {
		resultString.WriteString(fmt.Sprintf("%s - %s", err.Field, err.Err.Error()))
	}
	return resultString.String()
}

//nolint:gocognit
func Validate(v interface{}) error {
	val := reflect.ValueOf(v)

	// if kind of v is pointer then dereference the v
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// check if v is a struct
	if val.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	valT := val.Type()

	resultErrors := make(ValidationErrors, 0, valT.NumField())

	for i := 0; i < valT.NumField(); i++ {
		field := valT.Field(i)
		valField := val.Field(i)

		// check if the field of the struct is public
		if field.PkgPath != "" {
			continue
		}

		validateField, ok := field.Tag.Lookup("validate")
		if !ok {
			continue
		}

		switch field.Type.Kind() { //nolint:exhaustive
		case reflect.Int:
			err := validateInt(valField.Int(), validateField)
			resultErrors, err = checkError(resultErrors, field.Name, err)
			if err != nil {
				return err
			}
		case reflect.String:
			err := validateString(valField.String(), validateField)
			resultErrors, err = checkError(resultErrors, field.Name, err)
			if err != nil {
				return err
			}
		case reflect.Slice:
			switch field.Type.Elem().Kind() { //nolint:exhaustive
			case reflect.Int:
				for i := 0; i < valField.Len(); i++ {
					elem := valField.Index(i).Int()
					err := validateInt(elem, validateField)
					resultErrors, err = checkError(resultErrors, field.Name, err)
					if err != nil {
						return err
					}
				}
			case reflect.String:
				for i := 0; i < valField.Len(); i++ {
					elem := valField.Index(i).String()
					err := validateString(elem, validateField)
					resultErrors, err = checkError(resultErrors, field.Name, err)
					if err != nil {
						return err
					}
				}
			default:
				return ErrWrongType
			}
		default:
			return ErrWrongType
		}
	}

	return resultErrors
}

//nolint:gocognit
func validateInt(fieldVal int64, validateField string) error {
	// a map for checking the first condition encountered
	possibleConditions := map[string]bool{
		"min": false,
		"max": false,
		"in":  false,
	}

	validateOptions := strings.Split(validateField, "|")

	for _, opt := range validateOptions {
		parts := strings.Split(opt, ":")
		if len(parts) != 2 {
			return ErrInvalidValidationParameters
		}

		parameter := parts[0]
		condition := parts[1]

		switch parameter {
		// in has the highest priority
		case "in":
			if ok := possibleConditions["in"]; !ok {
				possibleConditions["in"] = true

				isFound := false

				numberSet := strings.Split(condition, ",")
				for _, valStr := range numberSet {
					valInt, err := strconv.Atoi(valStr)
					if err != nil {
						return fmt.Errorf("in: %w", ErrValueNotNumber)
					}
					if int64(valInt) == fieldVal {
						isFound = true
						break
					}
				}

				if !isFound {
					return ErrNumberNotInSet
				}
			}
		case "min":
			if ok := possibleConditions["min"]; !ok {
				if possibleConditions["in"] {
					return fmt.Errorf("%w: min with in", ErrMultipleCondition)
				}

				possibleConditions["min"] = true

				value, err := strconv.Atoi(condition)
				if err != nil {
					return fmt.Errorf("min: %w", ErrValueNotNumber)
				}

				if fieldVal < int64(value) {
					return fmt.Errorf("%w, value: %d", ErrTooSmallValue, fieldVal)
				}
			}
		case "max":
			if ok := possibleConditions["max"]; !ok {
				if possibleConditions["in"] {
					return fmt.Errorf("%w: max with in", ErrMultipleCondition)
				}

				possibleConditions["max"] = true

				value, err := strconv.Atoi(condition)
				if err != nil {
					return fmt.Errorf("max: %w", ErrValueNotNumber)
				}

				if fieldVal > int64(value) {
					return fmt.Errorf("%w, value: %d", ErrTooBigValue, fieldVal)
				}
			}
		}
	}
	return nil
}

//nolint:gocognit
func validateString(fieldVal string, validateField string) error {
	// a map for checking the first condition encountered
	possibleConditions := map[string]bool{
		"len":    false,
		"regexp": false,
		"in":     false,
	}

	validateOptions := strings.Split(validateField, "|")

	for _, opt := range validateOptions {
		parts := strings.Split(opt, ":")
		if len(parts) != 2 {
			return ErrInvalidValidationParameters
		}

		parameter := parts[0]
		condition := parts[1]

		switch parameter {
		// in has the highest priority
		case "in":
			if ok := possibleConditions["in"]; !ok {
				possibleConditions["in"] = true

				isFound := false

				numberSet := strings.Split(condition, ",")
				for _, val := range numberSet {
					if fieldVal == val {
						isFound = true
						break
					}
				}

				if !isFound {
					return ErrNumberNotInSet
				}
			}
		case "len":
			if ok := possibleConditions["len"]; !ok {
				if possibleConditions["in"] {
					return fmt.Errorf("%w: len with in", ErrMultipleCondition)
				}
				possibleConditions["len"] = true

				mustLen, err := strconv.Atoi(condition)
				if err != nil {
					return fmt.Errorf("len: %w", ErrValueNotNumber)
				}

				if utf8.RuneCountInString(fieldVal) != mustLen {
					return fmt.Errorf("%w, value: %s", ErrInvalidLen, fieldVal)
				}
			}
		case "regexp":
			if ok := possibleConditions["regexp"]; !ok {
				if possibleConditions["in"] {
					return fmt.Errorf("%w: regexp with in", ErrMultipleCondition)
				}

				possibleConditions["regexp"] = true

				match, err := regexp.MatchString(condition, fieldVal)
				if err != nil {
					return ErrCompileRegExp
				}

				if !match {
					return fmt.Errorf("%w, value: %s", ErrStrNotMatch, fieldVal)
				}
			}
		}
	}
	return nil
}

func checkError(outResult ValidationErrors, fieldName string, err error) (ValidationErrors, error) {
	systemErrors := map[error]struct{}{
		ErrCompileRegExp:                         {},
		fmt.Errorf("len: %w", ErrValueNotNumber): {},
		ErrInvalidValidationParameters:           {},
		fmt.Errorf("in: %w", ErrValueNotNumber):  {},
		fmt.Errorf("min: %w", ErrValueNotNumber): {},
		fmt.Errorf("max: %w", ErrValueNotNumber): {},
	}

	if _, ok := systemErrors[err]; ok {
		return outResult, err
	}

	outResult = append(outResult, ValidationError{
		Field: fieldName,
		Err:   err,
	})

	return outResult, nil
}
