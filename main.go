package main

import (
	"fmt"
	"os"
)

func main() {
	opts, err := readOptions(os.Args[1:])
	StopIf(err)

	structResult, err := opts.command.Execute(opts.otherArgs)
	StopIf(err)

	result, err := opts.outputFormat.formatter.Format(structResult)
	StopIf(err)

	if *opts.clipboard {
		err = copyToClipboard(result)
		StopIf(err)
	}

	fmt.Println(result)
}
