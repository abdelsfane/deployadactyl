package state

import (
	C "github.com/compozed/deployadactyl/config"
	"github.com/compozed/deployadactyl/controller/deployer"
	I "github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/structs"
)

type AuthResolverConstructor func(authConfig C.Config) I.AuthResolver

func NewAuthResolver(config C.Config) I.AuthResolver {
	return AuthResolver{Config: config}
}

type AuthResolver struct {
	Config C.Config
}

func (a AuthResolver) Resolve(authorization I.Authorization, environment structs.Environment, deploymentLogger I.DeploymentLogger) (I.Authorization, error) {
	deploymentLogger.Debug("checking for basic auth")
	if authorization.Username == "" && authorization.Password == "" {
		if environment.Authenticate == false {
			authorization.Username = a.Config.Username
			authorization.Password = a.Config.Password
		} else {
			return I.Authorization{}, deployer.BasicAuthError{}
		}
	}
	return authorization, nil
}
