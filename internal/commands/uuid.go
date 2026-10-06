package commands

import (
	"errors"
	"flag"
	"fmt"
	"io"

	"github.com/google/uuid"
)

type UUIDVersion string

func (f *UUIDVersion) String() string {
	return string(*f)
}

func (f *UUIDVersion) Set(valore string) error {
	switch valore {
	case "4", "7":
		*f = UUIDVersion(valore)
		return nil
	default:
		return errors.New("UUID version should be '4' or '7' (default 4)")
	}
}

// Uuid command ------------------------------------------------------------------------

type UUIDResult struct {
	UUID string
}

func (r UUIDResult) String() string {
	return r.UUID
}

type UUIDCommand struct{}

func (cmd UUIDCommand) Execute(otherArgs map[string]string) (fmt.Stringer, error) {
	switch otherArgs["version"] {
	case "4":
		code, err := uuid.NewRandom()
		if err != nil {
			return nil, err
		}
		return UUIDResult{
			UUID: code.String(),
		}, nil
	case "7":
		code, err := uuid.NewV7()
		if err != nil {
			return nil, err
		}
		return UUIDResult{
			UUID: code.String(),
		}, nil
	default:
		return nil, errors.New("Invalid version: " + otherArgs["version"])
	}
}

func (cmd UUIDCommand) DiscardOutput() bool {
	return false
}

// UUIDCliSubcommand ---------------------------------------------------------------------

type UUIDCliSubcommand struct{}

func (s UUIDCliSubcommand) Name() string {
	return "uuid"
}

func (s UUIDCliSubcommand) Description() string {
	return "Generates a UUID v4 or v7"
}

func (s UUIDCliSubcommand) createFlagSet(outputWriter io.Writer) (*flag.FlagSet, *UUIDVersion) {
	uuidCmd := flag.NewFlagSet(s.Name(), flag.ContinueOnError)
	uuidCmd.SetOutput(outputWriter)
	defaultVersion := "4"
	uuidVersion := UUIDVersion(defaultVersion)
	uuidCmd.Var(&uuidVersion, "version", "UUID version (4 or 7)")
	uuidCmd.Usage = func() {
		fmt.Fprintf(uuidCmd.Output(), "Usage: generate uuid [-version uuid_version_number]\n\nOptions:\n")
		uuidCmd.PrintDefaults()
	}
	return uuidCmd, &uuidVersion
}

func (s UUIDCliSubcommand) Parse(args []string, outputWriter io.Writer) (Command, map[string]string, error) {
	uuidCmd, uuidVersion := s.createFlagSet(outputWriter)
	err := uuidCmd.Parse(args)
	if err != nil {
		return nil, nil, err
	}
	otherArgs := map[string]string{
		"version": uuidVersion.String(),
	}
	return UUIDCommand{}, otherArgs, nil
}

func (s UUIDCliSubcommand) PrintDefaults(outputWriter io.Writer) {
	uuidCmd, _ := s.createFlagSet(outputWriter)
	uuidCmd.PrintDefaults()
}
