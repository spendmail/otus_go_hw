package hw09structvalidator

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"regexp"
)

// Validation tag name.
const TagName = "validate"

// Struct for the validation error.
type ValidationError struct {
	Field string
	Err   error
}

// Validation errors slice type.
type ValidationErrors []ValidationError

// Implementation of the error interface by validation accumulator.
func (v ValidationErrors) Error() string {
	var buffer bytes.Buffer

	for _, validationError := range v {
		buffer.WriteString(fmt.Sprintf("field: %v, error: %v\n", validationError.Field, validationError.Err))
	}

	return buffer.String()
}

var (
	validationErrors           ValidationErrors
	ErrNotAStruct              = errors.New("given type is not a struct")
	ErrWrongFieldType          = errors.New("wrong field type")
	ErrWrongValidatorValue     = errors.New("wrong validator value")
	ErrNotImplementedValidator = errors.New("validator is not implemented")
	tagParsingRegexp           = regexp.MustCompile(`\s*([\w]+)\s*:\s*([^|]+)\s*`)
)

// Parses tag value with regular expressions.
func parseValidateTag(s string) [][]string {
	return tagParsingRegexp.FindAllStringSubmatch(s, -1)
}

// Explores given structure and returns either error or slice of validation errors.
func Validate(v interface{}) error {
	// Initializing a slice for validation errors.
	validationErrors = ValidationErrors{}

	// Getting reflect.Value of the struct.
	value := reflect.ValueOf(v)

	// If pointer given, getting appropriate value.
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	// Checking whether a struct given or not.
	if value.Kind() != reflect.Struct {
		return ErrNotAStruct
	}

	// Getting reflect type of the struct.
	valueType := value.Type()

	// Struct types loop.
	for i := 0; i < value.NumField(); i++ {
		// Getting the field of the struct.
		field := valueType.Field(i)

		// Getting the tag of the struct
		tag := valueType.Field(i).Tag

		// Searching for the validation tag.
		tagValue, ok := tag.Lookup(TagName)
		if !ok {
			continue
		}

		// If validation tag is empty.
		if tagValue == "" {
			continue
		}

		// Splitting up the tag by values.
		tagParts := parseValidateTag(tagValue)

		// Selecting an appropriate validator.
		for _, tagPart := range tagParts {
			validatorName, validatorValue := tagPart[1], tagPart[2]
			switch validatorName {
			case ValidatorLen:
				ValidateLen(field.Name, value.Field(i), validatorValue)
			case ValidatorRegexp:
				ValidateRegexp(field.Name, value.Field(i), validatorValue)
			case ValidatorIn:
				ValidateIn(field.Name, value.Field(i), validatorValue)
			case ValidatorMin:
				ValidateMin(field.Name, value.Field(i), validatorValue)
			case ValidatorMax:
				ValidateMax(field.Name, value.Field(i), validatorValue)
			default:
				err := fmt.Errorf("%w: %s validator is not implemented", ErrNotImplementedValidator, validatorName)
				validationErrors = append(validationErrors, ValidationError{Field: field.Name, Err: err})
			}
		}
	}

	return validationErrors
}
