package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"golang.design/x/clipboard"
)

func StopIf(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func copyToClipboard(s string) error {
	err := clipboard.Init()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	ch, err := clipboard.Write(ctx, clipboard.FmtText, []byte(s))
	StopIf(err)

	select {
	case <-ch:
	case <-ctx.Done():
	}

	return nil
}
