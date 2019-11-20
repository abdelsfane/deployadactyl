package mocks

import (
	"io"

	"github.com/compozed/deployadactyl/controller/deployer/bluegreen"
	"github.com/compozed/deployadactyl/interfaces"
	S "github.com/compozed/deployadactyl/structs"
)

// PushManager handmade mock for tests.
type PushManager struct {
	SetUpCall struct {
		Called  bool
		Returns struct {
			Err error
		}
	}
	OnStartCall struct {
		Called  bool
		Returns struct {
			Err error
		}
	}
	CreatePusherCall struct {
		TimesCalled int
		Returns     struct {
			Pushers []interfaces.Action
			Error   []error
		}
	}
	OnFinishCall struct {
		Called   bool
		Received struct {
			Environment S.Environment
			Response    io.ReadWriter
			Error       error
		}
		Returns struct {
			DeployResponse interfaces.DeployResponse
		}
	}
	CleanUpCall struct {
		Called bool
	}
	InitiallyErrorCall struct {
		Received struct {
			Errs []error
		}
		Returns struct {
			Err error
		}
	}
}

type FileSystemCleaner struct {
	RemoveAllCall struct {
		Called   bool
		Received struct {
			Path string
		}
		Returns struct {
			Error error
		}
	}
}

// CreatePusher mock method.

func (p *FileSystemCleaner) RemoveAll(path string) error {
	p.RemoveAllCall.Called = true

	p.RemoveAllCall.Received.Path = path

	return p.RemoveAllCall.Returns.Error
}

func (p *PushManager) SetUp() error {
	p.SetUpCall.Called = true
	return p.SetUpCall.Returns.Err
}

func (p *PushManager) CleanUp() {
	p.CleanUpCall.Called = true
}

func (p *PushManager) OnStart() error {
	p.OnStartCall.Called = true

	return p.OnStartCall.Returns.Err
}

func (p *PushManager) OnFinish(env S.Environment, response io.ReadWriter, err error) interfaces.DeployResponse {
	p.OnFinishCall.Called = true
	p.OnFinishCall.Received.Environment = env
	p.OnFinishCall.Received.Response = response
	p.OnFinishCall.Received.Error = err

	return p.OnFinishCall.Returns.DeployResponse
}

func (p *PushManager) Create(environment S.Environment, response io.ReadWriter, foundationURL string) (interfaces.Action, error) {
	defer func() { p.CreatePusherCall.TimesCalled++ }()

	return p.CreatePusherCall.Returns.Pushers[p.CreatePusherCall.TimesCalled], p.CreatePusherCall.Returns.Error[p.CreatePusherCall.TimesCalled]
}

func (p *PushManager) InitiallyError(initiallyErrors []error) error {
	p.InitiallyErrorCall.Received.Errs = initiallyErrors

	return p.InitiallyErrorCall.Returns.Err
}

func (p *PushManager) ExecuteError(executeErrors []error) error {
	return bluegreen.PushError{PushErrors: executeErrors}
}

func (p *PushManager) UndoError(executeErrors, undoErrors []error) error {
	return bluegreen.RollbackError{PushErrors: executeErrors, RollbackErrors: undoErrors}
}

func (p *PushManager) SuccessError(successErrors []error) error {
	return bluegreen.FinishPushError{FinishPushError: successErrors}
}
