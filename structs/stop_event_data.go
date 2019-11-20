package structs

import "io"

type StopEventData struct {
	Response       io.ReadWriter
	DeploymentInfo *DeploymentInfo
}
