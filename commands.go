package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jacoelho/banking/iban"
)

type Command interface {
	Execute(otherArgs map[string]string) (fmt.Stringer, error)
	DiscardOutput() bool
}

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

// Uuid command ------------------------------------------------------------------------

type UUIDResult struct {
	UUID string
}

func (r UUIDResult) String() string {
	return r.UUID
}

type UUIDCommand struct{}

func (cmd UUIDCommand) Execute(otherArgs map[string]string) (fmt.Stringer, error) {
	switch otherArgs["version"] {
	case "4":
		code, err := uuid.NewRandom()
		if err != nil {
			return nil, err
		}
		return UUIDResult{
			UUID: code.String(),
		}, nil
	case "7":
		code, err := uuid.NewV7()
		if err != nil {
			return nil, err
		}
		return UUIDResult{
			UUID: code.String(),
		}, nil
	default:
		return nil, errors.New("Invalid version: " + otherArgs["version"])
	}
}

func (cmd UUIDCommand) DiscardOutput() bool {
	return false
}

// Password command ------------------------------------------------------------------------

type PasswordResult struct {
	Password string
}

func (p PasswordResult) String() string {
	return p.Password
}

func (p PasswordResult) ContainsOneOf(charSet []string) bool {
	for _, char := range charSet {
		if strings.Contains(p.Password, char) {
			return true
		}
	}
	return false
}

type PasswordCommand struct{}

func (cmd PasswordCommand) Execute(otherArgs map[string]string) (fmt.Stringer, error) {
	length, err := strconv.Atoi(otherArgs["length"])
	if err != nil {
		return nil, err
	}

	charSets := [][]string{uppercaseLetters(), lowercaseLetters(), numbers(), symbols()}

	numberOfCharSets := len(charSets)

	if length < numberOfCharSets {
		return nil, errors.New("Invalid length: " + otherArgs["length"])
	}

	var result strings.Builder

	for _, charSet := range charSets {
		result.WriteString(randomChar(charSet))
	}

	for i := 0; i < (length - numberOfCharSets); i++ {
		randomCharset := rand.Intn(numberOfCharSets)
		result.WriteString(randomChar(charSets[randomCharset]))
	}

	return PasswordResult{shuffle(result.String())}, nil
}

func (cmd PasswordCommand) DiscardOutput() bool {
	return true
}

func randomChar(charSet []string) string {
	randomIndex := rand.Intn(len(charSet))
	return charSet[randomIndex]
}

func symbols() []string {
	return enumerateASCIIChars(6, '!')
}

func lowercaseLetters() []string {
	return enumerateASCIIChars(26, 'a')
}

func uppercaseLetters() []string {
	return enumerateASCIIChars(26, 'A')
}

func numbers() []string {
	return enumerateASCIIChars(10, '0')
}

func enumerateASCIIChars(numberOfChars int, startingChar rune) []string {
	letters := make([]string, numberOfChars)

	for i := range numberOfChars {
		letters[i] = string(rune(int(startingChar) + i))
	}

	return letters
}

func shuffle(s string) string {
	inRune := []rune(s)
	rand.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})
	return string(inRune)
}
