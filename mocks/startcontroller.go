package mocks

import (
	"bytes"
	"github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/request"
)

type StartController struct {
	StartDeploymentCall struct {
		Received struct {
			Deployment request.PutDeploymentRequest
			Response   *bytes.Buffer
		}
		Returns struct {
			DeployResponse interfaces.DeployResponse
		}
		Writes string
		Called bool
	}
}

func (c *StartController) StartDeployment(deployment request.PutDeploymentRequest, response *bytes.Buffer) (deployResponse interfaces.DeployResponse) {
	c.StartDeploymentCall.Called = true
	c.StartDeploymentCall.Received.Deployment = deployment
	c.StartDeploymentCall.Received.Deployment.Request.Data = deployment.Request.Data
	c.StartDeploymentCall.Received.Response = response

	if c.StartDeploymentCall.Writes != "" {
		response.Write([]byte(c.StartDeploymentCall.Writes))
	}

	return c.StartDeploymentCall.Returns.DeployResponse
}
