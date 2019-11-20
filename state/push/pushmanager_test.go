package push_test

import (
	"bytes"
	"encoding/base64"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/compozed/deployadactyl/constants"
	"github.com/compozed/deployadactyl/controller/deployer"
	"github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/mocks"
	"github.com/compozed/deployadactyl/randomizer"
	. "github.com/compozed/deployadactyl/state/push"
	"github.com/compozed/deployadactyl/structs"
	"github.com/go-errors/errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/op/go-logging"
	"github.com/spf13/afero"
)

var _ = Describe("Actioncreator", func() {
	var (
		logBuffer         *bytes.Buffer
		log               interfaces.Logger
		fetcher           *mocks.Fetcher
		eventManager      *mocks.EventManager
		pusherCreator     *PushManager
		fileSystemCleaner *mocks.FileSystemCleaner
		response          io.ReadWriter
	)

	BeforeEach(func() {
		logBuffer = bytes.NewBuffer([]byte{})
		log = interfaces.DefaultLogger(logBuffer, logging.DEBUG, "deployer tests")

		fetcher = &mocks.Fetcher{}
		eventManager = &mocks.EventManager{}
		fileSystemCleaner = &mocks.FileSystemCleaner{}

		response = NewBuffer()
		pusherCreator = &PushManager{
			Fetcher:      fetcher,
			Logger:       interfaces.DeploymentLogger{log, randomizer.StringRunes(10)},
			EventManager: eventManager,
			DeployEventData: structs.DeployEventData{
				DeploymentInfo: &structs.DeploymentInfo{},
				Response:       response,
			},
			FileSystemCleaner: fileSystemCleaner,
			CFContext:         interfaces.CFContext{},
			Auth:              interfaces.Authorization{},
			Environment:       structs.Environment{Instances: 0},
		}
	})

	Describe("Setup", func() {
		Context("content-type is JSON", func() {

			manifest := `---
applications:
- instances: 2`
			encodedManifest := base64.StdEncoding.EncodeToString([]byte(manifest))

			It("should extract manifest from the request", func() {
				fetcher.FetchCall.Returns.AppPath = "newAppPath"

				deploymentInfo := structs.DeploymentInfo{
					Manifest:    encodedManifest,
					ContentType: "application/json",
				}
				pusherCreator.DeployEventData.DeploymentInfo = &deploymentInfo

				pusherCreator.SetUp()

				Expect(pusherCreator.DeployEventData.DeploymentInfo.Manifest).To(Equal(manifest))
				logBytes, _ := ioutil.ReadAll(logBuffer)
				Eventually(string(logBytes)).Should(ContainSubstring("deploying from json request"))
			})
			It("should fetch and return app path", func() {
				fetcher.FetchCall.Returns.AppPath = "newAppPath"

				deploymentInfo := structs.DeploymentInfo{
					Manifest:    encodedManifest,
					ArtifactURL: "https://artifacturl.com",
					ContentType: "application/json",
				}
				pusherCreator.DeployEventData.DeploymentInfo = &deploymentInfo

				pusherCreator.SetUp()

				Expect(pusherCreator.DeployEventData.DeploymentInfo.AppPath).To(Equal("newAppPath"))
				Expect(fetcher.FetchCall.Received.ArtifactURL).To(Equal(deploymentInfo.ArtifactURL))
				Expect(fetcher.FetchCall.Received.Manifest).To(Equal(manifest))

			})
			It("should error when artifact cannot be fetched", func() {
				fetcher.FetchCall.Returns.Error = errors.New("fetch error")

				deploymentInfo := structs.DeploymentInfo{
					Manifest:    encodedManifest,
					ArtifactURL: "https://artifacturl.com",
					ContentType: "application/json",
				}
				pusherCreator.DeployEventData.DeploymentInfo = &deploymentInfo

				err := pusherCreator.SetUp()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("unzipped app path failed: fetch error"))
			})
			It("should retrieve instances from manifest", func() {
				fetcher.FetchCall.Returns.AppPath = "newAppPath"

				deploymentInfo := structs.DeploymentInfo{
					Manifest:    encodedManifest,
					ContentType: "application/json",
				}
				pusherCreator.DeployEventData.DeploymentInfo = &deploymentInfo

				pusherCreator.SetUp()

				Expect(pusherCreator.DeployEventData.DeploymentInfo.Instances).To(Equal(uint16(2)))
			})
			Context("ArtifactRetrievalStartEvent", func() {
				It("calls EmitEvent", func() {
					fetcher.FetchArtifactFromRequestCall.Returns.Manifest = `---
applications:
- name: "blah"
  instances: 2`
					pusherCreator.SetUp()

					Expect(reflect.TypeOf(eventManager.EmitEventCall.Received.Events[0])).To(Equal(reflect.TypeOf(ArtifactRetrievalStartEvent{})))
				})
				It("passes the CFContext", func() {
					fetcher.FetchArtifactFromRequestCall.Returns.Manifest = `---
applications:
- name: "blah"
  instances: 2`
					pusherCreator.CFContext = interfaces.CFContext{
						Environment:  randomizer.StringRunes(10),
						Space:        randomizer.StringRunes(10),
						Organization: randomizer.StringRunes(10),
						Application:  randomizer.StringRunes(10),
					}

					pusherCreator.SetUp()

					Expect(eventManager.EmitEventCall.Received.Events[0].(ArtifactRetrievalStartEvent).CFContext).To(Equal(pusherCreator.CFContext))
				})
				It("passes the Authorization", func() {
					fetcher.FetchArtifactFromRequestCall.Returns.Manifest = `---
applications:
- name: "blah"
  instances: 2`
					pusherCreator.Auth = interfaces.Authorization{
						Username: randomizer.StringRunes(10),
						Password: randomizer.StringRunes(10),
					}

					pusherCreator.SetUp()

					Expect(eventManager.EmitEventCall.Received.Events[0].(ArtifactRetrievalStartEvent).Auth).To(Equal(pusherCreator.Auth))
				})
				It("passes the environment, response, and data", func() {
					fetcher.FetchArtifactFromRequestCall.Returns.Manifest = `---
applications:
- name: "blah"
  instances: 2`
					environment := structs.Environment{Instances: 0}

					pusherCreator.Environment = environment
					pusherCreator.DeployEventData.Response = response
					pusherCreator.DeployEventData.DeploymentInfo.Data = make(map[string]interface{})

					pusherCreator.SetUp()

					event := eventManager.EmitEventCall.Received.Events[0].(ArtifactRetrievalStartEvent)
					Expect(event.Environment).To(Equal(pusherCreator.Environment))
					Expect(event.Response).To(Equal(pusherCreator.DeployEventData.Response))
					Expect(event.Data).To(Equal(pusherCreator.DeployEventData.DeploymentInfo.Data))
				})
				It("passes the manifest and artifactURL", func() {

					pusherCreator.DeployEventData.DeploymentInfo.Manifest = encodedManifest
					pusherCreator.DeployEventData.DeploymentInfo.ArtifactURL = "theArtifactURL"
					pusherCreator.DeployEventData.DeploymentInfo.ContentType = "application/json"

					pusherCreator.SetUp()

					event := eventManager.EmitEventCall.Received.Events[0].(ArtifactRetrievalStartEvent)
					Expect(event.Manifest).To(Equal(manifest))
					Expect(event.ArtifactURL).To(Equal(pusherCreator.DeployEventData.DeploymentInfo.ArtifactURL))
				})
				Context("if error is returned", func() {
					It("returns an error", func() {

						eventManager.EmitEventCall.Returns.Error = []error{errors.New("a test error")}

						err := pusherCreator.SetUp()

						Expect(err).To(HaveOccurred())
						Expect(reflect.TypeOf(err)).To(Equal(reflect.TypeOf(deployer.EventError{})))
					})
				})

			})

			Context("ArtifactRetrievalSuccessEvent", func() {
				It("calls EmitEvent", func() {
					fetcher.FetchArtifactFromRequestCall.Returns.Manifest = `---
applications:
- name: "blah"
  instances: 2`
					pusherCreator.SetUp()

					Expect(reflect.TypeOf(eventManager.EmitEventCall.Received.Events[1])).To(Equal(reflect.TypeOf(ArtifactRetrievalSuccessEvent{})))
				})
				It("passes the CFContext", func() {
					fetcher.FetchArtifactFromRequestCall.Returns.Manifest = `---
applications:
- name: "blah"
  instances: 2`
					pusherCreator.CFContext = interfaces.CFContext{
						Environment:  randomizer.StringRunes(10),
						Space:        randomizer.StringRunes(10),
						Organization: randomizer.StringRunes(10),
						Application:  randomizer.StringRunes(10),
					}

					pusherCreator.SetUp()

					Expect(eventManager.EmitEventCall.Received.Events[1].(ArtifactRetrievalSuccessEvent).CFContext).To(Equal(pusherCreator.CFContext))
				})
				It("passes the Authorization", func() {
					fetcher.FetchArtifactFromRequestCall.Returns.Manifest = `---
applications:
- name: "blah"
  instances: 2`
					pusherCreator.Auth = interfaces.Authorization{
						Username: randomizer.StringRunes(10),
						Password: randomizer.StringRunes(10),
					}

					pusherCreator.SetUp()

					Expect(eventManager.EmitEventCall.Received.Events[1].(ArtifactRetrievalSuccessEvent).Auth).To(Equal(pusherCreator.Auth))
				})
				It("passes the environment, response, and data", func() {
					fetcher.FetchArtifactFromRequestCall.Returns.Manifest = `---
applications:
- name: "blah"
  instances: 2`
					environment := structs.Environment{Instances: 0}

					pusherCreator.Environment = environment
					pusherCreator.DeployEventData.Response = response
					pusherCreator.DeployEventData.DeploymentInfo.Data = make(map[string]interface{})

					pusherCreator.SetUp()

					event := eventManager.EmitEventCall.Received.Events[1].(ArtifactRetrievalSuccessEvent)
					Expect(event.Environment).To(Equal(pusherCreator.Environment))
					Expect(event.Response).To(Equal(pusherCreator.DeployEventData.Response))
					Expect(event.Data).To(Equal(pusherCreator.DeployEventData.DeploymentInfo.Data))
				})
				It("passes the manifest, artifactURL, and appPath", func() {

					pusherCreator.DeployEventData.DeploymentInfo.Manifest = encodedManifest
					pusherCreator.DeployEventData.DeploymentInfo.ArtifactURL = "theArtifactURL"
					pusherCreator.DeployEventData.DeploymentInfo.ContentType = "application/json"

					fetcher.FetchCall.Returns.AppPath = "new app path"

					pusherCreator.SetUp()

					event := eventManager.EmitEventCall.Received.Events[1].(ArtifactRetrievalSuccessEvent)
					Expect(event.Manifest).To(Equal(manifest))
					Expect(event.ArtifactURL).To(Equal(pusherCreator.DeployEventData.DeploymentInfo.ArtifactURL))
					Expect(event.AppPath).To(Equal("new app path"))
				})
				Context("if error is returned", func() {
					It("returns an error", func() {

						eventManager.EmitEventCall.Returns.Error = []error{nil, errors.New("a test error")}

						err := pusherCreator.SetUp()

						Expect(err).To(HaveOccurred())
						Expect(reflect.TypeOf(err)).To(Equal(reflect.TypeOf(deployer.EventError{})))
						Expect(err.(deployer.EventError).Type).To(Equal(ArtifactRetrievalSuccessEvent{}.Name()))
					})
				})

			})

			Context("ArtifactRetrievalFailureEvent", func() {
				It("calls EmitEvent", func() {

					fetcher.FetchArtifactFromRequestCall.Returns.Error = errors.New("a test error")
					pusherCreator.DeployEventData.DeploymentInfo.ContentType = "application/zip"

					pusherCreator.SetUp()

					Expect(reflect.TypeOf(eventManager.EmitEventCall.Received.Events[1])).To(Equal(reflect.TypeOf(ArtifactRetrievalFailureEvent{})))
				})
				It("passes the CFContext", func() {

					pusherCreator.CFContext = interfaces.CFContext{
						Environment:  randomizer.StringRunes(10),
						Space:        randomizer.StringRunes(10),
						Organization: randomizer.StringRunes(10),
						Application:  randomizer.StringRunes(10),
					}
					fetcher.FetchArtifactFromRequestCall.Returns.Error = errors.New("a test error")
					pusherCreator.DeployEventData.DeploymentInfo.ContentType = "application/zip"

					pusherCreator.SetUp()

					Expect(eventManager.EmitEventCall.Received.Events[1].(ArtifactRetrievalFailureEvent).CFContext).To(Equal(pusherCreator.CFContext))
				})
				It("passes the Authorization", func() {

					pusherCreator.Auth = interfaces.Authorization{
						Username: randomizer.StringRunes(10),
						Password: randomizer.StringRunes(10),
					}
					fetcher.FetchArtifactFromRequestCall.Returns.Error = errors.New("a test error")
					pusherCreator.DeployEventData.DeploymentInfo.ContentType = "application/zip"

					pusherCreator.SetUp()

					Expect(eventManager.EmitEventCall.Received.Events[1].(ArtifactRetrievalFailureEvent).Auth).To(Equal(pusherCreator.Auth))
				})
				It("passes the environment, response, and data", func() {
					environment := structs.Environment{Instances: 0}

					pusherCreator.Environment = environment
					pusherCreator.DeployEventData.Response = response
					pusherCreator.DeployEventData.DeploymentInfo.Data = make(map[string]interface{})

					fetcher.FetchArtifactFromRequestCall.Returns.Error = errors.New("a test error")
					pusherCreator.DeployEventData.DeploymentInfo.ContentType = "application/zip"

					pusherCreator.SetUp()

					event := eventManager.EmitEventCall.Received.Events[1].(ArtifactRetrievalFailureEvent)
					Expect(event.Environment).To(Equal(pusherCreator.Environment))
					Expect(event.Response).To(Equal(pusherCreator.DeployEventData.Response))
					Expect(event.Data).To(Equal(pusherCreator.DeployEventData.DeploymentInfo.Data))
				})
				It("passes the manifest and artifactURL", func() {

					pusherCreator.DeployEventData.DeploymentInfo.Manifest = encodedManifest
					pusherCreator.DeployEventData.DeploymentInfo.ArtifactURL = "theArtifactURL"
					pusherCreator.DeployEventData.DeploymentInfo.ContentType = "application/json"

					fetcher.FetchCall.Returns.Error = errors.New("a test error")

					pusherCreator.SetUp()

					event := eventManager.EmitEventCall.Received.Events[1].(ArtifactRetrievalFailureEvent)
					Expect(event.Manifest).To(Equal(manifest))
					Expect(event.ArtifactURL).To(Equal(pusherCreator.DeployEventData.DeploymentInfo.ArtifactURL))
				})

			})
		})

		Context("when instances is nil", func() {
			It("assigns environmental instances as the instance", func() {
				manifest := `---
applications:
- name: long-running-spring-app`
				encodedManifest := base64.StdEncoding.EncodeToString([]byte(manifest))
				pusherCreator.Environment = structs.Environment{Instances: 22}
				fetcher.FetchCall.Returns.AppPath = "newAppPath"

				deploymentInfo := structs.DeploymentInfo{
					Manifest:    encodedManifest,
					ArtifactURL: "https://artifacturl.com",
					ContentType: "application/json",
				}
				pusherCreator.DeployEventData.DeploymentInfo = &deploymentInfo

				pusherCreator.SetUp()

				Expect(pusherCreator.DeployEventData.DeploymentInfo.Instances).To(Equal(uint16(22)))
			})
		})

		Context("contentType is ZIP", func() {

			It("should extract manifest from the zip file", func() {
				fetcher.FetchArtifactFromRequestCall.Returns.Manifest = `---
applications:
- name: "blah"
  instances: 2`
				fetcher.FetchArtifactFromRequestCall.Returns.AppPath = "newAppPath"

				deploymentInfo := structs.DeploymentInfo{ContentType: "application/zip"}
				pusherCreator.DeployEventData.DeploymentInfo = &deploymentInfo

				pusherCreator.SetUp()

				Expect(pusherCreator.DeployEventData.DeploymentInfo.AppPath).To(Equal("newAppPath"))
				logBytes, _ := ioutil.ReadAll(logBuffer)
				Eventually(string(logBytes)).Should(ContainSubstring("deploying from archive request"))
			})

			Context("when instances are listed in the manifest", func() {
				It("should set the correct number of instances from the manifest", func() {
					fetcher.FetchArtifactFromRequestCall.Returns.AppPath = "newAppPath"
					fetcher.FetchArtifactFromRequestCall.Returns.Manifest = `---
applications:
- name: "blah"
  instances: 2
`

					deploymentInfo := structs.DeploymentInfo{ContentType: "application/zip"}
					pusherCreator.DeployEventData.DeploymentInfo = &deploymentInfo

					pusherCreator.SetUp()

					Expect(pusherCreator.DeployEventData.DeploymentInfo.Instances).To(Equal(uint16(2)))
					logBytes, _ := ioutil.ReadAll(logBuffer)
					Eventually(string(logBytes)).Should(ContainSubstring("deploying from archive request"))
				})
			})
			Context("when no instances are listed in the manifest", func() {
				It("should set the correct number of instances from the environment", func() {
					fetcher.FetchArtifactFromRequestCall.Returns.AppPath = "newAppPath"
					pusherCreator.Environment.Instances = 7
					fetcher.FetchArtifactFromRequestCall.Returns.Manifest = `---
applications:
- name: "blah"
`

					deploymentInfo := structs.DeploymentInfo{ContentType: "application/zip"}
					pusherCreator.DeployEventData.DeploymentInfo = &deploymentInfo

					pusherCreator.SetUp()

					Expect(pusherCreator.DeployEventData.DeploymentInfo.Instances).To(Equal(uint16(7)))
					logBytes, _ := ioutil.ReadAll(logBuffer)
					Eventually(string(logBytes)).Should(ContainSubstring("deploying from archive request"))
				})
			})

			It("should error when artifact cannot be fetched", func() {
				fetcher.FetchArtifactFromRequestCall.Returns.Error = errors.New("a test error")

				deploymentInfo := structs.DeploymentInfo{
					ContentType: "application/zip",
				}
				pusherCreator.DeployEventData.DeploymentInfo = &deploymentInfo

				err := pusherCreator.SetUp()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("unzipping request body error: a test error"))
			})
		})

	})

	Describe("OnStart", func() {
		Context("push.started Emit", func() {
			It("emits a push.started event", func() {
				pusherCreator.OnStart()

				Expect(eventManager.EmitCall.Received.Events[0].Type).Should(Equal(constants.PushStartedEvent))
			})
			It("logs the parameters", func() {
				deployInfo := pusherCreator.DeployEventData.DeploymentInfo
				deployInfo.ArtifactURL = randomizer.StringRunes(10)
				deployInfo.Username = randomizer.StringRunes(10)
				deployInfo.Environment = randomizer.StringRunes(10)
				deployInfo.Org = randomizer.StringRunes(10)
				deployInfo.Space = randomizer.StringRunes(10)
				deployInfo.AppName = randomizer.StringRunes(10)

				pusherCreator.OnStart()

				logBytes, _ := ioutil.ReadAll(logBuffer)
				Eventually(string(logBytes)).Should(ContainSubstring("Artifact URL: " + pusherCreator.DeployEventData.DeploymentInfo.ArtifactURL))
				Eventually(string(logBytes)).Should(ContainSubstring("Username:     " + pusherCreator.DeployEventData.DeploymentInfo.Username))
				Eventually(string(logBytes)).Should(ContainSubstring("Environment:  " + pusherCreator.DeployEventData.DeploymentInfo.Environment))
				Eventually(string(logBytes)).Should(ContainSubstring("Org:          " + pusherCreator.DeployEventData.DeploymentInfo.Org))
				Eventually(string(logBytes)).Should(ContainSubstring("Space:        " + pusherCreator.DeployEventData.DeploymentInfo.Space))
				Eventually(string(logBytes)).Should(ContainSubstring("AppName:      " + pusherCreator.DeployEventData.DeploymentInfo.AppName))
			})
			It("prints the parameters to the response", func() {
				deployInfo := pusherCreator.DeployEventData.DeploymentInfo
				deployInfo.ArtifactURL = randomizer.StringRunes(10)
				deployInfo.Username = randomizer.StringRunes(10)
				deployInfo.Environment = randomizer.StringRunes(10)
				deployInfo.Org = randomizer.StringRunes(10)
				deployInfo.Space = randomizer.StringRunes(10)
				deployInfo.AppName = randomizer.StringRunes(10)

				pusherCreator.OnStart()

				Eventually(response).Should(Say("Artifact URL: " + pusherCreator.DeployEventData.DeploymentInfo.ArtifactURL))
				Eventually(response).Should(Say("Username:     " + pusherCreator.DeployEventData.DeploymentInfo.Username))
				Eventually(response).Should(Say("Environment:  " + pusherCreator.DeployEventData.DeploymentInfo.Environment))
				Eventually(response).Should(Say("Org:          " + pusherCreator.DeployEventData.DeploymentInfo.Org))
				Eventually(response).Should(Say("Space:        " + pusherCreator.DeployEventData.DeploymentInfo.Space))
				Eventually(response).Should(Say("AppName:      " + pusherCreator.DeployEventData.DeploymentInfo.AppName))
			})
			Context("if Emit fails", func() {
				It("returns an error", func() {
					eventManager.EmitCall.Returns.Error = []error{errors.New("a test error")}

					err := pusherCreator.OnStart()

					Expect(reflect.TypeOf(err)).Should(Equal(reflect.TypeOf(deployer.EventError{})))
				})
				It("logs the error", func() {
					eventManager.EmitCall.Returns.Error = []error{errors.New("a test error")}

					pusherCreator.OnStart()

					logBytes, _ := ioutil.ReadAll(logBuffer)
					Eventually(string(logBytes)).Should(ContainSubstring("a test error"))
				})
			})
		})
		Context("calls EmitEvent", func() {
			It("emits a PushStartedEvent", func() {
				pusherCreator.OnStart()

				Expect(reflect.TypeOf(eventManager.EmitEventCall.Received.Events[0])).Should(Equal(reflect.TypeOf(PushStartedEvent{})))
			})
			It("passes CFContext", func() {
				pusherCreator.CFContext = interfaces.CFContext{
					Environment:  randomizer.StringRunes(10),
					Organization: randomizer.StringRunes(10),
					Space:        randomizer.StringRunes(10),
					Application:  randomizer.StringRunes(10),
				}

				pusherCreator.OnStart()

				event := eventManager.EmitEventCall.Received.Events[0].(PushStartedEvent)
				Expect(event.CFContext).To(Equal(pusherCreator.CFContext))
			})
			It("passes Authorization", func() {
				pusherCreator.Auth = interfaces.Authorization{
					Username: randomizer.StringRunes(10),
					Password: randomizer.StringRunes(10),
				}

				pusherCreator.OnStart()

				event := eventManager.EmitEventCall.Received.Events[0].(PushStartedEvent)
				Expect(event.Auth).To(Equal(pusherCreator.Auth))
			})
			It("passes Environment", func() {
				pusherCreator.Environment = structs.Environment{
					Name: randomizer.StringRunes(10),
				}

				pusherCreator.OnStart()

				event := eventManager.EmitEventCall.Received.Events[0].(PushStartedEvent)
				Expect(event.Environment).To(Equal(pusherCreator.Environment))
			})
			It("passes other params", func() {
				pusherCreator.DeployEventData.Response = response
				pusherCreator.DeployEventData.DeploymentInfo.Body = bytes.NewBuffer([]byte{})
				pusherCreator.DeployEventData.DeploymentInfo.ContentType = "content-type"

				pusherCreator.OnStart()

				event := eventManager.EmitEventCall.Received.Events[0].(PushStartedEvent)
				Expect(event.Response).To(Equal(pusherCreator.DeployEventData.Response))
				Expect(event.Body).To(Equal(pusherCreator.DeployEventData.DeploymentInfo.Body))
				Expect(event.ContentType).To(Equal("content-type"))
			})
			Context("if EmitEvent fails", func() {
				It("returns an error", func() {
					eventManager.EmitEventCall.Returns.Error = []error{errors.New("a test error")}

					err := pusherCreator.OnStart()

					Expect(reflect.TypeOf(err)).Should(Equal(reflect.TypeOf(deployer.EventError{})))
				})
				It("logs the error", func() {
					eventManager.EmitEventCall.Returns.Error = []error{errors.New("a test error")}

					pusherCreator.OnStart()

					logBytes, _ := ioutil.ReadAll(logBuffer)
					Eventually(string(logBytes)).Should(ContainSubstring("a test error"))
				})
			})
		})

	})

	Describe("CleanUp", func() {
		It("deletes all temp artifacts", func() {
			path := randomizer.StringRunes(10)
			pusherCreator.DeployEventData.DeploymentInfo.AppPath = path

			pusherCreator.CleanUp()

			Expect(fileSystemCleaner.RemoveAllCall.Received.Path).To(Equal(path))
		})
		It("really deletes all temp artifacts", func() {
			af := &afero.Afero{Fs: afero.NewMemMapFs()}
			pusherCreator.FileSystemCleaner = af

			directoryName, _ := af.TempDir("", "deployadactyl-")

			pusherCreator.CleanUp()

			exists, err := af.DirExists(directoryName)
			Expect(err).ToNot(HaveOccurred())

			Expect(exists).ToNot(BeTrue())
		})
	})

	Describe("OnFinish", func() {
		Context("when error occurs", func() {
			Context("and DisableRollback is true", func() {
				It("returns StatusOK", func() {
					env := structs.Environment{DisableRollback: true}
					err := errors.New("a test error")

					resp := pusherCreator.OnFinish(env, response, err)

					Expect(resp.StatusCode).To(Equal(http.StatusOK))
				})
				It("logs the failure", func() {
					env := structs.Environment{DisableRollback: true}
					err := errors.New("a test error")

					pusherCreator.OnFinish(env, response, err)

					logBytes, _ := ioutil.ReadAll(logBuffer)
					Eventually(string(logBytes)).Should(ContainSubstring("DisabledRollback true, returning status"))
				})
			})
			Context("and DisableRollback is false", func() {
				Context("and error is a login failure", func() {
					It("returns StatusBadRequest", func() {
						env := structs.Environment{DisableRollback: false}
						err := errors.New("the login failed")

						resp := pusherCreator.OnFinish(env, response, err)

						Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
					})
				})
				It("returns StatusInternalServerError", func() {
					env := structs.Environment{DisableRollback: false}
					err := errors.New("a test error")

					resp := pusherCreator.OnFinish(env, response, err)

					Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
				})
			})
		})
		Context("when no error occurs", func() {
			It("returns StatusOK", func() {
				env := structs.Environment{DisableRollback: false}

				resp := pusherCreator.OnFinish(env, response, nil)

				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})
			It("logs a successful deployment message", func() {
				env := structs.Environment{DisableRollback: false}

				pusherCreator.OnFinish(env, response, nil)
				logBytes, _ := ioutil.ReadAll(logBuffer)
				Eventually(string(logBytes)).Should(ContainSubstring("successfully deployed application"))
			})
			It("writes success to the output", func() {
				env := structs.Environment{DisableRollback: false}

				pusherCreator.OnFinish(env, response, nil)
				logBytes, _ := ioutil.ReadAll(response)
				Eventually(string(logBytes)).Should(ContainSubstring("Your deploy was successful!"))
			})
		})
	})
})
