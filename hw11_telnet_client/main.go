package main

import (
	"errors"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var ErrParamsParsing = errors.New("unable to parse input arguments, watch the example:\n./go-telnet --timeout=10s 127.0.0.1 8000")

// Parses input params, returns parsing error, if exists.
func parseParams() (string, time.Duration, error) {
	var timeoutRaw string
	flag.StringVar(&timeoutRaw, "timeout", "5s", "connection timeout")
	flag.Parse()

	arguments := flag.Args()

	if len(arguments) != 2 {
		return "", 0, ErrParamsParsing
	}

	address := net.JoinHostPort(arguments[0], arguments[1])

	timeout, err := time.ParseDuration(timeoutRaw)
	if err != nil {
		return "", 0, ErrParamsParsing
	}

	return address, timeout, nil
}

func main() {
	address, timeout, err := parseParams()
	if err != nil {
		log.Fatalln(err)
	}

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	err = client.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	doneChannel := make(chan interface{})
	signalsChannel := make(chan os.Signal, 1)

	signal.Notify(signalsChannel, syscall.SIGINT)

	go func() {
		<-signalsChannel
		close(doneChannel)
	}()

	go func() {
		err = client.Send()
		if err != nil {
			log.Println(err)
		}

		// If EOF has been sent.
		close(doneChannel)
	}()

	go func() {
		err = client.Receive()
		if err != nil {
			log.Println(err)
		}
		<-doneChannel
	}()

	<-doneChannel
}
