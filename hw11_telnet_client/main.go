package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

var (
	ErrParamsParsing = errors.New("unable to parse input arguments, watch the example:\n./go-telnet --timeout=10s 127.0.0.1 8000")
	timeout          string
	host             string
	port             string
)

func parseParams() error {
	flag.StringVar(&timeout, "timeout", "5s", "connection timeout")
	flag.Parse()

	args := flag.Args()

	if len(args) != 2 {
		return ErrParamsParsing
	}

	host = args[0]
	port = args[1]

	return nil
}

func main() {
	err := parseParams()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(timeout, host, port)
}
