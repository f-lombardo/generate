package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jacoelho/banking/iban"
)

const defaultCountry = "IT"

func main() {
	country := defaultCountry

	if len(os.Args) > 1 {
		country = strings.ToUpper(os.Args[1])
	}

	code, err := iban.Generate(country)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(code)
}
