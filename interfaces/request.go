package interfaces

type RequestDescriptor interface {
	GetId() string
	GetData() map[string]interface{}
	GetContext() CFContext
	SetContext(context CFContext) RequestDescriptor
	GetAuthorization() Authorization
	SetAuthorization(authorization Authorization) RequestDescriptor
	GetRequest() interface{}
	SetRequest(request interface{}) (RequestDescriptor, error)
}
