package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"github.com/anatol/clevis.go"
	_ "github.com/anatol/tang.go"
	"github.com/lestrrat-go/jwx/v3/jwe"
	flag "github.com/spf13/pflag"
)

const appName = "clevis"

var (
	version    string
	libVersion string
)

func encrypt(pin string, config string) (string, error) {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}
	output, err := clevis.Encrypt(input, pin, config)
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func decrypt() (string, error) {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}
	output, err := clevis.Decrypt(input)
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func inspect() (string, error) {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}
	message, err := jwe.Parse(input)
	if err != nil {
		return "", err
	}
	output, err := json.MarshalIndent(message, "", "  ")
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func errExit(code int, msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(code)
}

func init() {
	flag.BoolP("help", "h", false, "display usage help")
	flag.Bool("version", false, "display version")
	flag.CommandLine.SortFlags = false
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Pluggable framework for automated decryption\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s <command> [options...]\n\n", appName)
		fmt.Fprintf(os.Stderr, "Commands:\n")
		fmt.Fprintf(os.Stderr, "  %s decrypt         Decrypts using the policy defined at encryption time\n", appName)
		fmt.Fprintf(os.Stderr, "  %s encrypt sss     Encrypts using a Shamir's Secret Sharing policy\n", appName)
		fmt.Fprintf(os.Stderr, "  %s encrypt tang    Encrypts using a Tang binding server policy\n", appName)
		fmt.Fprintf(os.Stderr, "  %s inspect         Dumps raw JWE\n\n", appName)
		flag.PrintDefaults()
	}
	if version == "" {
		version = time.Now().Format("2006.1.2-dev")
	}
}

func main() {
	for _, arg := range os.Args[1:] {
		switch arg {
		case "--help", "-h":
			flag.Usage()
			os.Exit(0)
		case "--version":
			fmt.Printf("v%s %s/%s (anatol/clevis v%s)\n", version, runtime.GOOS, runtime.GOARCH, libVersion)
			os.Exit(0)
		}
	}
	flag.Parse()
	command := flag.Arg(0)
	if command == "" {
		errExit(2, "Missing command, use --help for usage")
	}

	switch flag.Arg(0) {
	case "encrypt", "e":
		pin := flag.Arg(1)
		config := flag.Arg(2)
		if pin == "" || config == "" {
			errExit(2, "Missing arguments, use --help for usage")
		}

		output, err := encrypt(pin, config)
		if err != nil {
			errExit(1, "Error: %v", err)
		}
		fmt.Println(output)

	case "decrypt", "d":
		output, err := decrypt()
		if err != nil {
			errExit(1, "Error: %v", err)
		}
		fmt.Println(output)

	case "inspect", "i":
		output, err := inspect()
		if err != nil {
			errExit(1, "Error: %v", err)
		}
		fmt.Println(output)

	case "luks":
		errExit(1, "LUKS is not currently supported")

	default:
		errExit(2, "Invalid command, use --help for usage")
	}
}
