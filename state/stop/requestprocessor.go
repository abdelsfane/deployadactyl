package stop

import (
	"bytes"

	"github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/request"
)

type StopRequestProcessorConstructor func(log interfaces.DeploymentLogger, controller request.StopController, request request.PutDeploymentRequest, buffer *bytes.Buffer) interfaces.RequestProcessor

func NewStopRequestProcessor(log interfaces.DeploymentLogger, sc request.StopController, request request.PutDeploymentRequest, buffer *bytes.Buffer) interfaces.RequestProcessor {
	return &StopRequestProcessor{
		StopController: sc,
		Request:        request,
		Response:       buffer,
		Log:            log,
	}
}

type StopRequestProcessor struct {
	StopController request.StopController
	Request        request.PutDeploymentRequest
	Response       *bytes.Buffer
	Log            interfaces.DeploymentLogger
}

func (c StopRequestProcessor) Process() interfaces.DeployResponse {
	return c.StopController.StopDeployment(c.Request, c.Response)
}
