package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

const SourceString = "Hello, OTUS!"

func main() {
	reversedString := stringutil.Reverse(SourceString)

	fmt.Println(reversedString)
}
