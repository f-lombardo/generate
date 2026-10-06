package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComputeCodiceFiscale(t *testing.T) {
	tests := []struct {
		testName     string
		firstName    string
		lastName     string
		gender       string
		day          int
		month        int
		year         int
		placeCode    string
		expectedCode string
	}{
		{
			testName:     "Mario Rossi",
			firstName:    "Mario",
			lastName:     "Rossi",
			gender:       "M",
			day:          10,
			month:        7,
			year:         1980,
			placeCode:    "H501",
			expectedCode: "RSSMRA80L10H501Z",
		},
		{
			testName:     "Anna Bianchi",
			firstName:    "Anna",
			lastName:     "Bianchi",
			gender:       "F",
			day:          5,
			month:        5,
			year:         1992,
			placeCode:    "E507",
			expectedCode: "BNCNNA92E45E507E",
		},
		{
			testName:     "Giuseppe Verdi",
			firstName:    "Giuseppe",
			lastName:     "Verdi",
			gender:       "M",
			day:          27,
			month:        1,
			year:         1975,
			placeCode:    "F205",
			expectedCode: "VRDGPP75A27F205T",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			cf := computeItalianFiscalCode(tt.firstName, tt.lastName, tt.gender, tt.day, tt.month, tt.year, tt.placeCode)
			assert.Equal(t, tt.expectedCode, cf)
		})
	}
}

func TestPersonCommand(t *testing.T) {
	cmd := PersonCommand{}

	result, err := cmd.Execute(map[string]string{"country": "IT"})
	require.NoError(t, err)

	person, ok := result.(PersonResult)
	require.True(t, ok)

	assert.NotEmpty(t, person.FirstName)
	assert.NotEmpty(t, person.LastName)
	assert.True(t, person.Gender == "M" || person.Gender == "F")
	assert.NotEmpty(t, person.DateOfBirth)
	assert.NotEmpty(t, person.PlaceOfBirth)
	assert.Len(t, person.FiscalCode, 16)
}

func TestPersonCommandGivesErrorForUnsupportedCountry(t *testing.T) {
	cmd := PersonCommand{}

	_, err := cmd.Execute(map[string]string{"country": "US"})
	require.Error(t, err)
	assert.Equal(t, "Invalid country: US. Only IT is supported at the moment", err.Error())
}
