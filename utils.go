package main

import (
	"fmt"
	"os"
)

func StopIf(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
