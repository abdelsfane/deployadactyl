package mocks

import (
	"io"

	"bytes"

	I "github.com/compozed/deployadactyl/interfaces"
	S "github.com/compozed/deployadactyl/structs"
)

// BlueGreener handmade mock for tests.
type BlueGreener struct {
	ExecuteCall struct {
		Write    string
		Received struct {
			ActionCreator I.ActionCreator
			Environment   S.Environment
			Out           io.Writer
		}
		Returns struct {
			Error I.DeploymentError
		}
	}
}

// Push mock method.
func (b *BlueGreener) Execute(actionCreator I.ActionCreator, environment S.Environment, out io.ReadWriter) error {
	b.ExecuteCall.Received.ActionCreator = actionCreator
	b.ExecuteCall.Received.Environment = environment
	b.ExecuteCall.Received.Out = out

	if b.ExecuteCall.Write != "" {
		bytes.NewBufferString(b.ExecuteCall.Write).WriteTo(out)
	}
	return b.ExecuteCall.Returns.Error
}
