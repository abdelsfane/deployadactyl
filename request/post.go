package request

import (
	"bytes"
	"errors"
	"github.com/compozed/deployadactyl/interfaces"
)

type PushController interface {
	RunDeployment(postDeploymentRequest PostDeploymentRequest, response *bytes.Buffer) (deployResponse interfaces.DeployResponse)
}

type PostRequest struct {
	ArtifactUrl          string                 `json:"artifact_url"`
	Manifest             string                 `json:"manifest"`
	EnvironmentVariables map[string]string      `json:"environment_variables"`
	HealthCheckEndpoint  string                 `json:"health_check_endpoint"`
	Data                 map[string]interface{} `json:"data"`
	UUID                 string                 `json:"uuid"`
}

type PostDeploymentRequest struct {
	interfaces.Deployment
	Request PostRequest
}

func (r PostDeploymentRequest) GetId() string {
	return r.Request.UUID
}

func (r PostDeploymentRequest) GetData() map[string]interface{} {
	return r.Request.Data
}

func (r PostDeploymentRequest) GetContext() interfaces.CFContext {
	return r.Deployment.CFContext
}

func (r PostDeploymentRequest) SetContext(newCFContext interfaces.CFContext) interfaces.RequestDescriptor {
	r.Deployment.CFContext = newCFContext
	return r
}

func (r PostDeploymentRequest) GetAuthorization() interfaces.Authorization {
	return r.Deployment.Authorization
}

func (r PostDeploymentRequest) SetAuthorization(newAuthorization interfaces.Authorization) interfaces.RequestDescriptor {
	r.Deployment.Authorization = newAuthorization
	return r
}

func (r PostDeploymentRequest) GetRequest() interface{} {
	return r.Request
}

func (r PostDeploymentRequest) SetRequest(request interface{}) (interfaces.RequestDescriptor, error) {
	postRequest, ok := request.(PostRequest)
	if !ok {
		return nil, InvalidArgumentError{Err: errors.New("PostDeploymentRequest.SetRequest requires a PostRequest")}
	}

	r.Request = postRequest
	return r, nil
}
