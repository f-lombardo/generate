package commands

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
