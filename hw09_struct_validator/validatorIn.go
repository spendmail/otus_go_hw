package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
)

// Validator name.
const ValidatorIn = "in"

var (
	InSplitter        = regexp.MustCompile(`,`)
	ErrValueNotExists = errors.New("value not exists")
)

// Validates the string typed field by the "in" validator.
func ValidateInString(fieldName string, value fmt.Stringer, validatorValue string) {
	validatorValues := InSplitter.Split(validatorValue, -1)
	validatorValuesMap := make(map[string]interface{}, len(validatorValues))
	for _, v := range validatorValues {
		validatorValuesMap[v] = nil
	}

	_, ok := validatorValuesMap[value.String()]
	if !ok {
		err := fmt.Errorf("%w: field value %s must be in (%s)", ErrValueNotExists, value.String(), validatorValue)
		validationErrors = append(validationErrors, ValidationError{Field: fieldName, Err: err})
	}
}

// Validates the int typed field by the "in" validator.
func ValidateInInt(fieldName string, value reflect.Value, validatorValue string) {
	validatorValues := InSplitter.Split(validatorValue, -1)
	validatorValuesMap := make(map[int64]interface{}, len(validatorValues))
	for _, v := range validatorValues {
		intValidatorValue, err := strconv.Atoi(v)
		if err != nil {
			err := fmt.Errorf("%w: can't change type of %s to int: %s", ErrWrongValidatorValue, validatorValue, err)
			validationErrors = append(validationErrors, ValidationError{Field: fieldName, Err: err})
			return
		}
		validatorValuesMap[int64(intValidatorValue)] = nil
	}

	_, ok := validatorValuesMap[value.Int()]
	if !ok {
		err := fmt.Errorf("%w: field value %d must be in (%s)", ErrValueNotExists, value.Int(), validatorValue)
		validationErrors = append(validationErrors, ValidationError{Field: fieldName, Err: err})
	}
}

// Validates the field by the "in" validator.
func ValidateIn(fieldName string, value reflect.Value, validatorValue string) {
	// Recursive call if the field is a slice of values.
	if value.Kind() == reflect.Slice {
		for j := 0; j < value.Len(); j++ {
			ValidateIn(fieldName, value.Index(j), validatorValue)
		}
		return
	}

	switch value.Type().Kind() {
	case reflect.String:
		ValidateInString(fieldName, value, validatorValue)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		ValidateInInt(fieldName, value, validatorValue)
	default:
		err := fmt.Errorf("%w: %s validator supports numeric types and strings, %s given", ErrWrongFieldType, ValidatorIn, value.Type().Kind())
		validationErrors = append(validationErrors, ValidationError{Field: fieldName, Err: err})
		return
	}
}
