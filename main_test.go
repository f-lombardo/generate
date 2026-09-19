package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const previousClipboardData = "previous clipboard data"

func TestGoodProgramExecutions(t *testing.T) {

	tests := []struct {
		testName       string
		args           []string
		outputVerifier func(string) error
	}{
		{
			testName: "iban default values",
			args:     []string{"iban"},
			outputVerifier: func(s string) error {
				actualResultWithoutNewLine := strings.TrimSuffix(s, "\n")
				if !strings.HasPrefix(actualResultWithoutNewLine, "IT") {
					return fmt.Errorf("expected to start with 'IT' prefix, while output is: %s", actualResultWithoutNewLine)
				}
				actualClipboardData, err := readFromClipboard()
				if err != nil {
					return err
				}
				assert.Equal(t, actualResultWithoutNewLine, actualClipboardData)
				return nil
			},
		},
		{
			testName: "iban with tab formatter",
			args:     []string{"-output", "tab", "iban"},
			outputVerifier: func(s string) error {
				actualResultWithoutNewLine := strings.TrimSuffix(s, "\n")
				if !strings.HasPrefix(actualResultWithoutNewLine, "IT") {
					return fmt.Errorf("expected to start with 'IT' prefix, while output is: %s", actualResultWithoutNewLine)
				}
				actualClipboardData, err := readFromClipboard()
				if err != nil {
					return err
				}
				assert.Equal(t, actualResultWithoutNewLine, actualClipboardData)
				return nil
			},
		},
		{
			testName: "iban default values without writing to clipboard",
			args:     []string{"-clipboard=false", "iban"},
			outputVerifier: func(s string) error {
				if !strings.HasPrefix(s, "IT") {
					return fmt.Errorf("expected to start with 'IT' prefix, while output is: %s", s)
				}
				clipboardData, err := readFromClipboard()
				if err != nil {
					return err
				}
				assert.Equal(t, previousClipboardData, clipboardData)
				return nil
			},
		},
		{
			testName: "iban with country code",
			args:     []string{"iban", "-country", "ES"},
			outputVerifier: func(s string) error {
				if !strings.HasPrefix(s, "ES") {
					return fmt.Errorf("expected to start with 'ES' prefix, while output is: %s", s)
				}
				return nil
			},
		},
		{
			testName: "iban with country code and JSON output",
			args:     []string{"-output", "json", "iban", "-country", "ES"},
			outputVerifier: func(s string) error {
				//if !strings.HasPrefix(s, "{\"IBAN\":\"ES") {
				//	return fmt.Errorf("expected to start with 'ES' prefix, while output is: %s", s)
				//}
				var object IbanResult
				decoder := json.NewDecoder(strings.NewReader(s))
				decoder.DisallowUnknownFields()
				if err := decoder.Decode(&object); err != nil {
					return fmt.Errorf("expected a valid JSON, while output is: %s", s)
				}
				if !strings.HasPrefix(object.IBAN, "ES") {
					return fmt.Errorf("expected to start with 'ES' prefix, while output is: %s", object.IBAN)
				}
				return nil
			},
		},
		{
			testName: "uuid default values",
			args:     []string{"uuid"},
			outputVerifier: func(s string) error {
				actualResultWithoutNewLine := strings.TrimSuffix(s, "\n")
				if err := validateUUID("4", actualResultWithoutNewLine); err != nil {
					return err
				}
				actualClipboardData, err := readFromClipboard()
				if err != nil {
					return err
				}
				assert.Equal(t, actualResultWithoutNewLine, actualClipboardData)
				return nil
			},
		},
		{
			testName: "uuid v7",
			args:     []string{"uuid", "-version", "7"},
			outputVerifier: func(s string) error {
				actualResultWithoutNewLine := strings.TrimSuffix(s, "\n")
				if err := validateUUID("7", actualResultWithoutNewLine); err != nil {
					return err
				}
				actualClipboardData, err := readFromClipboard()
				if err != nil {
					return err
				}
				assert.Equal(t, actualResultWithoutNewLine, actualClipboardData)
				return nil
			},
		},
		{
			testName: "password default values",
			args:     []string{"password"},
			outputVerifier: func(s string) error {
				actualResultWithoutNewLine := strings.TrimSuffix(s, "\n")
				assert.Empty(t, actualResultWithoutNewLine)
				actualClipboardData, err := readFromClipboard()
				if err != nil {
					return err
				}
				assert.Equal(t, 14, len(actualClipboardData))
				return nil
			},
		},
		{
			testName: "password with length value",
			args:     []string{"password", "-length", "5"},
			outputVerifier: func(s string) error {
				actualResultWithoutNewLine := strings.TrimSuffix(s, "\n")
				assert.Empty(t, actualResultWithoutNewLine)
				actualClipboardData, err := readFromClipboard()
				if err != nil {
					return err
				}
				assert.Equal(t, 5, len(actualClipboardData))
				return nil
			},
		},
		{
			testName: "version",
			args:     []string{"-version"},
			outputVerifier: func(s string) error {
				actualResultWithoutNewLine := strings.TrimSuffix(s, "\n")
				assert.True(t, strings.HasPrefix(actualResultWithoutNewLine, "Version "+version), "Wrong vesion: "+actualResultWithoutNewLine)
				return nil
			},
		},
		{
			testName: "vat default values",
			args:     []string{"vat"},
			outputVerifier: func(s string) error {
				actualResultWithoutNewLine := strings.TrimSuffix(s, "\n")
				re := regexp.MustCompile(`^\d{11}$`)
				assert.True(t, re.MatchString(actualResultWithoutNewLine), "Wrong VAT number: "+actualResultWithoutNewLine)
				return nil
			},
		},
	}

	for _, tt := range tests {
		// Before each
		err := copyToClipboard(previousClipboardData)
		require.NoError(t, err)

		t.Run(tt.testName, func(t *testing.T) {
			var outputWriter bytes.Buffer
			var errorWriter bytes.Buffer

			err := executeProgram(tt.args, &outputWriter, &errorWriter)

			require.NoError(t, err)

			require.Empty(t, errorWriter.String())

			if err = tt.outputVerifier(outputWriter.String()); err != nil {
				assert.Fail(t, err.Error())
			}
		})
	}
}

func TestWrongProgramExecutions(t *testing.T) {

	tests := []struct {
		testName        string
		args            []string
		expectedMessage string
		expectedError   string
	}{
		{
			testName:        "no parameters",
			args:            []string{},
			expectedMessage: "generate [global options] <command> [command options]",
			expectedError:   "no command specified",
		},
		{
			testName:        "wrong command",
			args:            []string{"wrong-command"},
			expectedMessage: "generate [global options] <command> [command options]",
			expectedError:   "Invalid command: wrong-command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			var outputWriter bytes.Buffer
			var errorWriter bytes.Buffer

			err := executeProgram(tt.args, &outputWriter, &errorWriter)

			require.Error(t, err)

			assert.Equal(t, tt.expectedError, err.Error())

			require.Empty(t, outputWriter.String())

			actualMessage := errorWriter.String()
			assert.True(t, strings.Contains(actualMessage, tt.expectedMessage), "Actual message: "+actualMessage)
		})
	}
}

func validateUUID(version string, s string) error {
	if err := uuid.Validate(s); err != nil {
		return err
	}

	if s[14] != version[0] {
		return fmt.Errorf("uuid %s is not of version %s", s, version)
	}

	return nil
}
