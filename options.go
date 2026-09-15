package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

// 1. Definiamo un nuovo tipo stringa
type OutputFormat struct {
	formatter Formatter
}

// 2. Implementiamo il metodo String() (richiesto da flag.Value)
func (f *OutputFormat) String() string {
	return fmt.Sprintf("%T", f)
}

// 3. Implementiamo il metodo Set() (richiesto da flag.Value)
// È qui che avviene la validazione in tempo reale!
func (f *OutputFormat) Set(valore string) error {
	switch valore {
	case "json":
		f.formatter = JsonFormatter{}
		return nil
	case "tab":
		f.formatter = TabFormatter{}
		return nil
	default:
		return errors.New("output format must be 'json' o 'tab' (default)")
	}
}

type UuidVersion string

func (f *UuidVersion) String() string {
	return string(*f)
}

func (f *UuidVersion) Set(valore string) error {
	switch valore {
	case "4", "7":
		*f = UuidVersion(valore)
		return nil
	default:
		return errors.New("UUID version should be '4' or '7' (default 4)")
	}
}

type Options struct {
	outputFormat OutputFormat
	clipboard    *bool
	command      Command
	otherArgs    map[string]string
}

func readOptions(args []string) (Options, error) {
	options := Options{
		clipboard: new(bool),
		otherArgs: make(map[string]string),
	}

	fs := flag.NewFlagSet("create", flag.ContinueOnError)

	fs.ErrorHandling()

	options.clipboard = fs.Bool("clipboard", true, "Copies the results to the system clipboard. E.g. --clipboard=false")

	options.outputFormat = defaultFormat()
	fs.Var(&options.outputFormat, "output", "Output format. Valid values: json, tab. Default value: tab")

	ibanCmd := flag.NewFlagSet("iban", flag.ContinueOnError)
	defaultCountry := "IT"
	inputCountry := ibanCmd.String("country", defaultCountry, "IBAN country code (e.g. IT, ES, NL)")
	ibanCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: generate iban [-country COUNTRY_CODE]\n\nOptions:\n")
		ibanCmd.PrintDefaults()
	}

	// uuid subcommand
	uuidCmd := flag.NewFlagSet("uuid", flag.ContinueOnError)
	defaultVersion := "4"
	uuidVersion := uuidCmd.String("version", defaultVersion, "UUID version (4 or 7)")
	uuidCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: generate uuid [-version 7]\n\nOptions:\n")
		uuidCmd.PrintDefaults()
	}

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: generate [global options] <command> [command options]\n\n")
		fmt.Fprintf(os.Stderr, "Global options:\n")
		fs.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nAvailable commands:\n\n")

		fmt.Fprintf(os.Stderr, "  iban\tGenerates a valid IBAN\n")
		fmt.Fprintf(os.Stderr, "  Command options:\n")
		ibanCmd.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")

		fmt.Fprintf(os.Stderr, "  uuid\tGenerates a UUID v4 or v7\n")
		fmt.Fprintf(os.Stderr, "  Command options:\n")
		uuidCmd.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}

	if err := fs.Parse(args); err != nil {
		return Options{}, err
	}

	// fs.Args() restituisce i parametri rimasti dopo aver tolto i flag globali.
	// Il primo elemento rimasto DOVREBBE essere il nostro sotto-comando.
	remainingArgs := fs.Args()

	if len(remainingArgs) < 1 {
		fs.Usage()
		return Options{}, errors.New("no command specified")
	}

	subcommand := remainingArgs[0]

	// iban subcommand
	switch subcommand {
	case "iban":
		// Passiamo al sotto-comando tutti gli argomenti che vengono DOPO di lui
		err := ibanCmd.Parse(remainingArgs[1:])
		if err != nil {
			return Options{}, err
		}
		options.command = IbanCommand{}
		options.otherArgs["country"] = strings.ToUpper(*inputCountry)

	case "uuid":
		// Passiamo al sotto-comando tutti gli argomenti che vengono DOPO di lui
		err := uuidCmd.Parse(remainingArgs[1:])
		if err != nil {
			return Options{}, err
		}
		options.command = UuidCommand{}
		options.otherArgs["version"] = *uuidVersion

	default:
		return Options{}, errors.New("Invalid command: " + subcommand)
	}

	return options, nil
}

func defaultFormat() OutputFormat {
	return OutputFormat{
		formatter: TabFormatter{},
	}
}

func jsonFormat() OutputFormat {
	return OutputFormat{
		formatter: JsonFormatter{},
	}
}
