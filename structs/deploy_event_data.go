package structs

import "io"

// DeployEventData has a RequestBody and DeploymentInfo.
type DeployEventData struct {
	// Writer is being deprecated in favor of using Response as a ReadWriter. 01/03/2017
	Writer io.Writer

	Response       io.ReadWriter
	DeploymentInfo *DeploymentInfo
	RequestBody    io.Reader
}
