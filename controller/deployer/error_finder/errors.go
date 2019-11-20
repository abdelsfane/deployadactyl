package error_finder

import "github.com/compozed/deployadactyl/interfaces"

func CreateLogMatchedError(description string, details []string, solution, code string) interfaces.LogMatchedError {
	return &logMatchedError{description: description, details: details, solution: solution, code: code}
}

type logMatchedError struct {
	description string
	details     []string
	solution    string
	code        string
}

func (e *logMatchedError) Code() string {
	return e.code
}

func (e *logMatchedError) Error() string {
	return e.description
}

func (e *logMatchedError) Details() []string {
	return e.details
}

func (e *logMatchedError) Solution() string {
	return e.solution
}
