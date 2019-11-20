package routemapper

import "fmt"

type MapRouteError struct {
	Route string
	Out   []byte
}

func (e MapRouteError) Error() string {
	return fmt.Sprintf("failed to map route: %s: %s", e.Route, string(e.Out))
}

type InvalidRouteError struct {
	Route string
}

func (e InvalidRouteError) Error() string {
	return fmt.Sprintf("invalid route provided, check that the domain exists in the foundation: %s", e.Route)
}

type ReadFileError struct {
	Err error
}

func (e ReadFileError) Error() string {
	return fmt.Sprintf("failed to read manifest file: %s", e.Err.Error())
}
