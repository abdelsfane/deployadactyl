package interfaces

import (
	"io"

	S "github.com/compozed/deployadactyl/structs"
)

type Action interface {
	Initially() error
	Verify() error
	Execute() error
	PostExecute() error
	Success() error
	Undo() error
	Finally() error
}

type ActionCreator interface {
	SetUp() error
	CleanUp()
	OnStart() error
	OnFinish(environment S.Environment, response io.ReadWriter, err error) DeployResponse
	Create(environment S.Environment, response io.ReadWriter, foundationURL string) (Action, error)
	InitiallyError(initiallyErrors []error) error
	ExecuteError(executeErrors []error) error
	UndoError(executeErrors, undoErrors []error) error
	SuccessError(successErrors []error) error
}
