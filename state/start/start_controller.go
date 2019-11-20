package start

import (
	"bytes"
	"fmt"
	"net/http"

	"io"

	"github.com/compozed/deployadactyl/controller/deployer"
	"github.com/compozed/deployadactyl/controller/deployer/bluegreen"
	I "github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/request"
	"github.com/compozed/deployadactyl/structs"
)

type StartControllerConstructor func(log I.DeploymentLogger, deployer I.Deployer, eventManager I.EventManager, errorFinder I.ErrorFinder, startManagerFactory I.StartManagerFactory, resolver I.AuthResolver, envResolver I.EnvResolver) request.StartController

func NewStartController(l I.DeploymentLogger, d I.Deployer, em I.EventManager, ef I.ErrorFinder, smf I.StartManagerFactory, resolver I.AuthResolver, envResolver I.EnvResolver) request.StartController {
	return &StartController{
		Deployer:            d,
		EventManager:        em,
		ErrorFinder:         ef,
		StartManagerFactory: smf,
		Log:                 l,
		AuthResolver:        resolver,
		EnvResolver:         envResolver,
	}
}

// StartController is used to determine the type of request and process it accordingly.
type StartController struct {
	Log                 I.DeploymentLogger
	StartManagerFactory I.StartManagerFactory
	Deployer            I.Deployer
	EventManager        I.EventManager
	ErrorFinder         I.ErrorFinder
	AuthResolver        I.AuthResolver
	EnvResolver         I.EnvResolver
}

//deployment *I.Deployment, data map[string]interface{}

func (c *StartController) StartDeployment(deployment request.PutDeploymentRequest, response *bytes.Buffer) (deployResponse I.DeployResponse) {
	cf := deployment.CFContext
	c.Log.Debugf("Preparing to start %s with UUID %s", cf.Application, c.Log.UUID)

	if deployment.Request.Data == nil {
		deployment.Request.Data = make(map[string]interface{})
	}

	environment, err := c.EnvResolver.Resolve(cf.Environment)
	if err != nil {
		fmt.Fprintln(response, err.Error())
		return I.DeployResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      err,
		}
	}
	auth, err := c.AuthResolver.Resolve(deployment.Authorization, environment, c.Log)
	if err != nil {
		return I.DeployResponse{
			StatusCode: http.StatusUnauthorized,
			Error:      err,
		}
	}

	deploymentInfo := &structs.DeploymentInfo{
		Org:          cf.Organization,
		Space:        cf.Space,
		AppName:      cf.Application,
		Environment:  cf.Environment,
		UUID:         c.Log.UUID,
		Domain:       environment.Domain,
		SkipSSL:      environment.SkipSSL,
		CustomParams: environment.CustomParams,
		Username:     auth.Username,
		Password:     auth.Password,
		Data:         deployment.Request.Data,
	}

	defer c.emitStartFinish(response, c.Log, cf, &auth, &environment, deployment.Request.Data, &deployResponse)
	defer c.emitStartSuccessOrFailure(response, c.Log, cf, &auth, &environment, deployment.Request.Data, &deployResponse)

	err = c.EventManager.EmitEvent(StartStartedEvent{
		CFContext:     cf,
		Authorization: auth,
		Environment:   environment,
		Data:          deployment.Request.Data,
		Response:      response,
		Log:           c.Log,
	})
	if err != nil {
		c.Log.Error(err)
		err = &bluegreen.InitializationError{err}
		return I.DeployResponse{
			StatusCode:     http.StatusInternalServerError,
			Error:          deployer.EventError{Type: "StartStartedEvent", Err: err},
			DeploymentInfo: deploymentInfo,
		}
	}

	deployEventData := structs.DeployEventData{Response: response, DeploymentInfo: deploymentInfo}

	manager := c.StartManagerFactory.StartManager(deployEventData)
	deployResponse = *c.Deployer.Deploy(deploymentInfo, environment, manager, response)
	return deployResponse
}

func (c StartController) emitStartFinish(response io.ReadWriter, deploymentLogger I.DeploymentLogger, cfContext I.CFContext, auth *I.Authorization, environment *structs.Environment, data map[string]interface{}, deployResponse *I.DeployResponse) {
	var event I.IEvent
	event = StartFinishedEvent{
		CFContext:     cfContext,
		Authorization: *auth,
		Data:          data,
		Environment:   *environment,
		Log:           deploymentLogger,
	}
	deploymentLogger.Debugf("emitting a %s event", event.Name())
	c.EventManager.EmitEvent(event)
}

func (c StartController) emitStartSuccessOrFailure(response io.ReadWriter, deploymentLogger I.DeploymentLogger, cfContext I.CFContext, auth *I.Authorization, environment *structs.Environment, data map[string]interface{}, deployResponse *I.DeployResponse) {
	var event I.IEvent

	if deployResponse.Error != nil {
		c.printErrors(response, &deployResponse.Error)
		event = StartFailureEvent{
			CFContext:     cfContext,
			Authorization: *auth,
			Environment:   *environment,
			Data:          data,
			Response:      response,
			Error:         deployResponse.Error,
			Log:           deploymentLogger,
		}

	} else {
		event = StartSuccessEvent{
			CFContext:     cfContext,
			Authorization: *auth,
			Environment:   *environment,
			Data:          data,
			Response:      response,
			Log:           deploymentLogger,
		}
	}
	deploymentLogger.Debugf("emitting a %s event", event.Name())
	eventErr := c.EventManager.EmitEvent(event)
	if eventErr != nil {
		deploymentLogger.Errorf("an error occurred when emitting a %s event: %s", event.Name(), eventErr)
		fmt.Fprintln(response, eventErr)
	}
}

func (c StartController) printErrors(response io.ReadWriter, err *error) {
	tempBuffer := bytes.Buffer{}
	tempBuffer.ReadFrom(response)
	fmt.Fprint(response, tempBuffer.String())

	errors := c.ErrorFinder.FindErrors(tempBuffer.String())
	if len(errors) > 0 {
		fmt.Fprintln(response)
		fmt.Fprintln(response, "<conveyor-error>")
		fmt.Fprintln(response, "********** Deployment Failure Detected **********")
		*err = errors[0]
		for _, error := range errors {
			fmt.Fprintln(response, "****")
			fmt.Fprintln(response)
			fmt.Fprintln(response, "The following error was found in the above logs: "+error.Error())
			fmt.Fprintln(response)
			fmt.Fprintln(response, "Error: "+error.Details()[0])
			fmt.Fprintln(response)
			fmt.Fprintln(response, "Potential solution: "+error.Solution())
			fmt.Fprintln(response)
			fmt.Fprintln(response, "****")
		}

		fmt.Fprintln(response, "*************************************************")
		fmt.Fprintln(response, "</conveyor-error>")
	}
}
