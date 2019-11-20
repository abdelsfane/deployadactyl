package prechecker

import "fmt"

type NoFoundationsConfiguredError struct{}

func (e NoFoundationsConfiguredError) Error() string {
	return "no foundations configured"
}

type InvalidGetRequestError struct {
	FoundationURL string
	Err           error
}

func (e InvalidGetRequestError) Error() string {
	return fmt.Sprintf("error building request to url %s: %s", e.FoundationURL, e.Err)
}

type FoundationUnavailableError struct {
	FoundationURL string
	Status        string
}

func (e FoundationUnavailableError) Error() string {
	return fmt.Sprintf("deploy aborted: one or more CF foundations unavailable: %s: %s", e.FoundationURL, e.Status)
}
