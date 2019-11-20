package mocks

import (
	"github.com/compozed/deployadactyl/interfaces"
	S "github.com/compozed/deployadactyl/structs"

	"github.com/compozed/deployadactyl/controller/deployer/bluegreen"
	"io"
)

type receivedCall struct {
	FoundationURL string
	Response      io.ReadWriter
}

type StopManager struct {
	CreateStopperCall struct {
		TimesCalled int
		Received    []receivedCall
		Returns     struct {
			Stoppers []interfaces.Action
			Error    []error
		}
	}
}

func (s *StopManager) SetUp() error {
	return nil
}

func (s *StopManager) OnStart() error {
	return nil
}

func (s *StopManager) OnFinish(env S.Environment, response io.ReadWriter, err error) interfaces.DeployResponse {
	return interfaces.DeployResponse{}
}

func (s *StopManager) CleanUp() {}

func (s *StopManager) InitiallyError(initiallyErrors []error) error {
	return bluegreen.LoginError{LoginErrors: initiallyErrors}
}

func (s *StopManager) Create(environment S.Environment, response io.ReadWriter, foundationURL string) (interfaces.Action, error) {
	defer func() { s.CreateStopperCall.TimesCalled++ }()

	received := receivedCall{
		FoundationURL: foundationURL,
		Response:      response,
	}
	s.CreateStopperCall.Received = append(s.CreateStopperCall.Received, received)

	return s.CreateStopperCall.Returns.Stoppers[s.CreateStopperCall.TimesCalled], s.CreateStopperCall.Returns.Error[s.CreateStopperCall.TimesCalled]
}

func (s *StopManager) ExecuteError(executeErrors []error) error {
	return bluegreen.StopError{Errors: executeErrors}
}

func (s *StopManager) UndoError(executeErrors, undoErrors []error) error {
	return bluegreen.RollbackStopError{StopErrors: executeErrors, RollbackErrors: undoErrors}
}

func (s *StopManager) SuccessError(successErrors []error) error {
	return bluegreen.FinishStopError{FinishStopErrors: successErrors}
}
