package state

import (
	"github.com/compozed/deployadactyl/config"
	"github.com/compozed/deployadactyl/controller/deployer"
	"github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/structs"
)

type EnvResolverConstructor func(envConfig config.Config) interfaces.EnvResolver

func NewEnvResolver(config config.Config) interfaces.EnvResolver {
	return EnvResolver{Config: config}
}

type EnvResolver struct {
	Config config.Config
}

func (e EnvResolver) Resolve(env string) (structs.Environment, error) {
	config := e.Config
	environment, ok := config.Environments[env]
	if !ok {
		return structs.Environment{}, deployer.EnvironmentNotFoundError{env}
	}
	return environment, nil
}
