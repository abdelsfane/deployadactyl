// Package geterrors is used build a list of missing values to reduce error checking.
package geterrors

import (
	"fmt"
	"strings"
)

// WrapFunc takes a function for getting values.
//
// Returns an ErrGetter.
func WrapFunc(get func(string) string) ErrGetter {
	return ErrGetter{get: get}
}

// ErrGetter has a get function and an array of missing keys for get calls.
type ErrGetter struct {
	get         func(string) string
	missingKeys []string
}

// Get takes a key value and uses the function from the WrapFunc method.
// If the key is missing it makes a slice of missing keys.
func (g *ErrGetter) Get(key string) string {
	val := g.get(key)

	if len(val) == 0 {
		g.missingKeys = append(g.missingKeys, key)
	}

	return val
}

// Err takes a message.
//
// Returns an error string prepended with a message and a list of the missing keys.
func (g *ErrGetter) Err(message string) error {
	if len(g.missingKeys) == 0 {
		return nil
	}

	return fmt.Errorf("%s: %s", message, strings.Join(g.missingKeys, ", "))
}
