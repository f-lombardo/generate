package main

import (
	"encoding/json"
	"fmt"
)

type Formatter interface {
	Format(structResult fmt.Stringer) (string, error)
}

type JsonFormatter struct {
}

func (j JsonFormatter) Format(structResult fmt.Stringer) (string, error) {
	result, err := json.Marshal(structResult)
	StopIf(err)
	return string(result), nil
}

type TabFormatter struct {
}

func (t TabFormatter) Format(structResult fmt.Stringer) (string, error) {
	return structResult.String(), nil
}
