package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
)

// Validator name.
const ValidatorRegexp = "regexp"

var (
	ErrRegexpCompilation = errors.New("regexp compilation error")
	ErrRegexpNotMatch    = errors.New("regexp matching error")
)

// Validates the field by the "regexp" validator.
func ValidateRegexp(fieldName string, value reflect.Value, validatorValue string) {
	// Recursive call if the field is a slice of values.
	if value.Kind() == reflect.Slice {
		for j := 0; j < value.Len(); j++ {
			ValidateRegexp(fieldName, value.Index(j), validatorValue)
		}
		return
	}

	if value.Kind() != reflect.String {
		err := fmt.Errorf("%w: %s validator supports string type, %s given", ErrWrongFieldType, ValidatorRegexp, value.Kind().String())
		validationErrors = append(validationErrors, ValidationError{Field: fieldName, Err: err})
		return
	}

	re, err := regexp.Compile(validatorValue)
	if err != nil {
		err := fmt.Errorf("%w: pattern %s can't be compiled", ErrRegexpCompilation, validatorValue)
		validationErrors = append(validationErrors, ValidationError{Field: fieldName, Err: err})
		return
	}

	if !re.MatchString(value.String()) {
		err := fmt.Errorf("%w: field value must match the pattern %s", ErrRegexpNotMatch, validatorValue)
		validationErrors = append(validationErrors, ValidationError{Field: fieldName, Err: err})
	}
}
