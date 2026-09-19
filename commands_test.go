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
		testName      string
		length        string
		expectedError string
	}{
		{
			testName:      "length is not a number",
			length:        "x",
			expectedError: `strconv.Atoi: parsing "x": invalid syntax`,
		},
		{
			testName:      "length is too small",
			length:        "2",
			expectedError: "Invalid length: 2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			command := PasswordCommand{}

			_, err := command.Execute(map[string]string{"length": tt.length})

			require.Error(t, err)
			assert.Equal(t, tt.expectedError, err.Error())
		})
	}
}

func TestItalianVatNumber(t *testing.T) {
	tests := []struct {
		vatNumberNoCheckDigit string
		checkDigit            string
	}{
		{
			vatNumberNoCheckDigit: "0764352056",
			checkDigit:            "7",
		},
		{
			vatNumberNoCheckDigit: "1091466015",
			checkDigit:            "3",
		},
		{
			vatNumberNoCheckDigit: "0181884043",
			checkDigit:            "9",
		},
	}
	for i, tt := range tests {
		t.Run("VAT example "+strconv.Itoa(i), func(t *testing.T) {
			mockRandom := mockRandom{vatNumberNoCheckDigit: tt.vatNumberNoCheckDigit}

			assert.Equal(t, tt.vatNumberNoCheckDigit+tt.checkDigit, italianVatNumber(mockRandom.function()))
		})
	}
}

type mockRandom struct {
	vatNumberNoCheckDigit string
	currentChar           int
}

func (r mockRandom) function() func(n int) int {
	return func(n int) int {
		numberToGenerate, _ := strconv.Atoi(string(r.vatNumberNoCheckDigit[r.currentChar]))

		r.currentChar++

		return numberToGenerate
	}
}
