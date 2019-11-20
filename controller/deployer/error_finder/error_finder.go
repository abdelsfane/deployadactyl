package error_finder

import (
	"github.com/compozed/deployadactyl/interfaces"
)

type ErrorFinder struct {
	Matchers []interfaces.ErrorMatcher
}

func (e *ErrorFinder) FindErrors(responseString string) []interfaces.LogMatchedError {
	errors := make([]interfaces.LogMatchedError, 0, 0)

	if len(e.Matchers) > 0 {
		for _, matcher := range e.Matchers {
			match := matcher.Match([]byte(responseString))
			if match != nil {
				errors = append(errors, match)
			}
		}
	}
	return errors
}
