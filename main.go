package main

import (
	"fmt"
	"io"
	"os"

	"github.com/f-lombardo/generate/internal/options"
	"github.com/f-lombardo/generate/internal/utils"
)

var (
	version = "dev"
)

func main() {
	err := executeProgram(os.Args[1:], os.Stdout, os.Stderr)
	utils.StopIf(err)
}

func executeProgram(args []string, stdout io.Writer, stderr io.Writer) error {
	opts, err := options.ReadOptions(args, stderr)
	if err != nil {
		return err
	}

	if *opts.Version {
		_, err := fmt.Fprintln(stdout, "Version "+version+" - git infos "+utils.GetGitHash())
		return err
	}

	structResult, err := opts.Command.Execute(opts.OtherArgs)
	if err != nil {
		return err
	}

	result, err := opts.OutputFormat.Formatter.Format(structResult)
	if err != nil {
		return err
	}

	if *opts.Clipboard {
		err = utils.CopyToClipboard(result)
		if err != nil {
			return err
		}
	}

	if !opts.Command.DiscardOutput() {
		_, err := fmt.Fprintln(stdout, result)
		return err
	}

	return nil
}
