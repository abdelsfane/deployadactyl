package creator

import (
	"bytes"

	"github.com/compozed/deployadactyl/artifetcher"
	"github.com/compozed/deployadactyl/artifetcher/extractor"
	"github.com/compozed/deployadactyl/controller/deployer"
	"github.com/compozed/deployadactyl/controller/deployer/bluegreen"
	"github.com/compozed/deployadactyl/controller/deployer/prechecker"
	"github.com/compozed/deployadactyl/eventmanager"
	I "github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/request"
	"github.com/compozed/deployadactyl/state/delete"
	"github.com/compozed/deployadactyl/state/push"
	"github.com/compozed/deployadactyl/state/start"
	"github.com/compozed/deployadactyl/state/stop"
	"github.com/compozed/deployadactyl/structs"
)

func newRequestCreator(c Creator, uuid string, b *bytes.Buffer) RequestCreator {
	logger := I.DeploymentLogger{UUID: uuid, Log: c.GetLogger()}
	var em I.EventManager
	if c.provider.NewEventManager != nil {
		em = c.provider.NewEventManager(logger, c.GetEventBindings().GetBindings())
	} else {
		em = eventmanager.NewEventManager(logger, c.GetEventBindings().GetBindings())
	}
	return RequestCreator{
		Creator:      c,
		EventManager: em,
		Buffer:       b,
		Log:          logger,
	}
}

type RequestCreator struct {
	Creator
	EventManager I.EventManager
	Buffer       *bytes.Buffer
	Log          I.DeploymentLogger
}

func (r *RequestCreator) CreateDeployer() I.Deployer {
	if r.provider.NewDeployer != nil {
		return r.provider.NewDeployer(r.CreateConfig(), r.CreateBlueGreener(), r.createPrechecker(), r.CreateEventManager(), r.createRandomizer(), r.createErrorFinder(), r.Log)
	}
	return deployer.NewDeployer(r.CreateConfig(), r.CreateBlueGreener(), r.createPrechecker(), r.CreateEventManager(), r.createRandomizer(), r.createErrorFinder(), r.Log)
}

func (r RequestCreator) CreateBlueGreener() I.BlueGreener {
	if r.provider.NewBlueGreen != nil {
		return r.provider.NewBlueGreen(r.Log)
	}
	return bluegreen.NewBlueGreen(r.Log)
}

func (r RequestCreator) CreateFetcher() I.Fetcher {
	if r.provider.NewFetcher != nil {
		return r.provider.NewFetcher(r.CreateFileSystem(), r.CreateExtractor(), r.Log)
	}
	return artifetcher.NewArtifetcher(r.CreateFileSystem(), r.CreateExtractor(), r.Log)
}

func (r RequestCreator) CreateExtractor() I.Extractor {
	if r.provider.NewExtractor != nil {
		return r.provider.NewExtractor(r.Log, r.CreateFileSystem())
	}
	return extractor.NewExtractor(r.Log, r.CreateFileSystem())
}

func (r *RequestCreator) CreateEventManager() I.EventManager {
	return r.EventManager
}

func (r *RequestCreator) createPrechecker() I.Prechecker {
	if r.provider.NewPrechecker != nil {
		return r.provider.NewPrechecker(r.CreateEventManager())
	}
	return prechecker.NewPrechecker(r.CreateEventManager())
}

type PushRequestCreatorConstructor func(creator Creator, uuid string, request request.PostDeploymentRequest, buffer *bytes.Buffer) I.RequestCreator

func NewPushRequestCreator(creator Creator, uuid string, request request.PostDeploymentRequest, buffer *bytes.Buffer) I.RequestCreator {
	return &PushRequestCreator{
		RequestCreator: newRequestCreator(creator, uuid, buffer),
		Request:        request,
	}
}

type PushRequestCreator struct {
	RequestCreator
	Request request.PostDeploymentRequest
}

func (r PushRequestCreator) CreateRequestProcessor() I.RequestProcessor {
	if r.provider.NewPushRequestProcessor != nil {
		return r.provider.NewPushRequestProcessor(r.Log, r.CreatePushController(), r.Request, r.Buffer)
	}
	return push.NewPushRequestProcessor(r.Log, r.CreatePushController(), r.Request, r.Buffer)
}

func (r PushRequestCreator) CreatePushController() request.PushController {
	if r.provider.NewPushController != nil {
		return r.provider.NewPushController(r.Log, r.CreateDeployer(), r.createSilentDeployer(), r.CreateEventManager(), r.createErrorFinder(), r, r.CreateAuthResolver(), r.CreateEnvResolver())
	}
	return push.NewPushController(r.Log, r.CreateDeployer(), r.createSilentDeployer(), r.CreateEventManager(), r.createErrorFinder(), r, r.CreateAuthResolver(), r.CreateEnvResolver())
}

func (r PushRequestCreator) PushManager(deployEventData structs.DeployEventData, auth I.Authorization, env structs.Environment, envVars map[string]string) I.ActionCreator {
	if r.provider.NewPushManager != nil {
		return r.provider.NewPushManager(r.Creator, r.CreateEventManager(), r.Log, r.CreateFetcher(), deployEventData, r.CreateFileSystem(), r.Request.CFContext, auth, env, envVars, r.CreateHealthChecker(), r.CreateRouteMapper())
	} else {
		return push.NewPushManager(r.Creator, r.CreateEventManager(), r.Log, r.CreateFetcher(), deployEventData, r.CreateFileSystem(), r.Request.CFContext, auth, env, envVars, r.CreateHealthChecker(), r.CreateRouteMapper())
	}
}

