package main

import (
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/anatol/clevis.go"
)

var (
	version string = "0.0.0-dev"
	// commitSha  string
	goVersion  string
	libVersion string
)

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
		fmt.Printf(
			"Clevis v%s %s/%s %s (anatol/clevis v%s)\n",
			version, runtime.GOOS, runtime.GOARCH, goVersion, libVersion,
		)

	default:
		usage("Error: unknown action")
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	if output != nil {
		fmt.Println(string(output))
	}
}

func usage(params ...string) {
	code := 1

	if len(params) > 0 {
		fmt.Fprintln(os.Stderr, params[0]+"\n")
		code = 2
	}

	fmt.Println("Usage: clevis encrypt <pin> <config>")
	fmt.Println("       clevis decrypt")
	fmt.Println("       clevis version")
	os.Exit(code)
}
