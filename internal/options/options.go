package options

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/f-lombardo/generate/internal/commands"
	"github.com/f-lombardo/generate/internal/formatters"
)

type OutputFormat struct {
	Formatter formatters.Formatter
}

// Required by flag.Value
func (f *OutputFormat) String() string {
	return fmt.Sprintf("%T", f)
}

// Set Required by flag.Value. Validation will be performed here
func (f *OutputFormat) Set(valore string) error {
	switch valore {
	case "json":
		f.Formatter = formatters.JSONFormatter{}
		return nil
	case "tab":
		f.Formatter = formatters.TabFormatter{}
		return nil
	default:
		return errors.New("output format must be 'json' o 'tab' (default)")
	}
}

type PasswordLength string

func (l *PasswordLength) String() string {
	return string(*l)
}

func (l *PasswordLength) Set(valore string) error {
	length, err := strconv.Atoi(valore)
	if err != nil {
		return err
	}
	numberOfCharSets := 4

	if length < numberOfCharSets {
		return errors.New("Invalid length: " + valore)
	}

	*l = PasswordLength(valore)
	return nil
}

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

type Options struct {
	OutputFormat OutputFormat
	Clipboard    *bool
	Version      *bool
	Command      commands.Command
	OtherArgs    map[string]string
}

// Passing nil to outputWriter is the same as passing os.Stderr
func ReadOptions(args []string, outputWriter io.Writer) (Options, error) {
	options := Options{
		Clipboard: new(bool),
		Version:   new(bool),
		OtherArgs: make(map[string]string),
	}

	fs := flag.NewFlagSet("create", flag.ContinueOnError)
	fs.SetOutput(outputWriter)

	options.Version = fs.Bool("version", false, "Prints current version and exits")

	options.Clipboard = fs.Bool("clipboard", true, "Copies the results to the system clipboard. E.g. --clipboard=false")

	options.OutputFormat = defaultFormat()
	fs.Var(&options.OutputFormat, "output", "Output format. Valid values: json, tab. Default value: tab")

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
	uuidVersion := UUIDVersion(defaultVersion)
	uuidCmd.Var(&uuidVersion, "version", "UUID version (4 or 7)")
	uuidCmd.Usage = func() {
		fmt.Fprintf(uuidCmd.Output(), "Usage: generate uuid [-version uuid_version_number]\n\nOptions:\n")
		uuidCmd.PrintDefaults()
	}

	// password subcommand
	passwordCmd := flag.NewFlagSet("password", flag.ContinueOnError)
	passwordCmd.SetOutput(outputWriter)
	defaultLength := "14"
	passwordLength := PasswordLength(defaultLength)
	passwordCmd.Var(&passwordLength, "length", "length of the password (min 4)")
	passwordCmd.Usage = func() {
		fmt.Fprintf(passwordCmd.Output(), "Usage: generate password [-length n]\n\nOptions:\n")
		passwordCmd.PrintDefaults()
	}

	// vat subcommand
	vatCmd := flag.NewFlagSet("vat", flag.ContinueOnError)
	vatCmd.SetOutput(outputWriter)
	vatInputCountry := vatCmd.String("country", defaultCountry, "VAT code country code (e.g. IT, ES, NL) (Only IT is supported at this time)")
	vatCmd.Usage = func() {
		fmt.Fprintf(vatCmd.Output(), "Usage: generate vat [-country COUNTRY_CODE]\n\nOptions:\n")
		vatCmd.PrintDefaults()
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

		fmt.Fprintf(fs.Output(), "  password\tGenerates a random password\n")
		fmt.Fprintf(fs.Output(), "  Command options:\n")
		passwordCmd.PrintDefaults()
		fmt.Fprintf(fs.Output(), "\n")

		fmt.Fprintf(fs.Output(), "  vat\tGenerates a valid VAT number\n")
		fmt.Fprintf(fs.Output(), "  Command options:\n")
		vatCmd.PrintDefaults()
		fmt.Fprintf(fs.Output(), "\n")
	}

	if err := fs.Parse(args); err != nil {
		return Options{}, err
	}

	// fs.Args() Returns the remaining parameters after removing the global flags.
	// The first remaining element should be our subcommand.
	remainingArgs := fs.Args()

	if len(remainingArgs) < 1 {
		if *options.Version {
			return options, nil
		}
		fs.Usage()
		return Options{}, errors.New("no command specified")
	}

	subcommand := remainingArgs[0]

	switch subcommand {
	case "iban":
		err := ibanCmd.Parse(remainingArgs[1:])
		if err != nil {
			return Options{}, err
		}
		options.Command = commands.IbanCommand{}
		options.OtherArgs["country"] = strings.ToUpper(*inputCountry)

	case "uuid":
		err := uuidCmd.Parse(remainingArgs[1:])
		if err != nil {
			return Options{}, err
		}
		options.Command = commands.UUIDCommand{}
		options.OtherArgs["version"] = uuidVersion.String()

	case "password":
		err := passwordCmd.Parse(remainingArgs[1:])
		if err != nil {
			return Options{}, err
		}
		options.Command = commands.PasswordCommand{}
		options.OtherArgs["length"] = passwordLength.String()

	case "vat":
		err := vatCmd.Parse(remainingArgs[1:])
		if err != nil {
			return Options{}, err
		}
		options.Command = commands.VatCommand{}
		options.OtherArgs["country"] = strings.ToUpper(*vatInputCountry)

	default:
		fs.Usage()
		return Options{}, errors.New("Invalid command: " + subcommand)
	}

	return options, nil
}

func defaultFormat() OutputFormat {
	return OutputFormat{
		Formatter: formatters.TabFormatter{},
	}
}

func jsonFormat() OutputFormat {
	return OutputFormat{
		Formatter: formatters.JSONFormatter{},
	}
}
