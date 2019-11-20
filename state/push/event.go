package push

import (
	"errors"
	"github.com/compozed/deployadactyl/eventmanager"
	"github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/structs"
	"io"
	"reflect"
)

type eventBinding struct {
	etype   reflect.Type
	handler func(event interface{}) error
}

func (s eventBinding) Accepts(event interface{}) bool {
	return reflect.TypeOf(event) == s.etype
}

func (b eventBinding) Emit(event interface{}) error {
	return b.handler(event)
}

type DeployStartedEvent struct {
	CFContext   interfaces.CFContext
	ArtifactURL string
	Body        io.Reader
	ContentType string
	Environment structs.Environment
	Auth        interfaces.Authorization
	Response    io.ReadWriter
	Data        map[string]interface{}
	Log         interfaces.DeploymentLogger
}

func (d DeployStartedEvent) Name() string {
	return "DeployStartedEvent"
}

func NewDeployStartEventBinding(handler func(event DeployStartedEvent) error) interfaces.Binding {
	return eventBinding{
		etype: reflect.TypeOf(DeployStartedEvent{}),
		handler: func(gevent interface{}) error {
			event, ok := gevent.(DeployStartedEvent)
			if ok {
				return handler(event)
			} else {
				return eventmanager.InvalidEventType{errors.New("invalid event type")}
			}
		},
	}
}

type DeployFinishedEvent struct {
	CFContext   interfaces.CFContext
	Body        io.Reader
	ContentType string
	Environment structs.Environment
	Auth        interfaces.Authorization
	Response    io.ReadWriter
	Data        map[string]interface{}
	Log         interfaces.DeploymentLogger
}

func (d DeployFinishedEvent) Name() string {
	return "DeployFinishEvent"
}

func NewDeployFinishedEventBinding(handler func(event DeployFinishedEvent) error) interfaces.Binding {
	return eventBinding{
		etype: reflect.TypeOf(DeployFinishedEvent{}),
		handler: func(gevent interface{}) error {
			event, ok := gevent.(DeployFinishedEvent)
			if ok {
				return handler(event)
			} else {
				return eventmanager.InvalidEventType{errors.New("invalid event type")}
			}
		},
	}
}

type DeploySuccessEvent struct {
	CFContext           interfaces.CFContext
	Body                io.Reader
	ContentType         string
	Environment         structs.Environment
	Auth                interfaces.Authorization
	Response            io.ReadWriter
	Data                map[string]interface{}
	HealthCheckEndpoint string
	ArtifactURL         string
	Log                 interfaces.DeploymentLogger
}

func (d DeploySuccessEvent) Name() string {
	return "DeploySuccessEvent"
}

func NewDeploySuccessEventBinding(handler func(event DeploySuccessEvent) error) interfaces.Binding {
	return eventBinding{
		etype: reflect.TypeOf(DeploySuccessEvent{}),
		handler: func(gevent interface{}) error {
			event, ok := gevent.(DeploySuccessEvent)
			if ok {
				return handler(event)
			} else {
				return eventmanager.InvalidEventType{errors.New("invalid event type")}
			}
		},
	}
}

type DeployFailureEvent struct {
	CFContext   interfaces.CFContext
	Body        io.Reader
	ContentType string
	Environment structs.Environment
	Auth        interfaces.Authorization
	Response    io.ReadWriter
	Data        map[string]interface{}
	Error       error
	Log         interfaces.DeploymentLogger
}

func (d DeployFailureEvent) Name() string {
	return "DeployFailureEvent"
}

func NewDeployFailureEventBinding(handler func(event DeployFailureEvent) error) interfaces.Binding {
	return eventBinding{
		etype: reflect.TypeOf(DeployFailureEvent{}),
		handler: func(gevent interface{}) error {
			event, ok := gevent.(DeployFailureEvent)
			if ok {
				return handler(event)
			} else {
				return eventmanager.InvalidEventType{errors.New("invalid event type")}
			}
		},
	}
}

type PushStartedEvent struct {
	CFContext            interfaces.CFContext
	Body                 io.Reader
	ContentType          string
	Environment          structs.Environment
	Auth                 interfaces.Authorization
	Response             io.ReadWriter
	Data                 map[string]interface{}
	Instances            uint16
	EnvironmentVariables map[string]string
	Manifest             string
	Log                  interfaces.DeploymentLogger
}

