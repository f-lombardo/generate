package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseGoodArgs(t *testing.T) {

	tests := []struct {
		testName        string
		args            []string
		expectedOptions Options
	}{
		{
			testName: "iban default values",
			args:     []string{"iban"},
			expectedOptions: Options{
				outputFormat: defaultFormat(),
				clipboard:    trueValuePointer(),
				version:      falseValuePointer(),
				command:      IbanCommand{},
				otherArgs:    map[string]string{"country": "IT"},
			},
		},
		{
			testName: "uuid default values",
			args:     []string{"uuid"},
			expectedOptions: Options{
				outputFormat: defaultFormat(),
				clipboard:    trueValuePointer(),
				version:      falseValuePointer(),
				command:      UUIDCommand{},
				otherArgs:    map[string]string{"version": "4"},
			},
		},
		{
			testName: "iban default values with global options",
			args:     []string{"-clipboard=false", "-output", "json", "iban"},
			expectedOptions: Options{
				outputFormat: jsonFormat(),
				clipboard:    falseValuePointer(),
				version:      falseValuePointer(),
				command:      IbanCommand{},
				otherArgs:    map[string]string{"country": "IT"},
			},
		},
		{
			testName: "uuid default values with global options",
			args:     []string{"-clipboard=false", "-output", "json", "uuid"},
			expectedOptions: Options{
				outputFormat: jsonFormat(),
				clipboard:    falseValuePointer(),
				version:      falseValuePointer(),
				command:      UUIDCommand{},
				otherArgs:    map[string]string{"version": "4"},
			},
		},
		{
			testName: "version flag",
			args:     []string{"-version"},
			expectedOptions: Options{
				outputFormat: defaultFormat(),
				clipboard:    trueValuePointer(),
				version:      trueValuePointer(),
				command:      nil,
				otherArgs:    map[string]string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			options, err := readOptions(tt.args, os.Stderr)

			require.NoError(t, err)
			require.NotNil(t, options)

			assert.Equal(t, tt.expectedOptions, options)
		})
	}
}

func TestHelpArgs(t *testing.T) {

	tests := []struct {
		testName        string
		args            []string
		expectedMessage string
	}{
		{
			testName:        "global help",
			args:            []string{"-help"},
			expectedMessage: "generate [global options] <command> [command options]",
		},
		{
			testName:        "iban help",
			args:            []string{"iban", "-help"},
			expectedMessage: "generate iban [-country COUNTRY_CODE]",
		},
		{
			testName:        "uuid help",
			args:            []string{"uuid", "-help"},
			expectedMessage: "generate uuid [-version uuid_version_number]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			var buf bytes.Buffer
			_, err := readOptions(tt.args, &buf)

			require.Error(t, err)

			assert.Equal(t, "flag: help requested", err.Error())

			actualMessage := buf.String()
			assert.True(t, strings.Contains(actualMessage, tt.expectedMessage), "Actual message: "+actualMessage)
		})
	}
}

func TestWrongArgs(t *testing.T) {

	tests := []struct {
		testName        string
		args            []string
		expectedMessage string
		expectedError   string
	}{
		{
			testName:        "no command",
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
		{
			testName:        "wrong uuid version",
			args:            []string{"uuid", "-version", "99"},
			expectedMessage: "generate uuid [-version uuid_version_number]",
			expectedError:   "invalid value \"99\" for flag -version: UUID version should be '4' or '7' (default 4)",
		},
		{
			testName:        "wrong iban option",
			args:            []string{"iban", "-wrong-option"},
			expectedMessage: "generate iban [-country COUNTRY_CODE]",
			expectedError:   "flag provided but not defined: -wrong-option",
		},
		{
			testName:        "wrong password option",
			args:            []string{"password", "-wrong-option"},
			expectedMessage: "generate password [-length n]",
			expectedError:   "flag provided but not defined: -wrong-option",
		},
		{
			testName:        "wrong length option",
			args:            []string{"password", "-length", "not-a-number"},
			expectedMessage: "generate password [-length n]",
			expectedError:   "invalid value \"not-a-number\" for flag -length: strconv.Atoi: parsing \"not-a-number\": invalid syntax",
		},
		{
			testName:        "small length option",
			args:            []string{"password", "-length", "1"},
			expectedMessage: "generate password [-length n]",
			expectedError:   "invalid value \"1\" for flag -length: Invalid length: 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			var buf bytes.Buffer
			_, err := readOptions(tt.args, &buf)

			require.Error(t, err)

			assert.Equal(t, tt.expectedError, err.Error())

			actualMessage := buf.String()
			assert.True(t, strings.Contains(actualMessage, tt.expectedMessage), "Actual message: "+actualMessage)
		})
	}
}

func trueValuePointer() *bool {
	result := new(bool)
	*result = true
	return result
}

func falseValuePointer() *bool {
	result := new(bool)
	*result = false
	return result
}
