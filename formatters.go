package main

import (
	"encoding/json"
	"fmt"
)

type Formatter interface {
	Format(structResult fmt.Stringer) (string, error)
}

type JSONFormatter struct {
}

func (j JSONFormatter) Format(structResult fmt.Stringer) (string, error) {
	result, err := json.Marshal(structResult)
	StopIf(err)
	return string(result), nil
}

type TabFormatter struct {
}

func (t TabFormatter) Format(structResult fmt.Stringer) (string, error) {
	return structResult.String(), nil
}
