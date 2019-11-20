package request

import (
	"bytes"
	"errors"
	"github.com/compozed/deployadactyl/interfaces"
)

type StartController interface {
	StartDeployment(request PutDeploymentRequest, response *bytes.Buffer) (deployResponse interfaces.DeployResponse)
}

type StopController interface {
	StopDeployment(request PutDeploymentRequest, response *bytes.Buffer) (deployResponse interfaces.DeployResponse)
}

type PutRequest struct {
	State string                 `json:"state"`
	Data  map[string]interface{} `json:"data"`
	UUID  string                 `json:"uuid"`
}

type PutDeploymentRequest struct {
	interfaces.Deployment
	Request PutRequest
}

func (r PutDeploymentRequest) GetId() string {
	return r.Request.UUID
}

func (r PutDeploymentRequest) GetData() map[string]interface{} {
	return r.Request.Data
}

func (r PutDeploymentRequest) GetContext() interfaces.CFContext {
	return r.Deployment.CFContext
}

func (r PutDeploymentRequest) SetContext(newCFContext interfaces.CFContext) interfaces.RequestDescriptor {
	r.Deployment.CFContext = newCFContext
	return r
}

func (r PutDeploymentRequest) GetAuthorization() interfaces.Authorization {
	return r.Deployment.Authorization
}

func (r PutDeploymentRequest) SetAuthorization(newAuthorization interfaces.Authorization) interfaces.RequestDescriptor {
	r.Deployment.Authorization = newAuthorization
	return r
}

func (r PutDeploymentRequest) GetRequest() interface{} {
	return r.Request
}

func (r PutDeploymentRequest) SetRequest(request interface{}) (interfaces.RequestDescriptor, error) {
	putRequest, ok := request.(PutRequest)
	if !ok {
		return nil, InvalidArgumentError{Err: errors.New("PutDeploymentRequest.SetRequest requires a PutRequest")}
	}

	r.Request = putRequest
	return r, nil
}
