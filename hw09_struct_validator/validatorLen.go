package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"unicode/utf8"
)

// Validator name.
const ValidatorLen = "len"

var ErrWrongFieldLen = errors.New("wrong field length")

// Validates the field by the "len" validator.
func ValidateLen(fieldName string, value reflect.Value, validatorValue string) {
	// Recursive call if the field is a slice of values.
	if value.Kind() == reflect.Slice {
		for j := 0; j < value.Len(); j++ {
			ValidateLen(fieldName, value.Index(j), validatorValue)
		}
		return
	}

	if value.Kind() != reflect.String {
		err := fmt.Errorf("%w: %s validator supports string type, %s given", ErrWrongFieldType, ValidatorLen, value.Kind().String())
		validationErrors = append(validationErrors, ValidationError{Field: fieldName, Err: err})
		return
	}

	validatorLength, err := strconv.Atoi(validatorValue)
	if err != nil {
		err := fmt.Errorf("%w: \"%s\" validator value can not be converted to int", err, ValidatorLen)
		validationErrors = append(validationErrors, ValidationError{Field: fieldName, Err: err})
		return
	}

	actualLength := utf8.RuneCountInString(value.String())
	if actualLength != validatorLength {
		err := fmt.Errorf("%w: field length should be %d characters, %d characters given", ErrWrongFieldLen, validatorLength, actualLength)
		validationErrors = append(validationErrors, ValidationError{Field: fieldName, Err: err})
	}
}
