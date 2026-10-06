package commands

import (
	"fmt"
)

type Command interface {
	Execute(otherArgs map[string]string) (fmt.Stringer, error)
	DiscardOutput() bool
}
