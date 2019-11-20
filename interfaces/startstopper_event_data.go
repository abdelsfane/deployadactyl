package interfaces

import (
	"io"
)

// PushEventData has a RequestBody and DeploymentInfo.
type StartStopEventData struct {
	FoundationURL string
	Context       CFContext
	Courier       interface{}
	Response      io.ReadWriter
}
