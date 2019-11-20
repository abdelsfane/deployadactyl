package start

import (
	"bytes"

	"github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/request"
)

type StartRequestProcessorConstructor func(log interfaces.DeploymentLogger, controller request.StartController, request request.PutDeploymentRequest, buffer *bytes.Buffer) interfaces.RequestProcessor

func NewStartRequestProcessor(log interfaces.DeploymentLogger, sc request.StartController, request request.PutDeploymentRequest, buffer *bytes.Buffer) interfaces.RequestProcessor {
	return &StartRequestProcessor{
		StartController: sc,
		Request:         request,
		Response:        buffer,
		Log:             log,
	}
}

type StartRequestProcessor struct {
	StartController request.StartController
	Request         request.PutDeploymentRequest
	Response        *bytes.Buffer
	Log             interfaces.DeploymentLogger
}

func (c StartRequestProcessor) Process() interfaces.DeployResponse {
	return c.StartController.StartDeployment(c.Request, c.Response)
}
