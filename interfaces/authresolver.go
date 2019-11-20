package interfaces

import (
	"github.com/compozed/deployadactyl/structs"
)

type AuthResolver interface {
	Resolve(authorization Authorization, environment structs.Environment, deploymentLogger DeploymentLogger) (Authorization, error)
}
