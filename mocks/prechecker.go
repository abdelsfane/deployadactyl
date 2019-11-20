package mocks

import (
	S "github.com/compozed/deployadactyl/structs"
)

// Prechecker handmade mock for tests.
type Prechecker struct {
	AssertAllFoundationsUpCall struct {
		Received struct {
			Environment S.Environment
		}
		Returns struct {
			Error error
		}
	}
}

// AssertAllFoundationsUp mock method.
func (p *Prechecker) AssertAllFoundationsUp(environment S.Environment) error {
	p.AssertAllFoundationsUpCall.Received.Environment = environment

	return p.AssertAllFoundationsUpCall.Returns.Error
}
