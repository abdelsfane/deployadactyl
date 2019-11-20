package interfaces

import (
	"github.com/compozed/deployadactyl/structs"
)

type StopManagerFactory interface {
	StopManager(deployEventData structs.DeployEventData) ActionCreator
}
