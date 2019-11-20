package start

import (
	"io"

	"fmt"
	"github.com/compozed/deployadactyl/controller/deployer/bluegreen"
	I "github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/state"
	S "github.com/compozed/deployadactyl/structs"
	"net/http"
	"regexp"
)

const successfulStart = `Your start was successful! (^_^)b

`

type StartManagerConstructor func(courierCreator I.CourierCreator, eventManager I.EventManager, logger I.DeploymentLogger, deployEventData S.DeployEventData) I.ActionCreator

func NewStartManager(c I.CourierCreator, em I.EventManager, l I.DeploymentLogger, d S.DeployEventData) I.ActionCreator {
	return &StartManager{
		CourierCreator:  c,
		EventManager:    em,
		Logger:          l,
		DeployEventData: d,
	}

}

type StartManager struct {
	CourierCreator  I.CourierCreator
	EventManager    I.EventManager
	Logger          I.DeploymentLogger
	DeployEventData S.DeployEventData
}

func (a StartManager) SetUp() error {
	return nil
}

func (a StartManager) OnStart() error {
	return nil
}

func (a StartManager) OnFinish(env S.Environment, response io.ReadWriter, err error) I.DeployResponse {
	if err != nil {
		fmt.Fprintf(response, "\nYour application was not successfully started on all foundations: %s\n\n", err.Error())
		if matched, _ := regexp.MatchString("login failed", err.Error()); matched {
			return I.DeployResponse{
				StatusCode: http.StatusBadRequest,
				Error:      err,
			}
		}
		return I.DeployResponse{
			Error:      err,
			StatusCode: http.StatusInternalServerError,
		}
	}

	a.Logger.Infof("successfully started application %s", a.DeployEventData.DeploymentInfo.AppName)
	fmt.Fprintf(response, "\n%s", successfulStart)

	return I.DeployResponse{StatusCode: http.StatusOK}
}

func (a StartManager) CleanUp() {}

func (a StartManager) Create(environment S.Environment, response io.ReadWriter, foundationURL string) (I.Action, error) {
	courier, err := a.CourierCreator.CreateCourier()
	if err != nil {
		a.Logger.Error(err)
		return &Starter{}, state.CourierCreationError{Err: err}
	}
	p := &Starter{
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
		Log:           a.Logger,
		FoundationURL: foundationURL,
		AppName:       a.DeployEventData.DeploymentInfo.AppName,
		Data:          a.DeployEventData.DeploymentInfo.Data,
	}

	return p, nil
}

func (a StartManager) InitiallyError(initiallyErrors []error) error {
	return bluegreen.LoginError{LoginErrors: initiallyErrors}
}

func (a StartManager) ExecuteError(executeErrors []error) error {
	return bluegreen.StartError{Errors: executeErrors}
}

func (a StartManager) UndoError(executeErrors, undoErrors []error) error {
	return bluegreen.RollbackStartError{StartErrors: executeErrors, RollbackErrors: undoErrors}
}

func (a StartManager) SuccessError(successErrors []error) error {
	return bluegreen.FinishStartError{FinishStartErrors: successErrors}
}
