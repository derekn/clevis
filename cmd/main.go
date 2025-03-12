package main

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"github.com/anatol/clevis.go"
	"github.com/urfave/cli/v2"
)

var (
	version    string = time.Now().Format("2006.1.2-dev")
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

func main() {
	app := &cli.App{
		Name:    "clevis",
		Usage:   "pluggable framework for automated decryption",
		Version: fmt.Sprintf("v%s %s/%s (anatol/clevis v%s)", version, runtime.GOOS, runtime.GOARCH, libVersion),
		Commands: []*cli.Command{
			{
				Name:      "encrypt",
				Usage:     "encrypts stdin",
				Aliases:   []string{"e"},
				ArgsUsage: "<pin> <config>",
				Action: func(ctx *cli.Context) error {
					if ctx.NArg() != 2 {
						_ = cli.ShowAppHelp(ctx)
						return fmt.Errorf("missing required arguments: <pin> <config>")
					}

					output, err := encrypt(ctx.Args().Get(0), ctx.Args().Get(1))
					if err != nil {
						return err
					}
					fmt.Println(output)
					return nil
				},
			},
			{
				Name:    "decrypt",
				Usage:   "decrypts stdin",
				Aliases: []string{"d"},
				Action: func(ctx *cli.Context) error {
					output, err := decrypt()
					if err != nil {
						return err
					}
					fmt.Println(output)
					return nil
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
