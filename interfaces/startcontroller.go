package interfaces

import (
	"github.com/compozed/deployadactyl/structs"
)

type StartManagerFactory interface {
	StartManager(deployEventData structs.DeployEventData) ActionCreator
}
