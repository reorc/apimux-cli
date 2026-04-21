package main

import (
	"context"
	"fmt"
	"os"

	"github.com/reorc/apimux-cli/internal/command"
)

func main() {
	root := command.NewRoot(os.Stdout, os.Stderr)
	exitCode, err := root.Execute(context.Background(), os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	os.Exit(exitCode)
}
