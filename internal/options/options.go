package options

import (
	"errors"
	"flag"
	"fmt"
	"io"

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

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: generate [global options] <command> [command options]\n\n")
		fmt.Fprintf(fs.Output(), "Global options:\n")
		fs.PrintDefaults()
		fmt.Fprintf(fs.Output(), "\nAvailable commands:\n\n")

		for _, sub := range commands.AllCliSubcommands() {
			fmt.Fprintf(fs.Output(), "  %s\t%s\n", sub.Name(), sub.Description())
			fmt.Fprintf(fs.Output(), "  Command options:\n")
			sub.PrintDefaults(fs.Output())
			fmt.Fprintf(fs.Output(), "\n")
		}
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

	subcommandName := remainingArgs[0]

	sub, found := commands.FindCliSubcommand(subcommandName)
	if !found {
		fs.Usage()
		return Options{}, errors.New("Invalid command: " + subcommandName)
	}

	cmd, otherArgs, err := sub.Parse(remainingArgs[1:], outputWriter)
	if err != nil {
		return Options{}, err
	}

	options.Command = cmd
	options.OtherArgs = otherArgs

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
