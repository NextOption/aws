package lambda

import (
	"fmt"
)

var (
	ErrNoInput        = fmt.Errorf("no input")
	ErrNoFunctionName = fmt.Errorf("no function name")
)
