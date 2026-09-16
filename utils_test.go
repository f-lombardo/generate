package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClipboard(t *testing.T) {
	s := "Test string"

	err := copyToClipboard(s)
	assert.NoError(t, err)

	actual, err := readFromClipboard()
	assert.NoError(t, err)

	assert.Equal(t, s, actual)
}
