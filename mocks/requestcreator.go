package mocks

import "github.com/compozed/deployadactyl/interfaces"

type RequestCreator struct {
	CreateRequestProcessorCall struct {
		Returns struct {
			Processor interfaces.RequestProcessor
		}
	}
}

func (rc *RequestCreator) CreateRequestProcessor() interfaces.RequestProcessor {
	return rc.CreateRequestProcessorCall.Returns.Processor
}
