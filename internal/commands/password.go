package commands

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

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
		result.WriteString(randomElement(charSet))
	}

	for i := 0; i < (length - numberOfCharSets); i++ {
		randomCharset := rand.Intn(numberOfCharSets)
		result.WriteString(randomElement(charSets[randomCharset]))
	}

	return PasswordResult{shuffle(result.String())}, nil
}

func (cmd PasswordCommand) DiscardOutput() bool {
	return true
}

func randomElement[T any](elements []T) T {
	randomIndex := rand.Intn(len(elements))
	return elements[randomIndex]
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
