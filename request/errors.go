package request

import "fmt"

type InvalidArgumentError struct {
	Err error
}

func (e InvalidArgumentError) Error() string {
	return fmt.Sprintf("Invalid Argument: %s", e.Err)
}
