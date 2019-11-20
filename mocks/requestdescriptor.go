package mocks

import "github.com/compozed/deployadactyl/interfaces"

type RequestDescriptor struct {
	GetIdCall struct {
		Returns struct {
			Id string
		}
	}
	GetDataCall struct {
		Returns struct {
			Data map[string]interface{}
		}
	}
	GetContextCall struct {
		Returns struct {
			Context interfaces.CFContext
		}
	}
	SetContextCall struct {
		Received struct {
			Context interfaces.CFContext
		}
		Returns struct {
			Descriptor interfaces.RequestDescriptor
		}
	}
	GetAuthorizationCall struct {
		Returns struct {
			Authorization interfaces.Authorization
		}
	}
	SetAuthorizationCall struct {
		Received struct {
			Authorization interfaces.Authorization
		}
		Returns struct {
			Descriptor interfaces.RequestDescriptor
		}
	}
	GetRequestCall struct {
		Returns struct {
			Request interface{}
		}
	}
	SetRequestCall struct {
		Received struct {
			Request interface{}
		}
		Returns struct {
			Descriptor interfaces.RequestDescriptor
			Err        error
		}
	}
}

func (d *RequestDescriptor) GetId() string {
	return d.GetIdCall.Returns.Id
}

func (d *RequestDescriptor) GetData() map[string]interface{} {
	return d.GetDataCall.Returns.Data
}

func (d *RequestDescriptor) GetContext() interfaces.CFContext {
	return d.GetContextCall.Returns.Context
}

func (d *RequestDescriptor) SetContext(context interfaces.CFContext) interfaces.RequestDescriptor {
	d.SetContextCall.Received.Context = context
	return d.SetContextCall.Returns.Descriptor
}

func (d *RequestDescriptor) GetAuthorization() interfaces.Authorization {
	return d.GetAuthorizationCall.Returns.Authorization
}

func (d *RequestDescriptor) SetAuthorization(authorization interfaces.Authorization) interfaces.RequestDescriptor {
	d.SetAuthorizationCall.Received.Authorization = authorization
	return d.SetAuthorizationCall.Returns.Descriptor
}

func (d *RequestDescriptor) GetRequest() interface{} {
	return d.GetRequestCall.Returns.Request
}

func (d *RequestDescriptor) SetRequest(request interface{}) (interfaces.RequestDescriptor, error) {
	d.SetRequestCall.Received.Request = request
	return d.SetRequestCall.Returns.Descriptor, d.SetRequestCall.Returns.Err
}
