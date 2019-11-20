package delete

import (
	"bytes"

	"github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/request"
)

type DeleteRequestProcessorConstructor func(log interfaces.DeploymentLogger, controller request.DeleteController, request request.DeleteDeploymentRequest, buffer *bytes.Buffer) interfaces.RequestProcessor

func NewDeleteRequestProcessor(log interfaces.DeploymentLogger, sc request.DeleteController, request request.DeleteDeploymentRequest, buffer *bytes.Buffer) interfaces.RequestProcessor {
	return &DeleteRequestProcessor{
		DeleteController: sc,
		Request:          request,
		Response:         buffer,
		Log:              log,
	}
}

type DeleteRequestProcessor struct {
	DeleteController request.DeleteController
	Request          request.DeleteDeploymentRequest
	Response         *bytes.Buffer
	Log              interfaces.DeploymentLogger
}

func (c DeleteRequestProcessor) Process() interfaces.DeployResponse {
	return c.DeleteController.DeleteDeployment(c.Request, c.Response)
}
