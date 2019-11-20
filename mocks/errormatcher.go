package mocks

import "github.com/compozed/deployadactyl/interfaces"

type ErrorMatcherMock struct {
	MatchCall struct {
		Returns interfaces.LogMatchedError
	}
	DescriptorCall struct {
		Returns string
	}
}

func (m *ErrorMatcherMock) Match(matchTo []byte) interfaces.LogMatchedError {
	return m.MatchCall.Returns
}

func (m *ErrorMatcherMock) Descriptor() string {
	return m.DescriptorCall.Returns
}
