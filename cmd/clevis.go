package main

import (
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/anatol/clevis.go"
	"github.com/urfave/cli/v2"
)

var (
	version    string = "0.0.0-dev"
	goVersion  string
	libVersion string
)

func encrypt(pin string, config string) error {
	input, _ := io.ReadAll(os.Stdin)
	output, err := clevis.Encrypt(input, pin, config)
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}

func decrypt() error {
	input, _ := io.ReadAll(os.Stdin)
	output, err := clevis.Decrypt(input)
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}

func errorExit(err error) {
	fmt.Fprintln(os.Stderr, "Error:", err)
	os.Exit(1)
}

func main() {
	app := &cli.App{
		Name:    "clevis",
		Usage:   "pluggable framework for automated decryption",
		Version: fmt.Sprintf("v%s %s/%s %s (anatol/clevis v%s)", version, runtime.GOOS, runtime.GOARCH, goVersion, libVersion),
		Commands: []*cli.Command{
			{
				Name:  "encrypt",
				Usage: "<pin> <config>",
				Action: func(ctx *cli.Context) error {
					if ctx.NArg() != 2 {
						return fmt.Errorf("missing required arguments: <pin> <config>")
					}
					err := encrypt(ctx.Args().Get(0), ctx.Args().Get(1))
					return err
				},
			},
			{
				Name: "decrypt",
				Action: func(ctx *cli.Context) error {
					err := decrypt()
					return err
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		errorExit(err)
	}
}
