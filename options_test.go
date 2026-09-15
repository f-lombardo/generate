package main

import (
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
				command:      UuidCommand{},
				otherArgs:    map[string]string{"version": "4"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			options, err := readOptions(tt.args)

			require.NoError(t, err)
			require.NotNil(t, options)

			assert.Equal(t, tt.expectedOptions, options)
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
	*result = true
	return result
}
