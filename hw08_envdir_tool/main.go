package main

import (
	"fmt"
	"os"
)

func main() {

	//args := os.Args[1:]
	path := os.Args[1]

	environmentMap, err := ReadDir(path)
	if err != nil {

	}

	for key, value := range environmentMap {
		fmt.Println(key, value.Value)
	}
}
