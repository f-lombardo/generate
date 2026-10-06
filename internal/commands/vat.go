package commands

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"strings"
)

// VatCommand ------------------------------------------------------------------------

type VatResult struct {
	Vat string
}

func (r VatResult) String() string {
	return r.Vat
}

type VatCommand struct{}

func (cmd VatCommand) Execute(otherArgs map[string]string) (fmt.Stringer, error) {
	country := otherArgs["country"]

	if "IT" != country {
		return nil, errors.New("Invalid country: " + country + ". Only IT VAT numbers are supported at the moment")
	}

	result := italianVatNumber(rand.Intn)

	return VatResult{Vat: result}, nil
}

func italianVatNumber(randomFunction func(n int) int) string {
	var result strings.Builder
	sum := 0

	for i := range 10 {
		number := randomFunction(10)
		result.WriteString(strconv.Itoa(number))

		if isEven(i + 1) {
			number *= 2
			if number > 9 {
				number -= 9
			}
		}
		sum += number
	}
	t := sum % 10
	checkDigit := (10 - t) % 10

	result.WriteString(strconv.Itoa(checkDigit))

	return result.String()
}

func isEven(i int) bool {
	return i%2 == 0
}

func (cmd VatCommand) DiscardOutput() bool {
	return false
}

// VatSubcommand ----------------------------------------------------------------------

type VatSubcommand struct{}

func (s VatSubcommand) Name() string {
	return "vat"
}

func (s VatSubcommand) Description() string {
	return "Generates a valid VAT number"
}

func (s VatSubcommand) createFlagSet(outputWriter io.Writer) (*flag.FlagSet, *string) {
	vatCmd := flag.NewFlagSet(s.Name(), flag.ContinueOnError)
	vatCmd.SetOutput(outputWriter)
	defaultCountry := "IT"
	vatInputCountry := vatCmd.String("country", defaultCountry, "VAT code country code (e.g. IT, ES, NL) (Only IT is supported at this time)")
	vatCmd.Usage = func() {
		fmt.Fprintf(vatCmd.Output(), "Usage: generate vat [-country COUNTRY_CODE]\n\nOptions:\n")
		vatCmd.PrintDefaults()
	}
	return vatCmd, vatInputCountry
}

func (s VatSubcommand) Parse(args []string, outputWriter io.Writer) (Command, map[string]string, error) {
	vatCmd, vatInputCountry := s.createFlagSet(outputWriter)
	err := vatCmd.Parse(args)
	if err != nil {
		return nil, nil, err
	}
	otherArgs := map[string]string{
		"country": strings.ToUpper(*vatInputCountry),
	}
	return VatCommand{}, otherArgs, nil
}

func (s VatSubcommand) PrintDefaults(outputWriter io.Writer) {
	vatCmd, _ := s.createFlagSet(outputWriter)
	vatCmd.PrintDefaults()
}
