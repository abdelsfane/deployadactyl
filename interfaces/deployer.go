package interfaces

import (
	"io"

	"github.com/compozed/deployadactyl/structs"
)

type DeployResponse struct {
	StatusCode     int
	DeploymentInfo *structs.DeploymentInfo
	Error          error
}

// Deployer interface.
type Deployer interface {
	Deploy(
		deploymentInfo *structs.DeploymentInfo,
		environment structs.Environment,
		actionCreator ActionCreator,
		response io.ReadWriter,
	) *DeployResponse
}
