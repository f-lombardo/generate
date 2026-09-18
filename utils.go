package main

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"strings"
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
	if err != nil {
		return err
	}

	select {
	case <-ch:
	case <-ctx.Done():
	}

	return nil
}

func readFromClipboard() (string, error) {
	err := clipboard.Init()
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	b, err := clipboard.Read(context.Background(), clipboard.FmtText)
	if err != nil {
		return "", err
	}

	select {
	case <-ctx.Done():
	}

	return string(b), nil
}

func getGitHash() string {
	info, available := debug.ReadBuildInfo()
	if !available {
		return ""
	}
	var result strings.Builder
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision", "vcs.time":
			result.WriteString(s.Value + " ")
		}
	}
	return result.String()
}
