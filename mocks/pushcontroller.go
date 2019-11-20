package mocks

import (
	"bytes"

	"github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/request"
)

type PushController struct {
	RunDeploymentCall struct {
		Received struct {
			Request  request.PostDeploymentRequest
			Response *bytes.Buffer
		}
		Returns struct {
			DeployResponse interfaces.DeployResponse
		}
		Writes string
		Called bool
	}
}

func (c *PushController) RunDeployment(deployment request.PostDeploymentRequest, response *bytes.Buffer) (deployResponse interfaces.DeployResponse) {
	c.RunDeploymentCall.Called = true
	c.RunDeploymentCall.Received.Request = deployment
	c.RunDeploymentCall.Received.Response = response

	if c.RunDeploymentCall.Writes != "" {
		response.Write([]byte(c.RunDeploymentCall.Writes))
	}

	return c.RunDeploymentCall.Returns.DeployResponse
}
