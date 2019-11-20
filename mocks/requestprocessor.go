package mocks

import (
	"bytes"

	"github.com/compozed/deployadactyl/interfaces"
)

type RequestProcessor struct {
	Response    *bytes.Buffer
	ProcessCall struct {
		TimesCalled int
		Returns     struct {
			Response interfaces.DeployResponse
		}
		Writes string
	}
}

func (c *RequestProcessor) Process() interfaces.DeployResponse {
	c.ProcessCall.TimesCalled++

	if c.Response != nil {
		c.Response.Write([]byte(c.ProcessCall.Writes))
	}
	return c.ProcessCall.Returns.Response
}
