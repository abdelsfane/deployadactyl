package stop

import (
	"fmt"
	"github.com/compozed/deployadactyl/controller/deployer/bluegreen"
	I "github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/state"
	S "github.com/compozed/deployadactyl/structs"
	"io"
	"net/http"
	"regexp"
)

const successfulStop = `Your stop was successful! (^_^)b

`

type StopManagerConstructor func(courierCreator I.CourierCreator, eventManager I.EventManager, log I.DeploymentLogger, deployEventData S.DeployEventData) I.ActionCreator

func NewStopManager(c I.CourierCreator, em I.EventManager, log I.DeploymentLogger, ded S.DeployEventData) I.ActionCreator {
	return &StopManager{
		CourierCreator:  c,
		EventManager:    em,
		Log:             log,
		DeployEventData: ded,
	}
}

type StopManager struct {
	CourierCreator  I.CourierCreator
	EventManager    I.EventManager
	Log             I.DeploymentLogger
	DeployEventData S.DeployEventData
}

func (a StopManager) Logger() I.DeploymentLogger {
	return a.Log
}

func (a StopManager) SetUp() error {
	return nil
}

func (a StopManager) OnStart() error {
	return nil
}

func (a StopManager) OnFinish(env S.Environment, response io.ReadWriter, err error) I.DeployResponse {
	if err != nil {
		fmt.Fprintf(response, "\nYour application was not successfully stopped on all foundations: %s\n\n", err.Error())
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

	a.Log.Infof("successfully stopped application %s", a.DeployEventData.DeploymentInfo.AppName)
	fmt.Fprintf(response, "\n%s", successfulStop)

	return I.DeployResponse{StatusCode: http.StatusOK}
}

func (a StopManager) CleanUp() {}

func (a StopManager) Create(environment S.Environment, response io.ReadWriter, foundationURL string) (I.Action, error) {
	courier, err := a.CourierCreator.CreateCourier()
	if err != nil {
		a.Log.Error(err)
		return &Stopper{}, state.CourierCreationError{Err: err}
	}
	p := &Stopper{
		Courier: courier,
		CFContext: I.CFContext{
			Environment:  environment.Name,
			Organization: a.DeployEventData.DeploymentInfo.Org,
			Space:        a.DeployEventData.DeploymentInfo.Space,
			Application:  a.DeployEventData.DeploymentInfo.AppName,
			SkipSSL:      a.DeployEventData.DeploymentInfo.SkipSSL,
		},
		Authorization: I.Authorization{
			Username: a.DeployEventData.DeploymentInfo.Username,
			Password: a.DeployEventData.DeploymentInfo.Password,
		},
		EventManager:  a.EventManager,
		Response:      response,
		Log:           a.Log,
		FoundationURL: foundationURL,
		AppName:       a.DeployEventData.DeploymentInfo.AppName,
	}

	return p, nil
}

func (a StopManager) InitiallyError(initiallyErrors []error) error {
	return bluegreen.LoginError{LoginErrors: initiallyErrors}
}

func (a StopManager) ExecuteError(executeErrors []error) error {
	return bluegreen.StopError{Errors: executeErrors}
}

func (a StopManager) UndoError(executeErrors, undoErrors []error) error {
	return bluegreen.RollbackStopError{StopErrors: executeErrors, RollbackErrors: undoErrors}
}

func (a StopManager) SuccessError(successErrors []error) error {
	return bluegreen.FinishStopError{FinishStopErrors: successErrors}
}