type StopRequestCreatorConstructor func(creator Creator, uuid string, request request.PutDeploymentRequest, buffer *bytes.Buffer) I.RequestCreator

func NewStopRequestCreator(creator Creator, uuid string, request request.PutDeploymentRequest, buffer *bytes.Buffer) I.RequestCreator {
	return &StopRequestCreator{
		RequestCreator: newRequestCreator(creator, uuid, buffer),
		Request:        request,
	}
}

type StopRequestCreator struct {
	RequestCreator
	Request request.PutDeploymentRequest
}

func (r StopRequestCreator) CreateRequestProcessor() I.RequestProcessor {
	if r.provider.NewStopRequestProcessor != nil {
		return r.provider.NewStopRequestProcessor(r.Log, r.CreateStopController(), r.Request, r.Buffer)
	}
	return stop.NewStopRequestProcessor(r.Log, r.CreateStopController(), r.Request, r.Buffer)
}

func (r StopRequestCreator) CreateStopController() request.StopController {
	if r.provider.NewStopController != nil {
		return r.provider.NewStopController(r.Log, r.CreateDeployer(), r.CreateEventManager(), r.createErrorFinder(), r, r.CreateAuthResolver(), r.CreateEnvResolver())
	}
	return stop.NewStopController(r.Log, r.CreateDeployer(), r.CreateEventManager(), r.createErrorFinder(), r, r.CreateAuthResolver(), r.CreateEnvResolver())
}

func (r StopRequestCreator) StopManager(deployEventData structs.DeployEventData) I.ActionCreator {
	if r.provider.NewStopManager != nil {
		return r.provider.NewStopManager(r.Creator, r.CreateEventManager(), r.Log, deployEventData)
	} else {
		return stop.NewStopManager(r.Creator, r.CreateEventManager(), r.Log, deployEventData)
	}
}

type StartRequestCreatorConstructor func(creator Creator, uuid string, request request.PutDeploymentRequest, buffer *bytes.Buffer) I.RequestCreator

func NewStartRequestCreator(creator Creator, uuid string, request request.PutDeploymentRequest, buffer *bytes.Buffer) I.RequestCreator {
	return &StartRequestCreator{
		RequestCreator: newRequestCreator(creator, uuid, buffer),
		Request:        request,
	}
}

type StartRequestCreator struct {
	RequestCreator
	Request request.PutDeploymentRequest
}

func (r StartRequestCreator) CreateRequestProcessor() I.RequestProcessor {
	if r.provider.NewStartRequestProcessor != nil {
		return r.provider.NewStartRequestProcessor(r.Log, r.CreateStartController(), r.Request, r.Buffer)
	}
	return start.NewStartRequestProcessor(r.Log, r.CreateStartController(), r.Request, r.Buffer)
}

func (r StartRequestCreator) CreateStartController() request.StartController {
	if r.provider.NewStartController != nil {
		return r.provider.NewStartController(r.Log, r.CreateDeployer(), r.CreateEventManager(), r.createErrorFinder(), r, r.CreateAuthResolver(), r.CreateEnvResolver())
	}
	return start.NewStartController(r.Log, r.CreateDeployer(), r.CreateEventManager(), r.createErrorFinder(), r, r.CreateAuthResolver(), r.CreateEnvResolver())
}

func (r StartRequestCreator) StartManager(deployEventData structs.DeployEventData) I.ActionCreator {
	if r.provider.NewStartManager != nil {
		return r.provider.NewStartManager(r.Creator, r.CreateEventManager(), r.Log, deployEventData)
	} else {
		return start.NewStartManager(r.Creator, r.CreateEventManager(), r.Log, deployEventData)
	}
}

type DeleteRequestCreatorConstructor func(creator Creator, uuid string, request request.DeleteDeploymentRequest, buffer *bytes.Buffer) I.RequestCreator

func NewDeleteRequestCreator(creator Creator, uuid string, request request.DeleteDeploymentRequest, buffer *bytes.Buffer) I.RequestCreator {
	return &DeleteRequestCreator{
		RequestCreator: newRequestCreator(creator, uuid, buffer),
		Request:        request,
	}
}

type DeleteRequestCreator struct {
	RequestCreator
	Request request.DeleteDeploymentRequest
}

func (r DeleteRequestCreator) CreateRequestProcessor() I.RequestProcessor {
	if r.provider.NewDeleteRequestProcessor != nil {
		return r.provider.NewDeleteRequestProcessor(r.Log, r.CreateDeleteController(), r.Request, r.Buffer)
	}
	return delete.NewDeleteRequestProcessor(r.Log, r.CreateDeleteController(), r.Request, r.Buffer)
}

func (r DeleteRequestCreator) CreateDeleteController() request.DeleteController {
	if r.provider.NewDeleteController != nil {
		return r.provider.NewDeleteController(r.Log, r.CreateDeployer(), r.CreateEventManager(), r.createErrorFinder(), r, r.CreateAuthResolver(), r.CreateEnvResolver())
	}
	return delete.NewDeleteController(r.Log, r.CreateDeployer(), r.CreateEventManager(), r.createErrorFinder(), r, r.CreateAuthResolver(), r.CreateEnvResolver())
}

func (r DeleteRequestCreator) DeleteManager(deployEventData structs.DeployEventData) I.ActionCreator {
	if r.provider.NewDeleteManager != nil {
		return r.provider.NewDeleteManager(r.Creator, r.CreateEventManager(), r.Log, deployEventData)
	} else {
		return delete.NewDeleteManager(r.Creator, r.CreateEventManager(), r.Log, deployEventData)
	}
}
