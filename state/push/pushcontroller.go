package push

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/compozed/deployadactyl/constants"
	"github.com/compozed/deployadactyl/controller/deployer"
	"github.com/compozed/deployadactyl/controller/deployer/bluegreen"
	I "github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/request"
	"github.com/compozed/deployadactyl/structs"
	"github.com/go-errors/errors"
)

type PushControllerConstructor func(log I.DeploymentLogger, deployer, silentDeployer I.Deployer, eventManager I.EventManager, errorFinder I.ErrorFinder, pushManagerFactory I.PushManagerFactory, resolver I.AuthResolver, envResolver I.EnvResolver) request.PushController

func NewPushController(l I.DeploymentLogger, d, sd I.Deployer, em I.EventManager, ef I.ErrorFinder, pmf I.PushManagerFactory, resolver I.AuthResolver, envResolver I.EnvResolver) request.PushController {
	return &PushController{
		Deployer:           d,
		SilentDeployer:     sd,
		EventManager:       em,
		ErrorFinder:        ef,
		PushManagerFactory: pmf,
		Log:                l,
		AuthResolver:       resolver,
		EnvResolver:        envResolver,
	}
}

type MissingParameterError struct {
	Err error
}

func (e MissingParameterError) Error() string {
	return e.Err.Error()
}

type PushController struct {
	Deployer           I.Deployer
	SilentDeployer     I.Deployer
	Log                I.DeploymentLogger
	EventManager       I.EventManager
	ErrorFinder        I.ErrorFinder
	PushManagerFactory I.PushManagerFactory
	AuthResolver       I.AuthResolver
	EnvResolver        I.EnvResolver
}

