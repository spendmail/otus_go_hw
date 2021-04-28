package main

import (
	"bufio"
	"io"
	"net"
	"time"
)

const protocol = "tcp"

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type TelnetClientStruct struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

// Attempts to connect to socket.
func (c *TelnetClientStruct) Connect() error {
	var err error
	c.conn, err = net.DialTimeout(protocol, c.address, c.timeout)
	if err != nil {
		return err
	}

	return nil
}

// Closes socket connection.
func (c *TelnetClientStruct) Close() error {
	if err := c.conn.Close(); err != nil {
		return err
	}

	return nil
}

// Scans input buffer, and sends bytes to socket.
func (c *TelnetClientStruct) Send() error {
	scanner := bufio.NewScanner(c.in)
	for {
		if scanner.Scan() {
			bytes := append(scanner.Bytes(), byte('\n'))
			_, err := c.conn.Write(bytes)
			if err != nil {
				return err
			}
		} else {
			return scanner.Err()
		}
	}
}

// Scans connection, and writes bytes to out buffer.
func (c *TelnetClientStruct) Receive() error {
	scanner := bufio.NewScanner(c.conn)
	for {
		if scanner.Scan() {
			bytes := append(scanner.Bytes(), byte('\n'))
			_, err := c.out.Write(bytes)
			if err != nil {
				return err
			}
		} else {
			return nil
		}
	}
}

// Returns an instance of TelnetClientStruct.
func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &TelnetClientStruct{address: address, timeout: timeout, in: in, out: out}
}
