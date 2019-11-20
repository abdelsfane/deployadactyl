package push

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/compozed/deployadactyl/constants"
	"github.com/compozed/deployadactyl/controller/deployer"
	"github.com/compozed/deployadactyl/controller/deployer/bluegreen"
	"github.com/compozed/deployadactyl/controller/deployer/manifestro"
	H "github.com/compozed/deployadactyl/eventmanager/handlers/healthchecker"
	R "github.com/compozed/deployadactyl/eventmanager/handlers/routemapper"
	I "github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/state"
	S "github.com/compozed/deployadactyl/structs"
)

const deploymentOutput = `Deployment Parameters:
Artifact URL: %s,
Username:     %s,
Environment:  %s,
Org:          %s,
Space:        %s,
AppName:      %s`

const successfulDeploy = `Your deploy was successful! (^_^)b
If you experience any problems after this point, check that you can manually push your application to Cloud Foundry on a lower environment.
It is likely that it is an error with your application and not with Deployadactyl.
Thanks for using Deployadactyl! Please push down pull up on your lap bar and exit to your left.

`

type PushManagerConstructor func(courierCreator I.CourierCreator, eventManager I.EventManager, log I.DeploymentLogger, fetcher I.Fetcher, deployEventData S.DeployEventData, fileSystemCleaner FileSystemCleaner, cfContext I.CFContext, auth I.Authorization, environment S.Environment, envVars map[string]string, healthChecker H.HealthChecker, routeMapper R.RouteMapper) I.ActionCreator

func NewPushManager(c I.CourierCreator, em I.EventManager, log I.DeploymentLogger, f I.Fetcher, ded S.DeployEventData, fcs FileSystemCleaner, cf I.CFContext, auth I.Authorization, env S.Environment, envVars map[string]string, healthChecker H.HealthChecker, routeMapper R.RouteMapper) I.ActionCreator {
	return &PushManager{
		CourierCreator:       c,
		EventManager:         em,
		Logger:               log,
		Fetcher:              f,
		DeployEventData:      ded,
		FileSystemCleaner:    fcs,
		CFContext:            cf,
		Auth:                 auth,
		Environment:          env,
		EnvironmentVariables: envVars,
		HealthChecker:        healthChecker,
		RouteMapper:          routeMapper,
	}
}

type FileSystemCleaner interface {
	RemoveAll(path string) error
}

type PushManager struct {
	CourierCreator       I.CourierCreator
	EventManager         I.EventManager
	Logger               I.DeploymentLogger
	Fetcher              I.Fetcher
	DeployEventData      S.DeployEventData
	FileSystemCleaner    FileSystemCleaner
	CFContext            I.CFContext
	Auth                 I.Authorization
	Environment          S.Environment
	EnvironmentVariables map[string]string
	HealthChecker        H.HealthChecker
	RouteMapper          R.RouteMapper
}

func (a *PushManager) SetUp() error {
	var (
		manifestString string
		instances      *uint16
		appPath        string
		err            error
	)

	var fetchFn func() (string, error)

	if a.DeployEventData.DeploymentInfo.ContentType == "application/json" {

		if a.DeployEventData.DeploymentInfo.Manifest != "" {
			manifest, err := base64.StdEncoding.DecodeString(a.DeployEventData.DeploymentInfo.Manifest)
			if err != nil {
				return state.ManifestError{}
			}
			manifestString = string(manifest)
		}

		fetchFn = func() (string, error) {
			a.Logger.Debug("deploying from json request")
			appPath, err = a.Fetcher.Fetch(a.DeployEventData.DeploymentInfo.ArtifactURL, manifestString)
			if err != nil {
				return "", state.AppPathError{Err: err}
			}
			return appPath, nil
		}
	} else {
		fetchFn = func() (string, error) {
			a.Logger.Debug("deploying from archive request")

			appPath, manifestString, err = a.Fetcher.FetchArtifactFromRequest(a.DeployEventData.RequestBody, a.DeployEventData.DeploymentInfo.ContentType)
			if err != nil {
				return "", state.UnzippingError{Err: err}
			}

			return appPath, nil
		}
	}

	var event I.IEvent
	event = ArtifactRetrievalStartEvent{
		CFContext:   a.CFContext,
		Auth:        a.Auth,
		Environment: a.Environment,
		Response:    a.DeployEventData.Response,
		Data:        a.DeployEventData.DeploymentInfo.Data,
		Manifest:    manifestString,
		ArtifactURL: a.DeployEventData.DeploymentInfo.ArtifactURL,
		Log:         a.Logger,
	}
	a.Logger.Debugf("emitting a %s event", event.Name())

	err = a.EventManager.EmitEvent(event)
	if err != nil {
		a.Logger.Error(err)
		err = &bluegreen.InitializationError{err}
		return deployer.EventError{Type: event.Name(), Err: err}
	}

	appPath, err = fetchFn()

	instances = manifestro.GetInstances(manifestString)
	if instances == nil {
		instances = &a.Environment.Instances
	}

	if err != nil {
		a.Logger.Error(err)
		event = ArtifactRetrievalFailureEvent{
			CFContext:   a.CFContext,
			Auth:        a.Auth,
			Environment: a.Environment,
			Response:    a.DeployEventData.Response,
			Data:        a.DeployEventData.DeploymentInfo.Data,
			Manifest:    manifestString,
			ArtifactURL: a.DeployEventData.DeploymentInfo.ArtifactURL,
			Log:         a.Logger,
		}
		a.EventManager.EmitEvent(event)
		return err
	}

	event = ArtifactRetrievalSuccessEvent{
		CFContext:            a.CFContext,
		Auth:                 a.Auth,
		Environment:          a.Environment,
		Response:             a.DeployEventData.Response,
		Data:                 a.DeployEventData.DeploymentInfo.Data,
		Manifest:             manifestString,
		ArtifactURL:          a.DeployEventData.DeploymentInfo.ArtifactURL,
		AppPath:              appPath,
		EnvironmentVariables: a.EnvironmentVariables,
		Log:                  a.Logger,
	}
	a.Logger.Debugf("emitting a %s event", event.Name())
	err = a.EventManager.EmitEvent(event)
	if err != nil {
		a.Logger.Error(err)
		err = &bluegreen.InitializationError{err}
		return deployer.EventError{Type: event.Name(), Err: err}
	}

	a.DeployEventData.DeploymentInfo.Manifest = manifestString
	a.DeployEventData.DeploymentInfo.AppPath = appPath
	a.DeployEventData.DeploymentInfo.Instances = *instances

	return nil
}

