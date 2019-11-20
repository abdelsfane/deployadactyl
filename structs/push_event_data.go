package structs

import (
	"io"
)

// PushEventData has a RequestBody and DeploymentInfo.
type PushEventData struct {
	AppPath         string
	FoundationURL   string
	TempAppWithUUID string

	DeploymentInfo *DeploymentInfo
	Courier        interface{}
	Response       io.ReadWriter
}
