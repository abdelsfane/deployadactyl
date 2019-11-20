package mocks

import (
	"github.com/compozed/deployadactyl/interfaces"
	S "github.com/compozed/deployadactyl/structs"

	"io"

	"github.com/compozed/deployadactyl/controller/deployer/bluegreen"
)

type DeleteManager struct {
	CreateStarterCall struct {
		TimesCalled int
		Received    []receivedCall
		Returns     struct {
			Starters []interfaces.Action
			Error    []error
		}
	}
}

func (s *DeleteManager) SetUp() error {
	return nil
}

func (s *DeleteManager) OnStart() error {
	return nil
}

func (s *DeleteManager) OnFinish(env S.Environment, response io.ReadWriter, err error) interfaces.DeployResponse {
	return interfaces.DeployResponse{}
}

func (s *DeleteManager) CleanUp() {}

func (s *DeleteManager) InitiallyError(initiallyErrors []error) error {
	return bluegreen.LoginError{LoginErrors: initiallyErrors}
}

func (s *DeleteManager) Create(environment S.Environment, response io.ReadWriter, foundationURL string) (interfaces.Action, error) {
	defer func() { s.CreateStarterCall.TimesCalled++ }()

	received := receivedCall{
		FoundationURL: foundationURL,
		Response:      response,
	}
	s.CreateStarterCall.Received = append(s.CreateStarterCall.Received, received)

	return s.CreateStarterCall.Returns.Starters[s.CreateStarterCall.TimesCalled], s.CreateStarterCall.Returns.Error[s.CreateStarterCall.TimesCalled]
}

func (s *DeleteManager) ExecuteError(executeErrors []error) error {
	return bluegreen.StartError{Errors: executeErrors}
}

func (s *DeleteManager) UndoError(executeErrors, undoErrors []error) error {
	return bluegreen.RollbackStartError{StartErrors: executeErrors, RollbackErrors: undoErrors}
}

func (s *DeleteManager) SuccessError(successErrors []error) error {
	return bluegreen.FinishStartError{FinishStartErrors: successErrors}
}
