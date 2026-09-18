package main

import (
	"fmt"
	"io"
	"os"
)

var (
	version = "dev"
)

func main() {
	err := executeProgram(os.Args[1:], os.Stdout, os.Stderr)
	StopIf(err)
}

func executeProgram(args []string, stdout io.Writer, stderr io.Writer) error {
	opts, err := readOptions(args, stderr)
	if err != nil {
		return err
	}

	if *opts.version {
		fmt.Fprintln(stdout, "Version "+version+" - git infos "+getGitHash())
		return nil
	}

	structResult, err := opts.command.Execute(opts.otherArgs)
	if err != nil {
		return err
	}

	result, err := opts.outputFormat.formatter.Format(structResult)
	if err != nil {
		return err
	}

	if *opts.clipboard {
		err = copyToClipboard(result)
		if err != nil {
			return err
		}
	}

	if !opts.command.DiscardOutput() {
		fmt.Fprintln(stdout, result)
	}

	return nil
}
