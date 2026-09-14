package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/jacoelho/banking/iban"
	"golang.design/x/clipboard"
)

// 1. Definiamo un nuovo tipo stringa
type outputFormat string

// 2. Implementiamo il metodo String() (richiesto da flag.Value)
func (f *outputFormat) String() string {
	return string(*f)
}

// 3. Implementiamo il metodo Set() (richiesto da flag.Value)
// È qui che avviene la validazione in tempo reale!
func (f *outputFormat) Set(valore string) error {
	switch valore {
	case "json", "tab":
		*f = outputFormat(valore)
		return nil
	default:
		return errors.New("output format must be 'json' o 'tab' (default)")
	}
}

type Options struct {
	outputFormat *outputFormat
	clipboard    *bool
	operation    string
	otherArgs    map[string]string
}

const defaultCountry = "IT"

func main() {
	opts, err := readOptions()
	if err != nil {
		panic(err)
	}

	result, err := executeCommand(opts)
	if err != nil {
		panic(err)
	}

	if *opts.clipboard {
		err = copyToClipboard(result)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println(result)
}

func copyToClipboard(result string) error {
	err := clipboard.Init()
	if err != nil {
		return err
	}

	ctx := context.Background()
	clipboard.Write(ctx, clipboard.FmtText, []byte(result))

	return nil
}

func executeCommand(opts Options) (string, error) {
	switch opts.operation {
	case "iban":
		code, err := iban.Generate(opts.otherArgs["country"])
		if err != nil {
			return "", err
		}

		return code, nil

	case "uuid":
		code, err := uuid.NewRandom()
		if err != nil {
			return "", err
		}
		return code.String(), nil

	default:
		return "", errors.New("Invalid command: " + opts.operation)
	}
}

func readOptions() (Options, error) {
	options := Options{
		outputFormat: new(outputFormat),
		clipboard:    new(bool),
		otherArgs:    make(map[string]string),
	}

	options.clipboard = flag.Bool("clipboard", true, "Copies the results to the system clipboard. E.g. --clipboard=false")

	*options.outputFormat = "tab"
	flag.Var(options.outputFormat, "output", "Output format. Valid values: json, tab. Default value: tab")

	// iban subcommand
	ibanCmd := flag.NewFlagSet("iban", flag.ExitOnError)
	inputCountry := ibanCmd.String("country", defaultCountry, "IBAN country code (e.g. IT, ES, NL")
	ibanCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: generate iban [-country COUNTRY_CODE]\n\nOpzioni:\n")
		ibanCmd.PrintDefaults()
	}

	// uuid subcommand
	uuid := flag.NewFlagSet("uuid", flag.ExitOnError)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: generate [global options] <command> [command options]\n\n")
		fmt.Fprintf(os.Stderr, "Global options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nAvailable commands:\n\n")

		fmt.Fprintf(os.Stderr, "  iban\tGenerates a valid IBAN\n")
		fmt.Fprintf(os.Stderr, "  Command options:\n")
		ibanCmd.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")

		fmt.Fprintf(os.Stderr, "  uuid\tGenerates a UUID v4\n\n")
	}

	// Questo legge tutto fino a quando non incontra qualcosa che non è un flag globale (es. il sotto-comando)
	flag.Parse()

	// flag.Args() restituisce i parametri rimasti dopo aver tolto i flag globali.
	// Il primo elemento rimasto DOVREBBE essere il nostro sotto-comando.
	remainingArgs := flag.Args()

	if len(remainingArgs) < 1 {
		flag.Usage()
		return Options{}, errors.New("no command specified")
	}

	subcommand := remainingArgs[0]

	switch subcommand {
	case "iban":
		// Passiamo al sotto-comando tutti gli argomenti che vengono DOPO di lui
		err := ibanCmd.Parse(remainingArgs[1:])
		if err != nil {
			return Options{}, err
		}
		options.operation = "iban"
		options.otherArgs["country"] = strings.ToUpper(*inputCountry)

	case "uuid":
		// Passiamo al sotto-comando tutti gli argomenti che vengono DOPO di lui
		err := uuid.Parse(remainingArgs[1:])
		if err != nil {
			return Options{}, err
		}
		options.operation = "uuid"

	default:
		return Options{}, errors.New("Invalid command: " + subcommand)
	}

	return options, nil
}
