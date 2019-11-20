package config

import "fmt"

type EnvironmentsNotSpecifiedError struct{}

func (e EnvironmentsNotSpecifiedError) Error() string {
	return "environments key not specified in the configuration"
}

type MissingParameterError struct{}

func (e MissingParameterError) Error() string {
	return "missing required parameter in the environments key"
}

type ParseYamlError struct {
	Err error
}

func (e ParseYamlError) Error() string {
	return fmt.Sprintf("cannot parse yaml file: %s", e.Err)
}
