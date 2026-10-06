package formatters

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type dummyResult struct {
	Value string `json:"value"`
}

func (d dummyResult) String() string {
	return d.Value
}

func TestTabFormatter(t *testing.T) {
	formatter := TabFormatter{}
	data := dummyResult{Value: "hello_world"}

	result, err := formatter.Format(data)
	require.NoError(t, err)
	assert.Equal(t, "hello_world", result)
}

func TestJSONFormatter(t *testing.T) {
	formatter := JSONFormatter{}
	data := dummyResult{Value: "hello_world"}

	result, err := formatter.Format(data)
	require.NoError(t, err)
	assert.JSONEq(t, `{"value":"hello_world"}`, result)
}
