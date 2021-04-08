package main

import (
	"log"
	"os"
)

func main() {
	// Parsing the arguments.
	path := os.Args[1]
	cmd := os.Args[2:]

	// Reading env dir.
	env, err := ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	// Executing commands.
	exitCode := RunCmd(cmd, env)
	os.Exit(exitCode)
}
