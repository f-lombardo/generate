package commands

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/jacoelho/banking/iban"
)

// IbanCommand ------------------------------------------------------------------------

type IbanResult struct {
	IBAN string
}

func (r IbanResult) String() string {
	return r.IBAN
}

type IbanCommand struct{}

func (cmd IbanCommand) Execute(otherArgs map[string]string) (fmt.Stringer, error) {
	code, err := iban.Generate(otherArgs["country"])
	if err != nil {
		return nil, err
	}

	return IbanResult{
		IBAN: code,
	}, nil
}

func (cmd IbanCommand) DiscardOutput() bool {
	return false
}

// IbanCliSubcommand ---------------------------------------------------------------------

type IbanCliSubcommand struct{}

func (s IbanCliSubcommand) Name() string {
	return "iban"
}

func (s IbanCliSubcommand) Description() string {
	return "Generates a valid IBAN"
}

func (s IbanCliSubcommand) createFlagSet(outputWriter io.Writer) (*flag.FlagSet, *string) {
	ibanCmd := flag.NewFlagSet(s.Name(), flag.ContinueOnError)
	ibanCmd.SetOutput(outputWriter)
	defaultCountry := "IT"
	inputCountry := ibanCmd.String("country", defaultCountry, "IBAN country code (e.g. IT, ES, NL)")
	ibanCmd.Usage = func() {
		fmt.Fprintf(ibanCmd.Output(), "Usage: generate iban [-country COUNTRY_CODE]\n\nOptions:\n")
		ibanCmd.PrintDefaults()
	}
	return ibanCmd, inputCountry
}

func (s IbanCliSubcommand) Parse(args []string, outputWriter io.Writer) (Command, map[string]string, error) {
	ibanCmd, inputCountry := s.createFlagSet(outputWriter)
	err := ibanCmd.Parse(args)
	if err != nil {
		return nil, nil, err
	}
	otherArgs := map[string]string{
		"country": strings.ToUpper(*inputCountry),
	}
	return IbanCommand{}, otherArgs, nil
}

func (s IbanCliSubcommand) PrintDefaults(outputWriter io.Writer) {
	ibanCmd, _ := s.createFlagSet(outputWriter)
	ibanCmd.PrintDefaults()
}
