package commands

import (
	"errors"
	"fmt"
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
