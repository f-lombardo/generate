package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"
)

type OutputFormat struct {
	formatter Formatter
}

// Required by flag.Value
func (f *OutputFormat) String() string {
	return fmt.Sprintf("%T", f)
}

// Required by flag.Value. Validation will be performed here
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

// Passing nil to outputWriter is the same as passing os.Stderr
func readOptions(args []string, outputWriter io.Writer) (Options, error) {
	options := Options{
		clipboard: new(bool),
		otherArgs: make(map[string]string),
	}

	fs := flag.NewFlagSet("create", flag.ContinueOnError)
	fs.SetOutput(outputWriter)

	options.clipboard = fs.Bool("clipboard", true, "Copies the results to the system clipboard. E.g. --clipboard=false")

	options.outputFormat = defaultFormat()
	fs.Var(&options.outputFormat, "output", "Output format. Valid values: json, tab. Default value: tab")

	// iban subcommand
	ibanCmd := flag.NewFlagSet("iban", flag.ContinueOnError)
	ibanCmd.SetOutput(outputWriter)
	defaultCountry := "IT"
	inputCountry := ibanCmd.String("country", defaultCountry, "IBAN country code (e.g. IT, ES, NL)")
	ibanCmd.Usage = func() {
		fmt.Fprintf(ibanCmd.Output(), "Usage: generate iban [-country COUNTRY_CODE]\n\nOptions:\n")
		ibanCmd.PrintDefaults()
	}

	// uuid subcommand
	uuidCmd := flag.NewFlagSet("uuid", flag.ContinueOnError)
	uuidCmd.SetOutput(outputWriter)
	defaultVersion := "4"
	uuidVersion := uuidCmd.String("version", defaultVersion, "UUID version (4 or 7)")
	uuidCmd.Usage = func() {
		fmt.Fprintf(ibanCmd.Output(), "Usage: generate uuid [-version uuid_version_number]\n\nOptions:\n")
		uuidCmd.PrintDefaults()
	}

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: generate [global options] <command> [command options]\n\n")
		fmt.Fprintf(fs.Output(), "Global options:\n")
		fs.PrintDefaults()
		fmt.Fprintf(fs.Output(), "\nAvailable commands:\n\n")

		fmt.Fprintf(fs.Output(), "  iban\tGenerates a valid IBAN\n")
		fmt.Fprintf(fs.Output(), "  Command options:\n")
		ibanCmd.PrintDefaults()
		fmt.Fprintf(fs.Output(), "\n")

		fmt.Fprintf(fs.Output(), "  uuid\tGenerates a UUID v4 or v7\n")
		fmt.Fprintf(fs.Output(), "  Command options:\n")
		uuidCmd.PrintDefaults()
		fmt.Fprintf(fs.Output(), "\n")
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
