package mocks

import (
	"fmt"
	"io"
)

// Pusher handmade mock for tests.
type Pusher struct {
	Response io.ReadWriter

	InitiallyCall struct {
		TimesCalled int
		Write       struct {
			Output string
		}
		Returns struct {
			Error error
		}
	}

	ExecuteCall struct {
		Write struct {
			Output string
		}
		Returns struct {
			Error error
		}
	}

	PostExecuteCall struct {
		Write struct {
			Output string
		}
		Returns struct {
			Error error
		}
	}

	VerifyCall struct {
		Returns struct {
			Error error
		}
	}

	UndoCall struct {
		Returns struct {
			Error error
		}
	}

	SuccessCall struct {
		Returns struct {
			Error error
		}
	}

	FinallyCall struct {
		Returns struct {
			Error error
		}
	}
}

// Login mock method.
func (p *Pusher) Initially() error {
	p.InitiallyCall.TimesCalled++
	fmt.Fprint(p.Response, p.InitiallyCall.Write.Output)

	return p.InitiallyCall.Returns.Error
}

// Push mock method.
func (p *Pusher) Execute() error {

	fmt.Fprint(p.Response, p.ExecuteCall.Write.Output)

	return p.ExecuteCall.Returns.Error
}

func (p *Pusher) PostExecute() error {

	fmt.Fprint(p.Response, p.PostExecuteCall.Write.Output)

	return p.PostExecuteCall.Returns.Error
}

func (p *Pusher) Verify() error {
	return p.VerifyCall.Returns.Error
}

// FinishPush mock method.
func (p *Pusher) Success() error {
	return p.SuccessCall.Returns.Error
}

// UndoPush mock method.
func (p *Pusher) Undo() error {
	return p.UndoCall.Returns.Error
}

// CleanUp mock method.
func (p *Pusher) Finally() error {
	return p.FinallyCall.Returns.Error
}
