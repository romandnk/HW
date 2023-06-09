package hw09structvalidator

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	ErrNotStruct                   = errors.New("input data is not a struct")
	ErrInvalidValidationParameters = errors.New("invalid validation parameters")
	ErrTooSmallValue               = errors.New("value must be bigger")
	ErrTooBigValue                 = errors.New("value must be smaller")
	ErrNumberNotInSet              = errors.New("number is not in specified set")
	ErrMultipleCondition           = errors.New("multiple conditions")
	ErrInvalidLen                  = errors.New("len of string does not match with condition")
	ErrCompileRegExp               = errors.New("error compiling regular expression")
	ErrStrNotMatch                 = errors.New("string does not match with regular expression")
	ErrWrongType                   = errors.New("can validate only int, string, []int, []string")
	ErrWrongInCond                 = errors.New("in must be array with numbers")
	ErrWrongMinCond                = errors.New("min must be a number")
	ErrWrongLenCond                = errors.New("len must be a number")
	ErrWrongMaxCond                = errors.New("max must be a number")
)

type ValidationError struct {
	Field string
	Err   error
}

func (v ValidationError) Error() string {
	return v.Err.Error()
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var resultString strings.Builder
	for _, err := range v {
		resultString.WriteString(err.Field + ": " + err.Err.Error())
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
	var myError ValidationError

	for i := 0; i < valT.NumField(); i++ {
		field := valT.Field(i)
		valField := val.Field(i)
		nameField := field.Name

		// check if the field of the struct is public
		if field.PkgPath != "" {
			continue
		}

		validateField, ok := field.Tag.Lookup("validate")
		if !ok {
			continue
		}

		switch field.Type.Kind() { //nolint:exhaustive
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			err := validateInt(valField.Int(), validateField, nameField)
			if err != nil {
				if errors.As(err, &myError) {
					resultErrors = append(resultErrors, myError)
				} else {
					return err
				}
			}
		case reflect.String:
			err := validateString(valField.String(), validateField, nameField)
			if err != nil {
				if errors.As(err, &myError) {
					resultErrors = append(resultErrors, myError)
				} else {
					return err
				}
			}
		case reflect.Slice:
			switch field.Type.Elem().Kind() { //nolint:exhaustive
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				for i := 0; i < valField.Len(); i++ {
					elem := valField.Index(i).Int()
					err := validateInt(elem, validateField, nameField)
					if err != nil {
						if errors.As(err, &myError) {
							resultErrors = append(resultErrors, myError)
						} else {
							return err
						}
					}
				}
			case reflect.String:
				for i := 0; i < valField.Len(); i++ {
					elem := valField.Index(i).String()
					err := validateString(elem, validateField, nameField)
					if err != nil {
						if errors.As(err, &myError) {
							resultErrors = append(resultErrors, myError)
						} else {
							return err
						}
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
func validateInt(fieldVal int64, validateField string, nameField string) error {
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
			return ValidationError{
				Field: nameField,
				Err:   ErrInvalidValidationParameters,
			}
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
						return ErrWrongInCond
					}
					if int64(valInt) == fieldVal {
						isFound = true
						break
					}
				}

				if !isFound {
					return ValidationError{
						Field: nameField,
						Err:   ErrNumberNotInSet,
					}
				}
			}
		case "min":
			if ok := possibleConditions["min"]; !ok {
				if possibleConditions["in"] {
					return ValidationError{
						Field: nameField,
						Err:   ErrMultipleCondition,
					}
				}

				possibleConditions["min"] = true

				value, err := strconv.Atoi(condition)
				if err != nil {
					return ErrWrongMinCond
				}

				if fieldVal < int64(value) {
					return ValidationError{
						Field: nameField,
						Err:   ErrTooSmallValue,
					}
				}
			}
		case "max":
			if ok := possibleConditions["max"]; !ok {
				if possibleConditions["in"] {
					return ValidationError{
						Field: nameField,
						Err:   ErrMultipleCondition,
					}
				}

				possibleConditions["max"] = true

				value, err := strconv.Atoi(condition)
				if err != nil {
					return ErrWrongMaxCond
				}

				if fieldVal > int64(value) {
					return ValidationError{
						Field: nameField,
						Err:   ErrTooBigValue,
					}
				}
			}
		}
	}
	return nil
}

//nolint:gocognit
func validateString(fieldVal string, validateField string, nameField string) error {
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
			return ValidationError{
				Field: nameField,
				Err:   ErrInvalidValidationParameters,
			}
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
					return ValidationError{
						Field: nameField,
						Err:   ErrNumberNotInSet,
					}
				}
			}
		case "len":
			if ok := possibleConditions["len"]; !ok {
				if possibleConditions["in"] {
					return ValidationError{
						Field: nameField,
						Err:   ErrMultipleCondition,
					}
				}
				possibleConditions["len"] = true

				mustLen, err := strconv.Atoi(condition)
				if err != nil {
					return ErrWrongLenCond
				}

				if utf8.RuneCountInString(fieldVal) != mustLen {
					return ValidationError{
						Field: nameField,
						Err:   ErrInvalidLen,
					}
				}
			}
		case "regexp":
			if ok := possibleConditions["regexp"]; !ok {
				if possibleConditions["in"] {
					return ValidationError{
						Field: nameField,
						Err:   ErrMultipleCondition,
					}
				}

				possibleConditions["regexp"] = true

				match, err := regexp.MatchString(condition, fieldVal)
				if err != nil {
					return ErrCompileRegExp
				}

				if !match {
					return ValidationError{
						Field: nameField,
						Err:   ErrStrNotMatch,
					}
				}
			}
		}
	}
	return nil
}
