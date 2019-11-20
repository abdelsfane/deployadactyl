package interfaces

import (
	"github.com/compozed/deployadactyl/structs"
)

type PushManagerFactory interface {
	PushManager(deployEventData structs.DeployEventData, auth Authorization, env structs.Environment, envVars map[string]string) ActionCreator
}
