package envvar_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/op/go-logging"
	"github.com/spf13/afero"

	. "github.com/compozed/deployadactyl/eventmanager/handlers/envvar"
	I "github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/state/push"
)

var _ = Describe("Env_Var_Handler", func() {
	var (
		eventHandler Envvarhandler
		logBuffer    *gbytes.Buffer
		log          I.DeploymentLogger
		ievent       push.ArtifactRetrievalSuccessEvent
		filesystem   = &afero.Afero{Fs: afero.NewMemMapFs()}
	)

	BeforeEach(func() {
		logBuffer = gbytes.NewBuffer()
		log = I.DeploymentLogger{Log: I.DefaultLogger(logBuffer, logging.DEBUG, "evn_var_handler_test")}
		ievent = push.ArtifactRetrievalSuccessEvent{
			Log: log,
		}
		eventHandler = Envvarhandler{FileSystem: filesystem}
	})

	Context("when an envvarhandler is called with event without deploy info", func() {
		It("it should succeed", func() {

			Expect(eventHandler.ArtifactRetrievalSuccessEventHandler(ievent)).To(Succeed())
		})
	})

	Context("when an envvarhandler is called with event without env variables", func() {
		It("it should be succeed", func() {

			ievent.EnvironmentVariables = nil

			Expect(eventHandler.ArtifactRetrievalSuccessEventHandler(ievent)).To(Succeed())
		})
	})

	Context("when an envvarhandler is called with event with env variables", func() {
		It("it should be succeed", func() {

			path := "/tmp"
			eventHandler.FileSystem.MkdirAll(path, 0755)

			envvars := make(map[string]string)
			envvars["one"] = "one"
			envvars["two"] = "two"

			ievent.AppPath = path
			ievent.EnvironmentVariables = envvars
			ievent.CFContext = I.CFContext{
				Application: "testApp",
			}

			//Process the event
			Expect(eventHandler.ArtifactRetrievalSuccessEventHandler(ievent)).To(Succeed())

			//Verify manifest was written and matches
			manifest, err := ReadManifest(path+"/manifest.yml", log, eventHandler.FileSystem)

			Expect(manifest).NotTo(BeNil())
			Expect(err).To(BeNil())
			Expect(manifest.Content.Applications[0].Name).To(Equal("testApp"))
			Expect(len(manifest.Content.Applications[0].Env)).To(Equal(2))
		})
	})

	Context("when an envvarhandler is called with bogus manifest in deploy info", func() {
		It("it should be fail", func() {

			content := `bork`

			envvars := make(map[string]string)
			envvars["one"] = "one"
			envvars["two"] = "two"

			ievent.AppPath = "/tmp"
			ievent.Manifest = content
			ievent.EnvironmentVariables = envvars
			ievent.CFContext = I.CFContext{
				Application: "testApp",
			}

			err := eventHandler.ArtifactRetrievalSuccessEventHandler(ievent)

			Expect(err).ToNot(BeNil())
		})
	})
})
