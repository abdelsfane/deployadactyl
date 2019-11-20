package mocks

import (
	"github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/structs"
)

// PushManager handmade mock for tests.
type PushManagerFactory struct {
	PushManagerCall struct {
		Called   bool
		Received struct {
			DeployEventData structs.DeployEventData
			Auth            interfaces.Authorization
			Environment     structs.Environment
			EnvVars         map[string]string
		}
		Returns struct {
			ActionCreator interfaces.ActionCreator
		}
	}
}

// CreatePusher mock method.

func (p *PushManagerFactory) PushManager(deployEventData structs.DeployEventData, auth interfaces.Authorization, env structs.Environment, envVars map[string]string) interfaces.ActionCreator {
	p.PushManagerCall.Called = true
	p.PushManagerCall.Received.DeployEventData = deployEventData
	p.PushManagerCall.Received.Auth = auth
	p.PushManagerCall.Received.Environment = env
	p.PushManagerCall.Received.EnvVars = envVars

	return p.PushManagerCall.Returns.ActionCreator
}

type StopManagerFactory struct {
	StopManagerCall struct {
		Called   bool
		Received struct {
			DeployEventData structs.DeployEventData
		}
		Returns struct {
			ActionCreater interfaces.ActionCreator
		}
	}
}

func (s *StopManagerFactory) StopManager(DeployEventData structs.DeployEventData) interfaces.ActionCreator {
	s.StopManagerCall.Called = true
	s.StopManagerCall.Received.DeployEventData = DeployEventData

	return s.StopManagerCall.Returns.ActionCreater
}

type StartManagerFactory struct {
	StartManagerCall struct {
		Called   bool
		Received struct {
			DeployEventData structs.DeployEventData
		}
		Returns struct {
			ActionCreater interfaces.ActionCreator
		}
	}
}

func (t *StartManagerFactory) StartManager(DeployEventData structs.DeployEventData) interfaces.ActionCreator {
	t.StartManagerCall.Called = true
	t.StartManagerCall.Received.DeployEventData = DeployEventData

	return t.StartManagerCall.Returns.ActionCreater
}

type DeleteManagerFactory struct {
	DeleteManagerCall struct {
		Called   bool
		Received struct {
			DeployEventData structs.DeployEventData
		}
		Returns struct {
			ActionCreater interfaces.ActionCreator
		}
	}
}

func (t *DeleteManagerFactory) DeleteManager(DeployEventData structs.DeployEventData) interfaces.ActionCreator {
	t.DeleteManagerCall.Called = true
	t.DeleteManagerCall.Received.DeployEventData = DeployEventData

	return t.DeleteManagerCall.Returns.ActionCreater
}
