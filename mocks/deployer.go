package mocks

import (
	"fmt"
	"io"

	I "github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/structs"
)

// Deployer handmade mock for tests.
type Deployer struct {
	DeployCall struct {
		Called   int
		Received struct {
			DeploymentInfo *structs.DeploymentInfo
			Env            structs.Environment
			ActionCreator  I.ActionCreator
			Response       io.ReadWriter
		}
		Write struct {
			Output string
		}
		Returns struct {
			Error      error
			StatusCode int
		}
	}
}

// Deploy mock method.
func (d *Deployer) Deploy(deploymentInfo *structs.DeploymentInfo, env structs.Environment, actionCreator I.ActionCreator, out io.ReadWriter) *I.DeployResponse {
	d.DeployCall.Called++

	d.DeployCall.Received.DeploymentInfo = deploymentInfo
	d.DeployCall.Received.Env = env
	d.DeployCall.Received.ActionCreator = actionCreator

	d.DeployCall.Received.Response = out

	fmt.Fprint(out, d.DeployCall.Write.Output)

	response := &I.DeployResponse{
		StatusCode:     d.DeployCall.Returns.StatusCode,
		Error:          d.DeployCall.Returns.Error,
		DeploymentInfo: deploymentInfo,
	}

	return response
}
