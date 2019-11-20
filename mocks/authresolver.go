package mocks

import (
	I "github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/structs"
)

type AuthResolver struct {
	ResolveCall struct {
		Received struct {
			Authorization    I.Authorization
			Environment      structs.Environment
			DeploymentLogger I.DeploymentLogger
		}
		Returns struct {
			Authorization I.Authorization
			Error         error
		}
	}
}

func (a *AuthResolver) Resolve(authorization I.Authorization, environment structs.Environment, deploymentLogger I.DeploymentLogger) (I.Authorization, error) {
	a.ResolveCall.Received.Authorization = authorization
	a.ResolveCall.Received.Environment = environment
	a.ResolveCall.Received.DeploymentLogger = deploymentLogger

	return a.ResolveCall.Returns.Authorization, a.ResolveCall.Returns.Error
}
