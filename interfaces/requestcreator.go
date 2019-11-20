package interfaces

type RequestProcessor interface {
	Process() DeployResponse
}

type RequestCreator interface {
	CreateRequestProcessor() RequestProcessor
}
