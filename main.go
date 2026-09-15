package main

import (
	"fmt"
	"os"
)

func main() {
	executeProgram(os.Stderr, os.Stdout)
}

func executeProgram(stderr *os.File, stdout *os.File) {
	opts, err := readOptions(os.Args[1:], stderr)
	StopIf(err)

	structResult, err := opts.command.Execute(opts.otherArgs)
	StopIf(err)

	result, err := opts.outputFormat.formatter.Format(structResult)
	StopIf(err)

	if *opts.clipboard {
		err = copyToClipboard(result)
		StopIf(err)
	}

	fmt.Fprintln(stdout, result)
}
