package hw10programoptimization

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

// Returns a map with domain statistics.
func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	rd := bufio.NewReader(r)
	result := make(DomainStat)
	user := &User{}

	for {
		// Reads line in bytes.
		line, readErr := rd.ReadBytes('\n')

		// Returns error, if the error is not io.EOF.
		if readErr != nil && !errors.Is(readErr, io.EOF) {
			return nil, readErr
		}

		// Unmarshalling using easyJson generated unmarshaler.
		if err := user.UnmarshalJSON(line); err != nil {
			return nil, err
		}

		// Collecting the number of domains.
		if strings.Contains(user.Email, "."+domain) {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}

		// Break in case of the end of file.
		if errors.Is(readErr, io.EOF) {
			break
		}
	}

	return result, nil
}
