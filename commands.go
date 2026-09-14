package main

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jacoelho/banking/iban"
)

type Command interface {
	Execute(otherArgs map[string]string) (fmt.Stringer, error)
}

// IbanCommand ------------------------------------------------------------------------

type IbanResult struct {
	IBAN string
}

func (r IbanResult) String() string {
	return r.IBAN
}

type IbanCommand struct{}

func (cmd IbanCommand) Execute(otherArgs map[string]string) (fmt.Stringer, error) {
	code, err := iban.Generate(otherArgs["country"])
	if err != nil {
		return nil, err
	}

	return IbanResult{
		IBAN: code,
	}, nil
}

// Uuid command ------------------------------------------------------------------------

type UUIDResult struct {
	UUID string
}

func (r UUIDResult) String() string {
	return r.UUID
}

type UuidCommand struct{}

func (cmd UuidCommand) Execute(otherArgs map[string]string) (fmt.Stringer, error) {
	switch otherArgs["version"] {
	case "4":
		code, err := uuid.NewRandom()
		if err != nil {
			return nil, err
		}
		return UUIDResult{
			UUID: code.String(),
		}, nil
	case "7":
		code, err := uuid.NewV7()
		if err != nil {
			return nil, err
		}
		return UUIDResult{
			UUID: code.String(),
		}, nil
	default:
		return nil, errors.New("Invalid version: " + otherArgs["version"])
	}
}
