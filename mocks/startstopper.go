package mocks

type StartStopper struct {
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
	UndoCall struct {
		Write   string
		Returns struct {
			Error error
		}
	}
	VerifyCall struct {
		Write   string
		Returns struct {
			Error error
		}
	}
	SuccessCall struct {
		Write   string
		Returns struct {
			Error error
		}
	}
	FinallyCall struct {
		Write   string
		Returns struct {
			Error error
		}
	}
}

func (s *StartStopper) Initially() error {

	return s.InitiallyCall.Returns.Error
}

func (s *StartStopper) Verify() error {

	return s.VerifyCall.Returns.Error
}

func (s *StartStopper) Finally() error {

	return s.FinallyCall.Returns.Error
}

func (s *StartStopper) Success() error {

	return s.SuccessCall.Returns.Error
}

func (s *StartStopper) Execute() error {

	return s.ExecuteCall.Returns.Error
}

func (s *StartStopper) PostExecute() error {

	return s.PostExecuteCall.Returns.Error
}

func (s *StartStopper) Undo() error {

	return s.UndoCall.Returns.Error
}
