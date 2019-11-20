package envvar

import (
	"github.com/spf13/afero"

	"github.com/compozed/deployadactyl/state/push"
)

type Envvarhandler struct {
	//Logger     I.Logger
	FileSystem *afero.Afero
}

func (handler Envvarhandler) ArtifactRetrievalSuccessEventHandler(event push.ArtifactRetrievalSuccessEvent) error {

	event.Log.Debugf("Environment Variable Handler Processing Event => %+v", event)

	if event.EnvironmentVariables == nil || len(event.EnvironmentVariables) == 0 {
		event.Log.Info("No Deployment Info or Environment Variables to process!")
		return nil
	}

	m, err := CreateManifest(event.CFContext.Application, event.Manifest, handler.FileSystem, event.Log)

	if err != nil {
		event.Log.Errorf("Error Parsing Manifest! Details: %v", err)
		return err
	}

	//Add any Environment variables
	addEnvResult, _ := m.AddEnvironmentVariables(event.EnvironmentVariables)

	if m.Content.Applications[0].Path != "" || addEnvResult {

		//Ensure path is empty. We are using a local/tmp file system with exploded contents for the deploy!
		m.Content.Applications[0].Path = ""

		//Re-Write the m
		m.WriteManifest(event.AppPath, true)
	}

	return nil
}
