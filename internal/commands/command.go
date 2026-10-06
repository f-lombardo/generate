package commands

import (
	"fmt"
	"io"
)

type Command interface {
	Execute(otherArgs map[string]string) (fmt.Stringer, error)
	DiscardOutput() bool
}

type CliSubcommand interface {
	Name() string
	Description() string
	Parse(args []string, outputWriter io.Writer) (Command, map[string]string, error)
	PrintDefaults(outputWriter io.Writer)
}

func AllCliSubcommands() []CliSubcommand {
	return []CliSubcommand{
		IbanSubcommand{},
		UUIDSubcommand{},
		PasswordSubcommand{},
		VatSubcommand{},
	}
}

func FindCliSubcommand(name string) (CliSubcommand, bool) {
	for _, sub := range AllCliSubcommands() {
		if sub.Name() == name {
			return sub, true
		}
	}
	return nil, false
}
