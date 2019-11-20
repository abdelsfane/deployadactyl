package mocks

// Action handmade mock for tests.
type Action struct {
	InitiallyCall struct {
		Returns struct {
			Error error
		}
	}
	ExecuteCall struct {
		Returns struct {
			Error error
		}
	}
	PostExecuteCall struct {
		Returns struct {
			Error error
		}
	}
	VerifyCall struct {
		Returns struct {
			Error error
		}
	}
	SuccessCall struct {
		Returns struct {
			Error error
		}
	}
	UndoCall struct {
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

// Action mock method.
func (a *Action) Initially() error {

	return a.InitiallyCall.Returns.Error
}

func (a *Action) Execute() error {
	return a.ExecuteCall.Returns.Error
}

func (a *Action) PostExecute() error {
	return a.ExecuteCall.Returns.Error
}

func (a *Action) Verify() error {

	return a.VerifyCall.Returns.Error
}

func (a *Action) Success() error {

	return a.SuccessCall.Returns.Error
}

func (a *Action) Undo() error {

	return a.UndoCall.Returns.Error
}

func (a *Action) Finally() error {

	return a.FinallyCall.Returns.Error
}
