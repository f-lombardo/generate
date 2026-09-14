package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

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

type uuidVersion string

func (f *uuidVersion) String() string {
	return string(*f)
}

func (f *uuidVersion) Set(valore string) error {
	switch valore {
	case "4", "7":
		*f = uuidVersion(valore)
		return nil
	default:
		return errors.New("UUID version should be '4' or '7' (default 4)")
	}
}

type Options struct {
	outputFormat *outputFormat
	clipboard    *bool
	operation    string
	otherArgs    map[string]string
}

func stopIf(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func copyToClipboard(result string) error {
	err := clipboard.Init()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	ch, err := clipboard.Write(ctx, clipboard.FmtText, []byte(result))
	stopIf(err)

	select {
	case <-ch:
	case <-ctx.Done():
	}

	return nil
}

type IbanResult struct {
	IBAN string
}

func (r IbanResult) String() string {
	return r.IBAN
}

type UUIDResult struct {
	UUID string
}

func (r UUIDResult) String() string {
	return r.UUID
}

func executeCommand(opts Options) (fmt.Stringer, error) {
	switch opts.operation {
	case "iban":
		code, err := iban.Generate(opts.otherArgs["country"])
		if err != nil {
			return nil, err
		}

		return IbanResult{
			IBAN: code,
		}, nil

	case "uuid":
		switch opts.otherArgs["version"] {
		case "4":
			code, err := uuid.NewRandom()
			if err != nil {
				return nil, err
			}
			return UUIDResult{
				UUID: code.String(),
			}, nil
		case "7":
			code, err := uuid.NewRandom()
			if err != nil {
				return nil, err
			}
			return UUIDResult{
				UUID: code.String(),
			}, nil
		default:
			return nil, errors.New("Invalid version: " + opts.otherArgs["version"])
		}
	default:
		return nil, errors.New("Invalid command: " + opts.operation)
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
	defaultCountry := "IT"
	inputCountry := ibanCmd.String("country", defaultCountry, "IBAN country code (e.g. IT, ES, NL)")
	ibanCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: generate iban [-country COUNTRY_CODE]\n\nOptions:\n")
		ibanCmd.PrintDefaults()
	}

	// uuid subcommand
	uuidCmd := flag.NewFlagSet("uuid", flag.ExitOnError)
	defaultVersion := "4"
	uuidVersion := uuidCmd.String("version", defaultVersion, "UUID version (4 or 7)")
	ibanCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: generate uuid [-version 7]\n\nOptions:\n")
		ibanCmd.PrintDefaults()
	}

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
		err := uuidCmd.Parse(remainingArgs[1:])
		if err != nil {
			return Options{}, err
		}
		options.operation = "uuid"
		options.otherArgs["version"] = *uuidVersion

	default:
		return Options{}, errors.New("Invalid command: " + subcommand)
	}

	return options, nil
}

func formatForOutput(structResult fmt.Stringer, format *outputFormat) (string, error) {
	switch format.String() {
	case "json":
		result, err := json.Marshal(structResult)
		stopIf(err)
		return string(result), nil
	case "tab":
		return structResult.String(), nil
	default:
		return "", errors.New("Unknown format: " + format.String())
	}
}

func main() {
	opts, err := readOptions()
	stopIf(err)

	structResult, err := executeCommand(opts)
	stopIf(err)

	result, err := formatForOutput(structResult, opts.outputFormat)
	stopIf(err)

	if *opts.clipboard {
		err = copyToClipboard(result)
		stopIf(err)
	}

	fmt.Println(result)
}
