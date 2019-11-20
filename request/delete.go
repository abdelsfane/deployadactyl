package request

import (
	"bytes"
	"github.com/compozed/deployadactyl/interfaces"
	"github.com/go-errors/errors"
)

type DeleteController interface {
	DeleteDeployment(request DeleteDeploymentRequest, response *bytes.Buffer) (deployResponse interfaces.DeployResponse)
}

type DeleteRequest struct {
	State string                 `json:"state"`
	Data  map[string]interface{} `json:"data"`
	UUID  string                 `json:"uuid"`
}

type DeleteDeploymentRequest struct {
	interfaces.Deployment
	Request DeleteRequest
}

func (r DeleteDeploymentRequest) GetId() string {
	return r.Request.UUID
}

func (r DeleteDeploymentRequest) GetData() map[string]interface{} {
	return r.Request.Data
}

func (r DeleteDeploymentRequest) GetContext() interfaces.CFContext {
	return r.Deployment.CFContext
}

func (r DeleteDeploymentRequest) SetContext(newCFContext interfaces.CFContext) interfaces.RequestDescriptor {
	r.Deployment.CFContext = newCFContext
	return r
}

func (r DeleteDeploymentRequest) GetAuthorization() interfaces.Authorization {
	return r.Deployment.Authorization
}

func (r DeleteDeploymentRequest) SetAuthorization(newAuthorization interfaces.Authorization) interfaces.RequestDescriptor {
	r.Deployment.Authorization = newAuthorization
	return r
}

func (r DeleteDeploymentRequest) GetRequest() interface{} {
	return r.Request
}

func (r DeleteDeploymentRequest) SetRequest(request interface{}) (interfaces.RequestDescriptor, error) {
	deleteRequest, ok := request.(DeleteRequest)
	if !ok {
		return nil, InvalidArgumentError{Err: errors.New("DeleteDeploymentRequest.SetRequest requires a DeleteRequest")}
	}

	r.Request = deleteRequest
	return r, nil
}
