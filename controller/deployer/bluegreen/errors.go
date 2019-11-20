package bluegreen

import (
	"errors"
	"fmt"
)

type LoginError struct {
	LoginErrors []error
}

func (e LoginError) Error() string {
	errs := makeErrorString(e.LoginErrors)
	return fmt.Sprintf("login failed: %s", errs)
}

func (e LoginError) Code() string {
	return "LoginError"
}

type PushError struct {
	PushErrors []error
}

func (e PushError) Error() string {
	errs := makeErrorString(e.PushErrors)
	return fmt.Sprintf("push failed: %s", errs)
}

func (e PushError) Code() string {
	return "PushError"
}

type RollbackError struct {
	PushErrors     []error
	RollbackErrors []error
}

type RollbackStopError struct {
	StopErrors     []error
	RollbackErrors []error
}

func (e RollbackError) Error() string {
	var (
		pushErrs       = makeErrorString(e.PushErrors)
		rollbackErrors = makeErrorString(e.RollbackErrors)
	)

	return fmt.Sprintf("push failed: %s: rollback failed: %s", pushErrs, rollbackErrors)
}

func (e RollbackStopError) Error() string {
	var (
		stopErrs           = makeErrorString(e.StopErrors)
		rollbackStopErrors = makeErrorString(e.RollbackErrors)
	)

	return fmt.Sprintf("stop failed: %s: rollback failed: %s", stopErrs, rollbackStopErrors)
}

func (e RollbackError) Code() string {
	return "RollbackError"
}

type FinishPushError struct {
	FinishPushError []error
}

func (e FinishPushError) Error() string {
	var (
		finishPushErrors = makeErrorString(e.FinishPushError)
	)

	return fmt.Sprintf("finish push failed: %s", finishPushErrors)
}

func (e FinishPushError) Code() string {
	return "FinishPushError"
}

type StartStopError struct {
	Err error
}

func (e StartStopError) Error() string {
	return e.Err.Error()
}

type InitializationError struct {
	Err error
}

func (e InitializationError) Error() string {
	return e.Err.Error()
}

func (e InitializationError) Code() string {
	return "InitError"
}

type FinishStopError struct {
	FinishStopErrors []error
}

func (e FinishStopError) Error() string {
	finishStopErrors := makeErrorString(e.FinishStopErrors)

	return fmt.Sprintf("finish stop failed: %s", finishStopErrors)
}

type StopError struct {
	Errors []error
}

func (e StopError) Error() string {
	errs := makeErrorString(e.Errors)
	return fmt.Sprintf("stop failed: %s", errs)
}

func (e StopError) Code() string {
	return "StopError"
}

type FinishDeployError struct {
	Err error
}

func (e FinishDeployError) Error() string {
	return e.Err.Error()
}

func (e FinishDeployError) Code() string {
	return "FinishDeployError"
}

func makeErrorString(manyErrors []error) error {
	var result string
	for i, e := range manyErrors {
		if len(e.Error()) != 0 {
			if i == 0 {
				result = e.Error()
			} else {
				result = fmt.Sprintf("%s: %s", result, e.Error())
			}
		}
	}

	return errors.New(result)
}

type FinishStartError struct {
	FinishStartErrors []error
}

func (e FinishStartError) Error() string {
	finishStartErrors := makeErrorString(e.FinishStartErrors)

	return fmt.Sprintf("finish stop failed: %s", finishStartErrors)
}

type StartError struct {
	Errors []error
}

func (e StartError) Error() string {
	errs := makeErrorString(e.Errors)
	return fmt.Sprintf("start failed: %s", errs)
}

func (e StartError) Code() string {
	return "StartError"
}

type RollbackStartError struct {
	StartErrors    []error
	RollbackErrors []error
}

func (e RollbackStartError) Error() string {
	var (
		startErrs           = makeErrorString(e.StartErrors)
		rollbackStartErrors = makeErrorString(e.RollbackErrors)
	)

	return fmt.Sprintf("start failed: %s: rollback failed: %s", startErrs, rollbackStartErrors)
}

type FinishDeleteError struct {
	FinishDeleteErrors []error
}

func (e FinishDeleteError) Error() string {
	finishDeleteErrors := makeErrorString(e.FinishDeleteErrors)

	return fmt.Sprintf("finish delete failed: %s", finishDeleteErrors)
}

type DeleteError struct {
	Errors []error
}

func (e DeleteError) Error() string {
	errs := makeErrorString(e.Errors)
	return fmt.Sprintf("delete failed: %s", errs)
}

func (e DeleteError) Code() string {
	return "DeleteError"
}

type RollbackDeleteError struct {
	DeleteErrors   []error
	RollbackErrors []error
}

func (e RollbackDeleteError) Error() string {
	var (
		startErrs           = makeErrorString(e.DeleteErrors)
		rollbackStartErrors = makeErrorString(e.RollbackErrors)
	)

	return fmt.Sprintf("delete failed: %s: rollback failed: %s", startErrs, rollbackStartErrors)
}
