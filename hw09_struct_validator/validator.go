package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
)

//package hw09structvalidator

const TagName = "validate"

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("implement me")
}

func parseValidateTag(s string) {

	fmt.Printf("%T = \"%v\"\n", s, s) //string = min:18|max:50
	re := regexp.MustCompile(`\s*([\w]+)\s*:\s*([^|]+)\s*`)
	fmt.Printf("%q\n", re.FindAllStringSubmatch(s, -1))
}

func Validate(v interface{}) error {

	value := reflect.ValueOf(v)

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return errors.New("not a struct")
	}

	valueType := value.Type()

	for i := 0; i < value.NumField(); i++ {
		fmt.Println("==========")
		//valueField := value.Field(i)
		//valueTypeField := valueType.Field(i)
		tag := valueType.Field(i).Tag

		//fmt.Printf("%T = %v\n", valueField, valueField)         // reflect.Value = 20
		//fmt.Printf("%T = %v\n", valueTypeField, valueTypeField) // reflect.StructField = {Age  int validate:"min:18|max:50" 32 [2] false}
		//fmt.Printf("%T = %v\n", tag, tag)                       // reflect.StructTag = validate:"min:18|max:50"

		tagValue, ok := tag.Lookup(TagName)
		if !ok {
			continue
		}

		if tagValue == "" {
			continue
		}

		//fmt.Printf("%T = %v\n", tagValue, tagValue) //string = min:18|max:50

		parseValidateTag(tagValue)
		//res, err := parseValidateTag(tagValue)
		//if err != nil {
		//
		//}
	}

	//fmt.Println(value)
	//fmt.Println(value.Type())
	//fmt.Println(value.Kind())
	//fmt.Println(value.Kind() == reflect.Ptr)
	fmt.Println("====================================================")

	return nil
}

type UserRole string

type User struct {
	ID     string `json:"id" validate:"len:36"`
	Name   string
	Age    int      `validate:"min:18|max:50"`
	Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
	Role   UserRole `validate:"in:admin,stuff"`
	Phones []string `validate:"len:11"`
	meta   json.RawMessage
}

func main() {

	user := User{
		ID:     "123456789012345678901234567890123456",
		Name:   "Username",
		Age:    20,
		Email:  "qwe@asd.ru",
		Role:   "admin",
		Phones: []string{"12345678901", "12345678901"},
		meta:   []byte("qwerty"),
	}

	err := Validate(user)
	if err != nil {
		fmt.Println(err)
	}

	err = Validate(&user)

	if err != nil {
		fmt.Println(err)
	}

	err = Validate(nil)

	if err != nil {
		fmt.Println(err)
	}
}