func (a PushManager) OnStart() error {
	info := a.DeployEventData.DeploymentInfo
	deploymentMessage := fmt.Sprintf(deploymentOutput, info.ArtifactURL, info.Username, info.Environment, info.Org, info.Space, info.AppName)

	a.Logger.Info(deploymentMessage)
	fmt.Fprintln(a.DeployEventData.Response, deploymentMessage)

	err := a.EventManager.Emit(I.Event{Type: constants.PushStartedEvent, Data: &a.DeployEventData})
	if err != nil {
		a.Logger.Error(err)
		err = &bluegreen.InitializationError{err}
		return deployer.EventError{Type: constants.PushStartedEvent, Err: err}
	}

	event := PushStartedEvent{
		CFContext:   a.CFContext,
		Auth:        a.Auth,
		Environment: a.Environment,
		Body:        info.Body,
		Response:    a.DeployEventData.Response,
		ContentType: info.ContentType,
		Data:        info.Data,
		Instances:   info.Instances,
		Log:         a.Logger,
	}
	err = a.EventManager.EmitEvent(event)
	if err != nil {
		a.Logger.Error(err)
		err = &bluegreen.InitializationError{err}
		return deployer.EventError{Type: event.Name(), Err: err}
	}
	return nil
}

func (a PushManager) OnFinish(env S.Environment, response io.ReadWriter, err error) I.DeployResponse {
	if err != nil {
		if env.DisableRollback {
			a.Logger.Errorf("DisabledRollback %t, returning status %d and err %s", env.DisableRollback, http.StatusOK, err)
			return I.DeployResponse{
				StatusCode: http.StatusOK,
				Error:      err,
			}
		}

		if matched, _ := regexp.MatchString("login failed", err.Error()); matched {
			return I.DeployResponse{
				StatusCode: http.StatusBadRequest,
				Error:      err,
			}
		}

		return I.DeployResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      err,
		}
	}
	a.Logger.Infof("successfully deployed application %s", a.DeployEventData.DeploymentInfo.AppName)
	fmt.Fprintf(response, "\n%s", successfulDeploy)

	return I.DeployResponse{StatusCode: http.StatusOK}
}

func (a PushManager) CleanUp() {
	a.FileSystemCleaner.RemoveAll(a.DeployEventData.DeploymentInfo.AppPath)
}

func (a PushManager) Create(environment S.Environment, response io.ReadWriter, foundationURL string) (I.Action, error) {

	courier, err := a.CourierCreator.CreateCourier()
	if err != nil {
		a.Logger.Error(err)
		return &Pusher{}, state.CourierCreationError{Err: err}
	}

	p := &Pusher{
		Courier:        courier,
		DeploymentInfo: *a.DeployEventData.DeploymentInfo,
		EventManager:   a.EventManager,
		Response:       response,
		Log:            a.Logger,
		FoundationURL:  foundationURL,
		AppPath:        a.DeployEventData.DeploymentInfo.AppPath,
		Environment:    environment,
		Fetcher:        a.Fetcher,
		CFContext:      a.CFContext,
		Auth:           a.Auth,
		HealthChecker:  a.HealthChecker,
		RouteMapper:    a.RouteMapper,
	}

	return p, nil
}

func (a PushManager) InitiallyError(initiallyErrors []error) error {
	return bluegreen.LoginError{LoginErrors: initiallyErrors}
}

func (a PushManager) ExecuteError(executeErrors []error) error {
	return bluegreen.PushError{PushErrors: executeErrors}
}

func (a PushManager) UndoError(executeErrors, undoErrors []error) error {
	return bluegreen.RollbackError{PushErrors: executeErrors, RollbackErrors: undoErrors}
}

func (a PushManager) SuccessError(successErrors []error) error {
	return bluegreen.FinishPushError{FinishPushError: successErrors}
}
