// Package bluegreen is responsible for concurrently pushing an application to multiple Cloud Foundry instances.
package bluegreen

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	I "github.com/compozed/deployadactyl/interfaces"
	S "github.com/compozed/deployadactyl/structs"
)

// BlueGreen has a PushManager to creater pushers for blue green deployments.
type BlueGreenConstructor func(log I.DeploymentLogger) I.BlueGreener

func NewBlueGreen(log I.DeploymentLogger) I.BlueGreener {
	return &BlueGreen{
		Log: log,
	}
}

type BlueGreen struct {
	Log I.DeploymentLogger
}

// Push will login to all the Cloud Foundry instances provided in the Config and then push the application to all the instances concurrently.
// If the application fails to start in any of the instances it handles rolling back the application in every instance, unless it is the first deploy.
func (bg BlueGreen) Execute(actionCreator I.ActionCreator, environment S.Environment, response io.ReadWriter) error {

	actors := make([]actor, len(environment.Foundations))
	buffers := make([]*bytes.Buffer, len(environment.Foundations))

	for i, foundationURL := range environment.Foundations {
		buffers[i] = &bytes.Buffer{}

		action, err := actionCreator.Create(environment, buffers[i], foundationURL)
		if err != nil {
			return InitializationError{err}
		}
		defer action.Finally()

		actors[i] = NewActor(action)
		defer close(actors[i].Commands)
	}

	defer func() {
		for _, buffer := range buffers {
			fmt.Fprintf(response, "\n%s Cloud Foundry Output %s\n", strings.Repeat("-", 19), strings.Repeat("-", 19))
			buffer.WriteTo(response)
		}

		fmt.Fprintf(response, "\n%s End Cloud Foundry Output %s\n", strings.Repeat("-", 17), strings.Repeat("-", 17))
	}()

	initLoginError := bg.commands(append(make([]actor, 0), actors[0]), func(action I.Action) error {
		return action.Initially()
	})
	if len(initLoginError) != 0 {
		return actionCreator.InitiallyError(initLoginError)
	}

	loginErrors := bg.commands(actors, func(action I.Action) error {
		return action.Initially()
	})

	if len(loginErrors) != 0 {
		return actionCreator.InitiallyError(loginErrors)
	}

	actionErrors := bg.commands(actors, func(action I.Action) error {
		return action.Execute()
	})

	if len(actionErrors) != 0 {
		return bg.processErrors(actionErrors, actors, actionCreator)
	}

	actionErrors = bg.commands(actors, func(action I.Action) error {
		return action.PostExecute()
	})

	if len(actionErrors) != 0 {
		return bg.processErrors(actionErrors, actors, actionCreator)
	}

	finishActionErrors := bg.commands(actors, func(action I.Action) error {
		return action.Success()
	})
	if len(finishActionErrors) != 0 {
		return actionCreator.SuccessError(finishActionErrors)
	}

	return nil
}

func (bg BlueGreen) commands(actors []actor, doFunc ActorCommand) (manyErrors []error) {

	for _, a := range actors {
		a.Commands <- doFunc
	}
	for _, a := range actors {
		if err := <-a.Errs; err != nil {
			manyErrors = append(manyErrors, err)
		}
	}
	return
}

func (bg BlueGreen) processErrors(actionErrors []error, actors []actor, actionCreator I.ActionCreator) error {
	bg.Log.Errorf("failed to execute action against all foundations - rolling back action")
	rollbackErrors := bg.commands(actors, func(action I.Action) error {
		return action.Undo()
	})

	if len(rollbackErrors) != 0 {
		return actionCreator.UndoError(actionErrors, rollbackErrors)
	}

	return actionCreator.ExecuteError(actionErrors)
}
