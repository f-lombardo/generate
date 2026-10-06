package commands

import (
	"fmt"
	"io"
)

type Command interface {
	Execute(otherArgs map[string]string) (fmt.Stringer, error)
	DiscardOutput() bool
}

type Subcommand interface {
	Name() string
	Description() string
	Parse(args []string, outputWriter io.Writer) (Command, map[string]string, error)
	PrintDefaults(outputWriter io.Writer)
}

func AllSubcommands() []Subcommand {
	return []Subcommand{
		IbanSubcommand{},
		UUIDSubcommand{},
		PasswordSubcommand{},
		VatSubcommand{},
	}
}

func FindSubcommand(name string) (Subcommand, bool) {
	for _, sub := range AllSubcommands() {
		if sub.Name() == name {
			return sub, true
		}
	}
	return nil, false
}
