package mocks

import (
	"bytes"
	"fmt"

	"github.com/gin-gonic/gin"

	I "github.com/compozed/deployadactyl/interfaces"
)

type Controller struct {
	RunDeploymentCall struct {
		Called   bool
		Received struct {
			Deployment *I.Deployment
			Response   *bytes.Buffer
		}
		Write struct {
			Output string
		}
		Returns I.DeployResponse
	}
	RunDeploymentViaHttpCall struct {
		Called   bool
		Received struct {
			Context *gin.Context
		}
	}
	PutRequestHandlerCall struct {
		Called   bool
		Received struct {
			Context *gin.Context
		}
	}
}

func (c *Controller) RunDeployment(deployment *I.Deployment, response *bytes.Buffer) I.DeployResponse {
	c.RunDeploymentCall.Called = true

	c.RunDeploymentCall.Received.Deployment = deployment
	c.RunDeploymentCall.Received.Response = response

	fmt.Fprint(response, c.RunDeploymentCall.Write.Output)

	return c.RunDeploymentCall.Returns
}

func (c *Controller) RunDeploymentViaHttp(g *gin.Context) {
	c.RunDeploymentViaHttpCall.Called = true

	c.RunDeploymentViaHttpCall.Received.Context = g
}

func (c *Controller) PutRequestHandler(g *gin.Context) {
	c.PutRequestHandlerCall.Called = true

	c.PutRequestHandlerCall.Received.Context = g
}
