package interfaces

import (
	"github.com/compozed/deployadactyl/structs"
)

type DeleteManagerFactory interface {
	DeleteManager(deployEventData structs.DeployEventData) ActionCreator
}
