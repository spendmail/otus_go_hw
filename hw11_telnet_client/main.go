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

var (
	ErrParamsParsing = errors.New("unable to parse input arguments, watch the example:\n./go-telnet --timeout=10s 127.0.0.1 8000")
	MsgShutdown      = "Bye"
)

// Parses input params, returns parsing error, if exists.
func parseParams() (string, time.Duration, error) {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 5*time.Second, "connection timeout")
	flag.Parse()

	arguments := flag.Args()

	if len(arguments) != 2 {
		return "", 0, ErrParamsParsing
	}

	address := net.JoinHostPort(arguments[0], arguments[1])

	return address, timeout, nil
}

func main() {
	// Parsing input params.
	address, timeout, err := parseParams()
	if err != nil {
		log.Fatalln(err)
	}

	// TelnetClient constructor.
	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	// Connecting to server.
	err = client.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	// Closing connection after all.
	defer func() {
		_ = client.Close()
	}()

	doneChannel := make(chan interface{})
	signalsChannel := make(chan os.Signal, 1)

	signal.Notify(signalsChannel, syscall.SIGINT)

	// Listening for SIGINT signal.
	go func() {
		<-signalsChannel
		doneChannel <- nil
	}()

	// Writing to socket.
	go func() {
		err = client.Send()
		if err != nil {
			log.Println(err)
		}

		// If EOF has been sent.
		doneChannel <- nil
	}()

	// Reading from socket.
	go func() {
		err = client.Receive()
		if err != nil {
			log.Println(err)
		}
		<-doneChannel
	}()

	<-doneChannel
	close(doneChannel)

	// Final message.
	_, _ = os.Stderr.Write([]byte(MsgShutdown))
}
