package main

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRandomPassword(t *testing.T) {
	command := PasswordCommand{}

	l := 8

	passwordResult, err := command.Execute(map[string]string{"length": strconv.Itoa(l)})

	require.NoError(t, err)

	assert.Equal(t, l, len(passwordResult.String()))

	pr := passwordResult.(PasswordResult)
	assert.True(t, pr.ContainsOneOf(symbols()), "No symbols in "+pr.String())
	assert.True(t, pr.ContainsOneOf(uppercaseLetters()), "No uppercase letters in "+pr.String())
	assert.True(t, pr.ContainsOneOf(lowercaseLetters()), "No lowercase letters in "+pr.String())
	assert.True(t, pr.ContainsOneOf(numbers()), "No numbers in "+pr.String())
}

func TestRandomPasswordGivesErrorForWrongParameter(t *testing.T) {
	tests := []struct {
		testName string
		length   string
	}{
		{
			testName: "length is not a number",
			length:   "x",
		},
		{
			testName: "length is too small",
			length:   "x",
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			command := PasswordCommand{}

			_, err := command.Execute(map[string]string{"length": "x"})

			require.Error(t, err)
		})
	}
}
