package delete

import (
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/compozed/deployadactyl/controller/deployer/bluegreen"
	I "github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/state"
	S "github.com/compozed/deployadactyl/structs"
)

const successfulDelete = `Your delete was successful! (^_^)b

`

type DeleteManagerConstructor func(courierCreator I.CourierCreator, eventManager I.EventManager, log I.DeploymentLogger, deployEventData S.DeployEventData) I.ActionCreator

func NewDeleteManager(c I.CourierCreator, em I.EventManager, log I.DeploymentLogger, ded S.DeployEventData) I.ActionCreator {
	return &DeleteManager{
		CourierCreator:  c,
		EventManager:    em,
		Log:             log,
		DeployEventData: ded,
	}
}

type DeleteManager struct {
	CourierCreator  I.CourierCreator
	EventManager    I.EventManager
	Log             I.DeploymentLogger
	DeployEventData S.DeployEventData
}

func (a DeleteManager) Logger() I.DeploymentLogger {
	return a.Log
}

func (a DeleteManager) SetUp() error {
	return nil
}

func (a DeleteManager) OnStart() error {
	return nil
}

func (a DeleteManager) OnFinish(env S.Environment, response io.ReadWriter, err error) I.DeployResponse {
	if err != nil {
		fmt.Fprintf(response, "\nYour application was not successfully delete on all foundations: %s\n\n", err.Error())
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

	a.Log.Infof("successfully deleted application %s", a.DeployEventData.DeploymentInfo.AppName)
	fmt.Fprintf(response, "\n%s", successfulDelete)

	return I.DeployResponse{StatusCode: http.StatusOK}
}

func (a DeleteManager) CleanUp() {}

func (a DeleteManager) Create(environment S.Environment, response io.ReadWriter, foundationURL string) (I.Action, error) {
	courier, err := a.CourierCreator.CreateCourier()
	if err != nil {
		a.Log.Error(err)
		return &Deleter{}, state.CourierCreationError{Err: err}
	}
	p := &Deleter{
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

func (a DeleteManager) InitiallyError(initiallyErrors []error) error {
	return bluegreen.LoginError{LoginErrors: initiallyErrors}
}

func (a DeleteManager) ExecuteError(executeErrors []error) error {
	return bluegreen.DeleteError{Errors: executeErrors}
}

func (a DeleteManager) UndoError(executeErrors, undoErrors []error) error {
	return bluegreen.RollbackDeleteError{DeleteErrors: executeErrors, RollbackErrors: undoErrors}
}

func (a DeleteManager) SuccessError(successErrors []error) error {
	return bluegreen.FinishDeleteError{FinishDeleteErrors: successErrors}
}
