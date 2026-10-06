package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClipboard(t *testing.T) {
	s := "Test string"

	err := CopyToClipboard(s)
	assert.NoError(t, err)

	actual, err := ReadFromClipboard()
	assert.NoError(t, err)

	assert.Equal(t, s, actual)
}
