package mocks

import (
	"github.com/compozed/deployadactyl/interfaces"
	S "github.com/compozed/deployadactyl/structs"

	"io"

	"github.com/compozed/deployadactyl/controller/deployer/bluegreen"
)

type StartManager struct {
	CreateStarterCall struct {
		TimesCalled int
		Received    []receivedCall
		Returns     struct {
			Starters []interfaces.Action
			Error    []error
		}
	}
}

func (s *StartManager) SetUp() error {
	return nil
}

func (s *StartManager) OnStart() error {
	return nil
}

func (s *StartManager) OnFinish(env S.Environment, response io.ReadWriter, err error) interfaces.DeployResponse {
	return interfaces.DeployResponse{}
}

func (s *StartManager) CleanUp() {}

func (s *StartManager) InitiallyError(initiallyErrors []error) error {
	return bluegreen.LoginError{LoginErrors: initiallyErrors}
}

func (s *StartManager) Create(environment S.Environment, response io.ReadWriter, foundationURL string) (interfaces.Action, error) {
	defer func() { s.CreateStarterCall.TimesCalled++ }()

	received := receivedCall{
		FoundationURL: foundationURL,
		Response:      response,
	}
	s.CreateStarterCall.Received = append(s.CreateStarterCall.Received, received)

	return s.CreateStarterCall.Returns.Starters[s.CreateStarterCall.TimesCalled], s.CreateStarterCall.Returns.Error[s.CreateStarterCall.TimesCalled]
}

func (s *StartManager) ExecuteError(executeErrors []error) error {
	return bluegreen.StartError{Errors: executeErrors}
}

func (s *StartManager) UndoError(executeErrors, undoErrors []error) error {
	return bluegreen.RollbackStartError{StartErrors: executeErrors, RollbackErrors: undoErrors}
}

func (s *StartManager) SuccessError(successErrors []error) error {
	return bluegreen.FinishStartError{FinishStartErrors: successErrors}
}
