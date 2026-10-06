package commands

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// Uuid command ------------------------------------------------------------------------

type UUIDResult struct {
	UUID string
}

func (r UUIDResult) String() string {
	return r.UUID
}

type UUIDCommand struct{}

func (cmd UUIDCommand) Execute(otherArgs map[string]string) (fmt.Stringer, error) {
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

func (cmd UUIDCommand) DiscardOutput() bool {
	return false
}