func (d PushStartedEvent) Name() string {
	return "PushStartedEvent"
}

func NewPushStartedEventBinding(handler func(event PushStartedEvent) error) interfaces.Binding {
	return eventBinding{
		etype: reflect.TypeOf(PushStartedEvent{}),
		handler: func(gevent interface{}) error {
			event, ok := gevent.(PushStartedEvent)
			if ok {
				return handler(event)
			} else {
				return eventmanager.InvalidEventType{errors.New("invalid event type")}
			}
		},
	}
}

type PushFinishedEvent struct {
	CFContext           interfaces.CFContext
	Auth                interfaces.Authorization
	Response            io.ReadWriter
	AppPath             string
	FoundationURL       string
	TempAppWithUUID     string
	Manifest            string
	Data                map[string]interface{}
	Courier             interfaces.Courier
	HealthCheckEndpoint string
	Log                 interfaces.DeploymentLogger
}

func (d PushFinishedEvent) Name() string {
	return "PushFinishedEvent"
}

func NewPushFinishedEventBinding(handler func(event PushFinishedEvent) error) interfaces.Binding {
	return eventBinding{
		etype: reflect.TypeOf(PushFinishedEvent{}),
		handler: func(gevent interface{}) error {
			event, ok := gevent.(PushFinishedEvent)
			if ok {
				return handler(event)
			} else {
				return eventmanager.InvalidEventType{errors.New("invalid event type")}
			}
		},
	}
}

type ArtifactRetrievalStartEvent struct {
	CFContext   interfaces.CFContext
	Auth        interfaces.Authorization
	Environment structs.Environment
	Response    io.ReadWriter
	Data        map[string]interface{}
	Manifest    string
	ArtifactURL string
	Log         interfaces.DeploymentLogger
}

func (d ArtifactRetrievalStartEvent) Name() string {
	return "ArtifactRetrievalStartEvent"
}

func NewArtifactRetrievalStartEventBinding(handler func(event ArtifactRetrievalStartEvent) error) interfaces.Binding {
	return eventBinding{
		etype: reflect.TypeOf(ArtifactRetrievalStartEvent{}),
		handler: func(gevent interface{}) error {
			event, ok := gevent.(ArtifactRetrievalStartEvent)
			if ok {
				return handler(event)
			} else {
				return eventmanager.InvalidEventType{errors.New("invalid event type")}
			}
		},
	}
}

type ArtifactRetrievalFailureEvent struct {
	CFContext   interfaces.CFContext
	Auth        interfaces.Authorization
	Environment structs.Environment
	Response    io.ReadWriter
	Data        map[string]interface{}
	Manifest    string
	ArtifactURL string
	Log         interfaces.DeploymentLogger
}

func (d ArtifactRetrievalFailureEvent) Name() string {
	return "ArtifactRetrievalFailureEvent"
}

func NewArtifactRetrievalFailureEventBinding(handler func(event ArtifactRetrievalFailureEvent) error) interfaces.Binding {
	return eventBinding{
		etype: reflect.TypeOf(ArtifactRetrievalFailureEvent{}),
		handler: func(gevent interface{}) error {
			event, ok := gevent.(ArtifactRetrievalFailureEvent)
			if ok {
				return handler(event)
			} else {
				return eventmanager.InvalidEventType{errors.New("invalid event type")}
			}
		},
	}
}

type ArtifactRetrievalSuccessEvent struct {
	CFContext            interfaces.CFContext
	Auth                 interfaces.Authorization
	Environment          structs.Environment
	Response             io.ReadWriter
	Data                 map[string]interface{}
	Manifest             string
	ArtifactURL          string
	AppPath              string
	EnvironmentVariables map[string]string
	Log                  interfaces.DeploymentLogger
}

func (d ArtifactRetrievalSuccessEvent) Name() string {
	return "ArtifactRetrievalSuccessEvent"
}

func NewArtifactRetrievalSuccessEventBinding(handler func(event ArtifactRetrievalSuccessEvent) error) interfaces.Binding {
	return eventBinding{
		etype: reflect.TypeOf(ArtifactRetrievalSuccessEvent{}),
		handler: func(gevent interface{}) error {
			event, ok := gevent.(ArtifactRetrievalSuccessEvent)
			if ok {
				return handler(event)
			} else {
				return eventmanager.InvalidEventType{errors.New("invalid event type")}
			}
		},
	}
}