// PUSH specific
func (c *PushController) RunDeployment(deployment request.PostDeploymentRequest, response *bytes.Buffer) (deployResponse I.DeployResponse) {
	cf := deployment.CFContext

	if deployment.Type == "application/json" && deployment.Request.ArtifactUrl == "" {
		c.Log.Error("artifact url is missing from request")
		return I.DeployResponse{
			StatusCode: http.StatusBadRequest,
			Error:      MissingParameterError{Err: errors.New("the following properties are missing: artifact_url")},
		}
	}

	if deployment.Request.Data == nil {
		deployment.Request.Data = make(map[string]interface{})
	}

	deploymentInfo := &structs.DeploymentInfo{
		Org:                  cf.Organization,
		Space:                cf.Space,
		AppName:              cf.Application,
		Environment:          cf.Environment,
		UUID:                 c.Log.UUID,
		Manifest:             deployment.Request.Manifest,
		ArtifactURL:          deployment.Request.ArtifactUrl,
		EnvironmentVariables: deployment.Request.EnvironmentVariables,
		HealthCheckEndpoint:  deployment.Request.HealthCheckEndpoint,
		Data:                 deployment.Request.Data,
	}

	c.Log.Debugf("Starting deploy of %s with UUID %s", cf.Application, deploymentInfo.UUID)
	c.Log.Debug("building deploymentInfo")

	body := ioutil.NopCloser(bytes.NewBuffer(*deployment.Body))
	if deployment.Type == "application/json" || deployment.Type == "application/zip" || deployment.Type == "application/x-tar" || deployment.Type == "application/x-gzip" {
		deploymentInfo.ContentType = deployment.Type
	} else {
		return I.DeployResponse{
			StatusCode: http.StatusBadRequest,
			Error:      deployer.InvalidContentTypeError{},
		}
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

	deploymentInfo.Username = auth.Username
	deploymentInfo.Password = auth.Password
	deploymentInfo.Domain = environment.Domain
	deploymentInfo.SkipSSL = environment.SkipSSL
	deploymentInfo.CustomParams = environment.CustomParams
	deploymentInfo.Data = deployment.Request.Data
	deploymentInfo.Body = body

	deployEventData := structs.DeployEventData{Response: response, DeploymentInfo: deploymentInfo, RequestBody: body}
	defer c.emitDeployFinish(&deployEventData, response, cf, auth, environment, &deployResponse, c.Log)
	defer c.emitDeploySuccessOrFailure(&deployEventData, response, cf, auth, environment, &deployResponse, c.Log)

	c.Log.Debugf("emitting a %s event", constants.DeployStartEvent)

	err = c.EventManager.Emit(I.Event{Type: constants.DeployStartEvent, Data: &deployEventData})
	if err != nil {
		c.Log.Error(err)
		err = &bluegreen.InitializationError{err}
		return I.DeployResponse{
			StatusCode:     http.StatusInternalServerError,
			Error:          deployer.EventError{Type: constants.DeployStartEvent, Err: err},
			DeploymentInfo: deploymentInfo,
		}
	}

	err = c.EventManager.EmitEvent(DeployStartedEvent{
		CFContext:   cf,
		Auth:        auth,
		Body:        body,
		ContentType: deploymentInfo.ContentType,
		Environment: environment,
		Response:    response,
		ArtifactURL: deploymentInfo.ArtifactURL,
		Data:        deploymentInfo.Data,
		Log:         c.Log,
	})
	if err != nil {
		c.Log.Error(err)
		err = &bluegreen.InitializationError{err}
		return I.DeployResponse{
			StatusCode:     http.StatusInternalServerError,
			Error:          deployer.EventError{Type: constants.DeployStartEvent, Err: err},
			DeploymentInfo: deploymentInfo,
		}
	}

	pusherCreator := c.PushManagerFactory.PushManager(deployEventData, auth, environment, deploymentInfo.EnvironmentVariables)

	reqChannel1 := make(chan *I.DeployResponse)
	reqChannel2 := make(chan *I.DeployResponse)
	defer close(reqChannel1)
	defer close(reqChannel2)

	go func() {
		reqChannel1 <- c.Deployer.Deploy(deploymentInfo, environment, pusherCreator, response)
	}()

	silentResponse := &bytes.Buffer{}
	if cf.Environment == os.Getenv("SILENT_DEPLOY_ENVIRONMENT") {
		go func() {
			reqChannel2 <- c.SilentDeployer.Deploy(deploymentInfo, environment, pusherCreator, silentResponse)
		}()
		<-reqChannel2
	}

	deployResponse = *<-reqChannel1

	return deployResponse
}

func (c *PushController) emitDeployFinish(deployEventData *structs.DeployEventData, response io.ReadWriter, cf I.CFContext, auth I.Authorization, environment structs.Environment, deployResponse *I.DeployResponse, deploymentLogger I.DeploymentLogger) {
	deploymentLogger.Debugf("emitting a %s event", constants.DeployFinishEvent)
	finishErr := c.EventManager.Emit(I.Event{Type: constants.DeployFinishEvent, Data: deployEventData})
	if finishErr != nil {
		fmt.Fprintln(response, finishErr)
		err := bluegreen.FinishDeployError{Err: fmt.Errorf("%s: %s", deployResponse.Error, deployer.EventError{constants.DeployFinishEvent, finishErr})}
		deployResponse.Error = err
		deployResponse.StatusCode = http.StatusInternalServerError
	}

	finishErr = c.EventManager.EmitEvent(DeployFinishedEvent{
		CFContext:   cf,
		Auth:        auth,
		Body:        deployEventData.RequestBody,
		ContentType: deployEventData.DeploymentInfo.ContentType,
		Environment: environment,
		Response:    deployEventData.Response,
		Data:        deployEventData.DeploymentInfo.Data,
		Log:         c.Log,
	})
	if finishErr != nil {
		fmt.Fprintln(response, finishErr)
		if finishErr != nil {
			fmt.Fprintln(response, finishErr)
			err := bluegreen.FinishDeployError{Err: fmt.Errorf("%s: %s", deployResponse.Error, deployer.EventError{constants.DeployFinishEvent, finishErr})}
			deployResponse.Error = err
			deployResponse.StatusCode = http.StatusInternalServerError
		}
	}
}

func (c PushController) emitDeploySuccessOrFailure(deployEventData *structs.DeployEventData, response io.ReadWriter, cf I.CFContext, auth I.Authorization, environment structs.Environment, deployResponse *I.DeployResponse, deploymentLogger I.DeploymentLogger) {
	deployEvent := I.Event{Type: constants.DeploySuccessEvent, Data: deployEventData}
	if deployResponse.Error != nil {
		c.printErrors(response, &deployResponse.Error)

		deployEvent.Type = constants.DeployFailureEvent
		deployEvent.Error = deployResponse.Error
	}
	deploymentLogger.Debug(fmt.Sprintf("emitting a %s event", deployEvent.Name()))
	eventErr := c.EventManager.Emit(deployEvent)
	if eventErr != nil {
		deploymentLogger.Errorf("an error occurred when emitting a %s event: %s", deployEvent.Name(), eventErr)
		fmt.Fprintln(response, eventErr)
		return
	}

	var event I.IEvent
	if deployResponse.Error != nil {
		event = DeployFailureEvent{
			CFContext:   cf,
			Auth:        auth,
			Body:        deployEventData.RequestBody,
			ContentType: deployEventData.DeploymentInfo.ContentType,
			Environment: environment,
			Response:    deployEventData.Response,
			Data:        deployEventData.DeploymentInfo.Data,
			Error:       deployResponse.Error,
			Log:         c.Log,
		}
	} else {
		event = DeploySuccessEvent{
			CFContext:           cf,
			Auth:                auth,
			Body:                deployEventData.RequestBody,
			ContentType:         deployEventData.DeploymentInfo.ContentType,
			Environment:         environment,
			Response:            deployEventData.Response,
			Data:                deployEventData.DeploymentInfo.Data,
			HealthCheckEndpoint: deployEventData.DeploymentInfo.HealthCheckEndpoint,
			ArtifactURL:         deployEventData.DeploymentInfo.ArtifactURL,
			Log:                 c.Log,
		}
	}
	deploymentLogger.Debug(fmt.Sprintf("emitting a %s event", event.Name()))
	eventErr = c.EventManager.EmitEvent(event)
	if eventErr != nil {
		deploymentLogger.Errorf("an error occurred when emitting a %s event: %s", event.Name(), eventErr)
		fmt.Fprintln(response, eventErr)
	}

}

func (c PushController) printErrors(response io.ReadWriter, err *error) {
	tempBuffer := bytes.Buffer{}
	tempBuffer.ReadFrom(response)
	fmt.Fprint(response, tempBuffer.String())

	errors := c.ErrorFinder.FindErrors(tempBuffer.String())
	fmt.Fprintln(response)
	fmt.Fprintln(response, "<conveyor-error>")
	fmt.Fprintln(response, "********** Deployment Failure Detected **********")
	if len(errors) > 0 {
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
	} else {
		c.Log.Info("Unknown Error in Cloud Foundry logs")
		fmt.Fprintln(response, "****")
		fmt.Fprintln(response)
		fmt.Fprintln(response, "Error: Your application failed to start")
		fmt.Fprintln(response)
		fmt.Fprintln(response, "Logs can be found in your console output above")
		fmt.Fprintln(response)
		fmt.Fprintln(response, "****")
	}
	fmt.Fprintln(response, "*************************************************")
	fmt.Fprintln(response, "</conveyor-error>")
}
