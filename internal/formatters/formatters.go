package formatters

import (
	"encoding/json"
	"fmt"

	"github.com/f-lombardo/generate/internal/utils"
)

type Formatter interface {
	Format(structResult fmt.Stringer) (string, error)
}

type JSONFormatter struct {
}

func (j JSONFormatter) Format(structResult fmt.Stringer) (string, error) {
	result, err := json.Marshal(structResult)
	utils.StopIf(err)
	return string(result), nil
}

type TabFormatter struct {
}

func (t TabFormatter) Format(structResult fmt.Stringer) (string, error) {
	return structResult.String(), nil
}
