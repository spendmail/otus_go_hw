package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// Validator name.
const ValidatorMax = "max"

var ErrValueExceededMaximum = errors.New("value exceeded maximum")

// Validates the field by the "max" validator.
func ValidateMax(fieldName string, value reflect.Value, validatorValue string) {
	// Recursive call if the field is a slice of values.
	if reflect.Slice == value.Kind() {
		for j := 0; j < value.Len(); j++ {
			ValidateMax(fieldName, value.Index(j), validatorValue)
		}
		return
	}

	intValidatorValue, err := strconv.Atoi(validatorValue)
	if err != nil {
		err := fmt.Errorf("%w: can't change type of %s to int: %s", ErrWrongValidatorValue, validatorValue, err)
		validationErrors = append(validationErrors, ValidationError{Field: fieldName, Err: err})
		return
	}

	switch value.Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if value.Int() > int64(intValidatorValue) {
			err := fmt.Errorf("%w: %s value must be less than %s", ErrValueExceededMaximum, ValidatorMax, validatorValue)
			validationErrors = append(validationErrors, ValidationError{Field: fieldName, Err: err})
			return
		}
	case reflect.Float32, reflect.Float64:
		if value.Float() > float64(intValidatorValue) {
			err := fmt.Errorf("%w: %s value must be less than %s", ErrValueExceededMaximum, ValidatorMax, validatorValue)
			validationErrors = append(validationErrors, ValidationError{Field: fieldName, Err: err})
			return
		}
	default:
		err := fmt.Errorf("%w: %s validator supports numeric types, %s given", ErrWrongFieldType, ValidatorMax, value.Type().Kind())
		validationErrors = append(validationErrors, ValidationError{Field: fieldName, Err: err})
		return
	}
}
