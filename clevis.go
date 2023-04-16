/*
Clevis front-end CLI
@author Derek Nicol
*/
package main

import (
	"fmt"
	"github.com/anatol/clevis.go"
	"io"
	"os"
)

const VERSION = "0.1.0"

func main() {
	var output []byte
	var err error

	if len(os.Args[1:]) < 1 {
		usage("Error: missing action")
	}
	action := os.Args[1]

	switch action {
	case "encrypt":
		if len(os.Args[2:]) != 2 {
			usage("Error: missing parameters")
		}

		input, _ := io.ReadAll(os.Stdin)
		pin := os.Args[2]
		config := os.Args[3]
		output, err = clevis.Encrypt(input, pin, config)

	case "decrypt":
		input, _ := io.ReadAll(os.Stdin)
		output, err = clevis.Decrypt(input)

	case "help":
		usage()

	case "version":
		fmt.Printf("Version %s", VERSION)

	default:
		usage("Error: unknown action")
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: %s", err)
		os.Exit(1)
	}

	fmt.Println(string(output))
}

func usage(params ...string) {
	code := 1

	if len(params) > 0 {
		fmt.Fprintln(os.Stderr, params[0]+"\n")
		code = 2
	}

	fmt.Println("Usage: clevis encrypt <pin> <config>")
	fmt.Println("       clevis decrypt")
	os.Exit(code)
}
