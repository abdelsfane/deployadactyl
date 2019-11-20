package state

import "fmt"

type CloudFoundryGetLogsError struct {
	CfTaskErr error
	CfLogErr  error
}

func (e CloudFoundryGetLogsError) Error() string {
	return fmt.Sprintf("%s: cannot get Cloud Foundry logs: %s", e.CfTaskErr, e.CfLogErr)
}

type DeleteApplicationError struct {
	ApplicationName string
	Out             []byte
}

func (e DeleteApplicationError) Error() string {
	return fmt.Sprintf("cannot delete %s: %s", e.ApplicationName, string(e.Out))
}

type LoginError struct {
	FoundationURL string
	Out           []byte
}

func (e LoginError) Error() string {
	return fmt.Sprintf("cannot login to %s: %s", e.FoundationURL, string(e.Out))
}

type RenameError struct {
	ApplicationName string
	Out             []byte
}

func (e RenameError) Error() string {
	return fmt.Sprintf("cannot rename %s: %s", e.ApplicationName, string(e.Out))
}

type PushError struct{}

func (e PushError) Error() string {
	return "check the Cloud Foundry output above for more information"
}

type MapRouteError struct {
	Out []byte
}

func (e MapRouteError) Error() string {
	return fmt.Sprintf("map route failed: %s", string(e.Out))
}

type UnmapRouteError struct {
	ApplicationName string
	Out             []byte
}

func (e UnmapRouteError) Error() string {
	return fmt.Sprintf("failed to unmap route for %s: %s", e.ApplicationName, string(e.Out))
}

type InvalidContentTypeError struct{}

func (e InvalidContentTypeError) Error() string {
	return "must be application/json or application/zip"
}

type AppPathError struct {
	Err error
}

func (e AppPathError) Error() string {
	return fmt.Sprintf("unzipped app path failed: %s", e.Err)
}

type ManifestError struct{}

func (e ManifestError) Error() string {
	return "manifest decoding error"
}

type UnzippingError struct {
	Err error
}

func (e UnzippingError) Error() string {
	return fmt.Sprintf("unzipping request body error: %s", e.Err)
}

type CourierCreationError struct {
	Err error
}

func (e CourierCreationError) Error() string {
	return fmt.Sprintf("failed to create Courier: %s", e.Err.Error())
}

type StartError struct {
	ApplicationName string
	Out             []byte
}

func (e StartError) Error() string {
	return fmt.Sprintf("cannot start %s: %s", e.ApplicationName, string(e.Out))
}

type StopError struct {
	ApplicationName string
	Out             []byte
}

func (e StopError) Error() string {
	return fmt.Sprintf("cannot stop %s: %s", e.ApplicationName, string(e.Out))
}

type DeleteError struct {
	ApplicationName string
	Out             []byte
}

func (e DeleteError) Error() string {
	return fmt.Sprintf("cannot delete %s: %s", e.ApplicationName, string(e.Out))
}

type ExistsError struct {
	ApplicationName string
}

func (e ExistsError) Error() string {
	return fmt.Sprintf("app %s doesn't exist", e.ApplicationName)
}
