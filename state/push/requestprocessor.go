package push

import (
	"bytes"

	"github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/request"
)

type PushRequestProcessorConstructor func(log interfaces.DeploymentLogger, controller request.PushController, request request.PostDeploymentRequest, buffer *bytes.Buffer) interfaces.RequestProcessor

func NewPushRequestProcessor(log interfaces.DeploymentLogger, pc request.PushController, request request.PostDeploymentRequest, buffer *bytes.Buffer) interfaces.RequestProcessor {
	return &PushRequestProcessor{
		PushController: pc,
		Request:        request,
		Response:       buffer,
		Log:            log,
	}
}

type PushRequestProcessor struct {
	PushController request.PushController
	Request        request.PostDeploymentRequest
	Response       *bytes.Buffer
	Log            interfaces.DeploymentLogger
}

func (c PushRequestProcessor) Process() interfaces.DeployResponse {
	return c.PushController.RunDeployment(c.Request, c.Response)
}
