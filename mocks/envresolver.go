package mocks

import (
	"github.com/compozed/deployadactyl/structs"
)

type EnvResolver struct {
	ResolveCall struct {
		Received struct {
			Environment string
		}
		Returns struct {
			Environment structs.Environment
			Error       error
		}
	}
}

func (e *EnvResolver) Resolve(environment string) (structs.Environment, error) {
	e.ResolveCall.Received.Environment = environment

	return e.ResolveCall.Returns.Environment, e.ResolveCall.Returns.Error
}
