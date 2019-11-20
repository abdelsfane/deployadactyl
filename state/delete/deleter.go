package delete

import (
	"io"

	"fmt"

	I "github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/state"
)

type Deleter struct {
	Courier       I.Courier
	CFContext     I.CFContext
	Authorization I.Authorization
	EventManager  I.EventManager
	Response      io.ReadWriter
	Log           I.DeploymentLogger
	FoundationURL string
	AppName       string
}

func (s Deleter) Verify() error {
	return nil
}

func (s Deleter) Success() error {
	return nil
}

func (s Deleter) Finally() error {
	return nil
}

// Login will login to a Cloud Foundry instance.
func (s Deleter) Initially() error {
	s.Log.Debugf(
		`logging into cloud foundry with parameters:
		foundation URL: %+v
		username: %+v
		org: %+v
		space: %+v`,
		s.FoundationURL, s.Authorization.Username, s.CFContext.Organization, s.CFContext.Space,
	)

	output, err := s.Courier.Login(
		s.FoundationURL,
		s.Authorization.Username,
		s.Authorization.Password,
		s.CFContext.Organization,
		s.CFContext.Space,
		s.CFContext.SkipSSL,
	)
	s.Response.Write(output)
	if err != nil {
		s.Log.Errorf("could not login to %s", s.FoundationURL)
		return state.LoginError{s.FoundationURL, output}
	}

	s.Log.Infof("logged into cloud foundry %s", s.FoundationURL)

	return nil
}

func (s Deleter) Execute() error {

	if s.Courier.Exists(s.AppName) != true {
		s.Log.Errorf("failed to delete app on foundation %s: application doesn't exist", s.FoundationURL)
		return state.ExistsError{ApplicationName: s.AppName}
	}

	s.Log.Infof("%s: deleting app %s", s.FoundationURL, s.AppName)

	output, err := s.Courier.Delete(s.AppName)
	if err != nil {
		s.Log.Errorf("failed to delete app on foundation %s: %s", s.FoundationURL, err.Error())
		return state.DeleteError{ApplicationName: s.AppName, Out: output}
	}
	s.Response.Write(output)

	s.Log.Infof("%s: successfully deleted app %s", s.FoundationURL, s.AppName)

	return nil
}

func (s Deleter) PostExecute() error {
	return nil
}

func (s Deleter) Undo() error {
	s.Response.Write([]byte(fmt.Sprintf("delete feature is unable to rollback: %s", s.AppName)))
	s.Log.Infof("%s: delete feature is unable to rollback: %s", s.FoundationURL, s.AppName)

	return nil
}
