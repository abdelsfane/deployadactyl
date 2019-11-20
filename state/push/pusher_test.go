package push_test

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/compozed/deployadactyl/mocks"
	"github.com/compozed/deployadactyl/randomizer"
	. "github.com/compozed/deployadactyl/state/push"
	S "github.com/compozed/deployadactyl/structs"
	"github.com/op/go-logging"

	"github.com/compozed/deployadactyl/eventmanager/handlers/healthchecker"
	"github.com/compozed/deployadactyl/eventmanager/handlers/routemapper"
	"github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/state"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/spf13/afero"
)

var _ = Describe("Pusher", func() {
	var (
		pusher       Pusher
		courier      *mocks.Courier
		eventManager *mocks.EventManager
		fetcher      *mocks.Fetcher
		client       *mocks.Client

		randomUsername      string
		randomPassword      string
		randomOrg           string
		randomSpace         string
		randomDomain        string
		randomAppPath       string
		randomAppName       string
		randomInstances     uint16
		randomUUID          string
		randomEndpoint      string
		randomFoundationURL string
		randomArtifactUrl   string
		tempAppWithUUID     string
		randomHostName      string
		skipSSL             bool
		deploymentInfo      S.DeploymentInfo
		response            *Buffer
		logBuffer           *Buffer
	)

	BeforeEach(func() {
		courier = &mocks.Courier{}
		eventManager = &mocks.EventManager{}
		fetcher = &mocks.Fetcher{}
		client = &mocks.Client{}

		randomFoundationURL = "randomFoundationURL-" + randomizer.StringRunes(10)
		randomUsername = "randomUsername-" + randomizer.StringRunes(10)
		randomPassword = "randomPassword-" + randomizer.StringRunes(10)
		randomOrg = "randomOrg-" + randomizer.StringRunes(10)
		randomSpace = "randomSpace-" + randomizer.StringRunes(10)
		randomDomain = "randomDomain-" + randomizer.StringRunes(10)
		randomAppPath = "randomAppPath-" + randomizer.StringRunes(10)
		randomAppName = "randomAppName-" + randomizer.StringRunes(10)
		randomEndpoint = "randomEndpoint-" + randomizer.StringRunes(10)
		randomArtifactUrl = "randomArtifactUrl-" + randomizer.StringRunes(10)
		randomUUID = randomizer.StringRunes(10)
		randomInstances = uint16(rand.Uint32())
		randomHostName = "randomHostName" + randomizer.StringRunes(10)

		tempAppWithUUID = randomAppName + TemporaryNameSuffix + randomUUID

		manifest := fmt.Sprintf(`
---
applications:
- name: example
  custom-routes:
  - route: %s
  - route: %s
  - route: %s

  env:
    CONVEYOR: 23432`,
			fmt.Sprintf("%s0.%s0", randomAppName, randomDomain),
			fmt.Sprintf("%s1.%s1", randomAppName, randomDomain),
			fmt.Sprintf("%s2.%s2", randomAppName, randomDomain),
		)

		response = NewBuffer()
		logBuffer = NewBuffer()

		eventManager.EmitCall.Returns.Error = append(eventManager.EmitCall.Returns.Error, nil)
		client.GetCall.Returns.Response.StatusCode = 200
		courier.DomainsCall.Returns.Domains = []string{randomDomain + "0", randomDomain + "1", randomDomain + "2"}

		healthChecker := healthchecker.HealthChecker{
			Client: client,
			OldURL: "api.cf",
			NewURL: "apps",
		}

		routeMapper := routemapper.RouteMapper{
			Courier:    courier,
			FileSystem: &afero.Afero{Fs: afero.NewMemMapFs()},
		}

		deploymentInfo = S.DeploymentInfo{
			Username:            randomUsername,
			Password:            randomPassword,
			Org:                 randomOrg,
			Space:               randomSpace,
			AppName:             randomAppName,
			SkipSSL:             skipSSL,
			Instances:           randomInstances,
			Domain:              randomDomain,
			UUID:                randomUUID,
			HealthCheckEndpoint: randomEndpoint,
			ArtifactURL:         randomArtifactUrl,
			Manifest:            manifest,
			ContentType:         "JSON",
			Data:                make(map[string]interface{}, 0),
		}

		pusher = Pusher{
			Courier:        courier,
			DeploymentInfo: deploymentInfo,
			EventManager:   eventManager,
			Response:       response,
			Log:            interfaces.DeploymentLogger{Log: interfaces.DefaultLogger(logBuffer, logging.DEBUG, "pusher_test")},
			FoundationURL:  randomFoundationURL,
			AppPath:        randomAppPath,
			Environment:    S.Environment{DisableRollback: false},
			Fetcher:        fetcher,
			CFContext:      interfaces.CFContext{},
			Auth:           interfaces.Authorization{},
			HealthChecker:  healthChecker,
			RouteMapper:    routeMapper,
		}
	})

	Describe("Initially", func() {
		Context("when login succeeds", func() {
			It("gives the correct info to the courier", func() {

				Expect(pusher.Initially()).To(Succeed())

				Expect(courier.LoginCall.Received.FoundationURL).To(Equal(randomFoundationURL))
				Expect(courier.LoginCall.Received.Username).To(Equal(randomUsername))
				Expect(courier.LoginCall.Received.Password).To(Equal(randomPassword))
				Expect(courier.LoginCall.Received.Org).To(Equal(randomOrg))
				Expect(courier.LoginCall.Received.Space).To(Equal(randomSpace))
				Expect(courier.LoginCall.Received.SkipSSL).To(Equal(skipSSL))
			})

			It("writes the output of the courier to the response", func() {
				courier.LoginCall.Returns.Output = []byte("login succeeded")

				Expect(pusher.Initially()).To(Succeed())

				Eventually(response).Should(Say("login succeeded"))
			})
		})

		Context("when login fails", func() {
			It("returns an error", func() {
				courier.LoginCall.Returns.Output = []byte("login output")
				courier.LoginCall.Returns.Error = errors.New("login error")

				err := pusher.Initially()
				Expect(err).To(MatchError(state.LoginError{randomFoundationURL, []byte("login output")}))
			})

			It("writes the output of the courier to the response", func() {
				courier.LoginCall.Returns.Output = []byte("login output")
				courier.LoginCall.Returns.Error = errors.New("login error")

				err := pusher.Initially()
				Expect(err).To(HaveOccurred())

				Eventually(response).Should(Say("login output"))
			})

			It("logs an error", func() {
				courier.LoginCall.Returns.Error = errors.New("login error")

				err := pusher.Initially()
				Expect(err).To(HaveOccurred())

				Eventually(logBuffer).Should(Say(fmt.Sprintf("could not login to %s", randomFoundationURL)))
			})
		})
	})

	Describe("Execute", func() {
		Context("with JSON request body", func() {
			Context("when the push succeeds", func() {
				It("pushes the new app", func() {
					courier.PushCall.Returns.Output = []byte("push succeeded")

					Expect(pusher.Execute()).To(Succeed())

					Expect(courier.PushCall.Received.AppName).To(Equal(tempAppWithUUID))
					Expect(courier.PushCall.Received.AppPath).To(Equal(randomAppPath))
					Expect(courier.PushCall.Received.Hostname).To(Equal(randomAppName))
					Expect(courier.PushCall.Received.Instances).To(Equal(randomInstances))

					Eventually(response).Should(Say("push succeeded"))

					Eventually(logBuffer).Should(Say(fmt.Sprintf("pushing app %s to %s", tempAppWithUUID, randomDomain)))
					Eventually(logBuffer).Should(Say(fmt.Sprintf("tempdir for app %s: %s", tempAppWithUUID, randomAppPath)))
					Eventually(logBuffer).Should(Say("output from Cloud Foundry"))
					Eventually(logBuffer).Should(Say("successfully deployed new build"))
				})
			})

			Context("when the push fails", func() {
				It("returns an error", func() {
					fetcher.FetchCall.Returns.AppPath = randomAppPath
					courier.PushCall.Returns.Error = errors.New("push error")

					err := pusher.Execute()

					Expect(err).To(MatchError(state.PushError{}))
				})

				It("gets logs from the courier", func() {
					fetcher.FetchCall.Returns.AppPath = randomAppPath
					courier.PushCall.Returns.Output = []byte("push output")
					courier.PushCall.Returns.Error = errors.New("push error")
					courier.LogsCall.Returns.Output = []byte("cf logs")

					Expect(pusher.Execute()).ToNot(Succeed())

					Eventually(response).Should(Say("push output"))
					Eventually(response).Should(Say("cf logs"))

					Eventually(logBuffer).Should(Say("logs from"))
				})

				Context("when the courier log call fails", func() {
					It("returns an error", func() {
						fetcher.FetchCall.Returns.AppPath = randomAppPath
						pushErr := errors.New("push error")
						logsErr := errors.New("logs error")

						courier.PushCall.Returns.Error = pushErr
						courier.LogsCall.Returns.Error = logsErr

						err := pusher.Execute()

						Expect(err).To(MatchError(state.CloudFoundryGetLogsError{pushErr, logsErr}))
					})
				})
			})
		})

		Context("with Zip request body", func() {
			Context("when the push succeeds", func() {
				It("pushes the new app", func() {
					pusher.DeploymentInfo.ContentType = "ZIP"
					courier.PushCall.Returns.Output = []byte("push succeeded")
					fetcher.FetchArtifactFromRequestCall.Returns.AppPath = randomAppPath

					Expect(pusher.Execute()).To(Succeed())

					Expect(courier.PushCall.Received.AppName).To(Equal(tempAppWithUUID))
					Expect(courier.PushCall.Received.AppPath).To(Equal(randomAppPath))
					Expect(courier.PushCall.Received.Hostname).To(Equal(randomAppName))
					Expect(courier.PushCall.Received.Instances).To(Equal(randomInstances))

					Eventually(response).Should(Say("push succeeded"))

					Eventually(logBuffer).Should(Say(fmt.Sprintf("pushing app %s to %s", tempAppWithUUID, randomDomain)))
					Eventually(logBuffer).Should(Say(fmt.Sprintf("tempdir for app %s: %s", tempAppWithUUID, randomAppPath)))
					Eventually(logBuffer).Should(Say("output from Cloud Foundry"))
					Eventually(logBuffer).Should(Say("successfully deployed new build"))
				})
			})
		})

		Context("with other besides zip and json request body type", func() {
			Context("when the push succeeds", func() {
				It("pushes the new app", func() {
					pusher.DeploymentInfo.ContentType = "ZIP"
					courier.PushCall.Returns.Output = []byte("push succeeded")
					fetcher.FetchArtifactFromRequestCall.Returns.AppPath = randomAppPath

					Expect(pusher.Execute()).To(Succeed())

					Expect(courier.PushCall.Received.AppName).To(Equal(tempAppWithUUID))
					Expect(courier.PushCall.Received.AppPath).To(Equal(randomAppPath))
					Expect(courier.PushCall.Received.Hostname).To(Equal(randomAppName))
					Expect(courier.PushCall.Received.Instances).To(Equal(randomInstances))

					Eventually(response).Should(Say("push succeeded"))

					Eventually(logBuffer).Should(Say(fmt.Sprintf("pushing app %s to %s", tempAppWithUUID, randomDomain)))
					Eventually(logBuffer).Should(Say(fmt.Sprintf("tempdir for app %s: %s", tempAppWithUUID, randomAppPath)))
					Eventually(logBuffer).Should(Say("output from Cloud Foundry"))
					Eventually(logBuffer).Should(Say("successfully deployed new build"))
				})
			})
		})

		It("should write foundation URL to log", func() {
			pusher.Execute()
			Eventually(logBuffer).Should(Say(randomFoundationURL + ": pushing app"))
			Eventually(logBuffer).Should(Say(randomFoundationURL + ": tempdir for app"))
			Eventually(logBuffer).Should(Say(randomFoundationURL + ": push output from Cloud Foundry"))
			Eventually(logBuffer).Should(Say(randomFoundationURL + ": successfully deployed"))
		})

		Context("when healthcheck route is created", func() {
			It("should call the healthchecker", func() {
				courier.PushCall.Returns.Output = []byte("push succeeded")

				Expect(pusher.Execute()).To(Succeed())

				Expect(client.GetCall.Received.URL).To(ContainSubstring(randomFoundationURL))
				Expect(client.GetCall.Received.URL).To(ContainSubstring(randomEndpoint))
			})
		})

		Context("when push call fails", func() {
			It("should write foundation URL to log", func() {
				courier.PushCall.Returns.Error = errors.New("an error")
				pusher.Execute()
				Eventually(logBuffer).Should(Say(randomFoundationURL + ": logs from"))
			})
		})
	})

	Describe("PostExecute", func() {
		It("calls MapRoute with the correct input", func() {
			courier.ExistsCall.Returns.Bool = true
			pusher.PostExecute()

			Expect(courier.MapRouteCall.Received.Hostname[0]).To(Equal(randomAppName + "0"))
			Expect(courier.MapRouteCall.Received.AppName[0]).To(Equal(randomAppName + "-new-build-" + randomUUID))
			Expect(courier.MapRouteCall.Received.Domain[0]).To(Equal(randomDomain + "0"))
		})

		It("writes logs about mapping routes", func() {
			courier.ExistsCall.Returns.Bool = true
			pusher.PostExecute()

			Eventually(logBuffer).Should(Say("mapping route for"))
			Eventually(logBuffer).Should(Say("application route created"))
		})

		Context("when no domain is set on the pusher struct", func() {
			It("calls MapRoute with the correct input", func() {
				courier.ExistsCall.Returns.Bool = true
				pusher.DeploymentInfo.Domain = ""
				pusher.PostExecute()

				Expect(courier.MapRouteCall.TimesCalled).To(Equal(3))
			})
		})

		It("calls CustomRouteMapper with correct input", func() {
			manifest := fmt.Sprintf(`
---
applications:
- name: example
  custom-routes:
  - route: %s
  - route: %s
  - route: %s

  env:
    CONVEYOR: 23432`,
				fmt.Sprintf("%s0.%s0", randomHostName, randomDomain),
				fmt.Sprintf("%s1.%s1", randomHostName, randomDomain),
				fmt.Sprintf("%s2.%s2", randomHostName, randomDomain),
			)

			pusher.DeploymentInfo.Manifest = manifest
			courier.ExistsCall.Returns.Bool = true
			courier.DomainsCall.Returns.Domains = []string{randomDomain + "0", randomDomain + "1", randomDomain + "2"}

			pusher.PostExecute()

			Expect(courier.DomainsCall.TimesCalled).To(Equal(1))
			Expect(courier.MapRouteCall.TimesCalled).To(Equal(4))
			Expect(courier.MapRouteCall.Received.AppName[0]).To(ContainSubstring(randomAppName + "-new-build-"))
			Expect(courier.MapRouteCall.Received.Domain[0]).To(Equal(randomDomain + "0"))
			Expect(courier.MapRouteCall.Received.Hostname[0]).To(Equal(randomHostName + "0"))
		})

		Context("when an error is returned from CustomRouteMapper", func() {
			It("returns the error", func() {
				manifest := fmt.Sprintf(`
---
applications:
- name: example
  custom-routes:`)

				pusher.DeploymentInfo.Manifest = manifest
				courier.ExistsCall.Returns.Bool = true
				courier.DomainsCall.Returns.Domains = []string{randomDomain + "0", randomDomain + "1", randomDomain + "2"}

				err := pusher.PostExecute()

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("Success", func() {
		It("renames the newly pushed app to the original name", func() {
			Expect(pusher.Success()).To(Succeed())

			Expect(courier.RenameCall.Received.AppName).To(Equal(randomAppName + TemporaryNameSuffix + randomUUID))
			Expect(courier.RenameCall.Received.AppNameVenerable).To(Equal(randomAppName))

			Eventually(logBuffer).Should(Say("renamed %s to %s", tempAppWithUUID, randomAppName))
		})

		Context("when rename fails", func() {
			It("returns an error", func() {
				courier.RenameCall.Returns.Output = []byte("rename output")
				courier.RenameCall.Returns.Error = errors.New("rename error")

				err := pusher.Success()
				Expect(err).To(MatchError(state.RenameError{randomAppName + TemporaryNameSuffix + randomUUID, []byte("rename output")}))

				Expect(courier.RenameCall.Received.AppName).To(Equal(randomAppName + TemporaryNameSuffix + randomUUID))
				Expect(courier.RenameCall.Received.AppNameVenerable).To(Equal(randomAppName))

				Eventually(logBuffer).Should(Say("could not rename %s to %s", tempAppWithUUID, randomAppName))
			})
		})

		Context("when the app exists", func() {
			BeforeEach(func() {
				courier.ExistsCall.Returns.Bool = true
			})

			It("checks the application exists", func() {
				Expect(pusher.Success()).To(Succeed())

				Expect(courier.ExistsCall.Received.AppName).To(Equal(randomAppName))
			})

			It("unmaps the load balanced route", func() {
				Expect(pusher.Success()).To(Succeed())

				Expect(courier.UnmapRouteCall.Received.AppName).To(Equal(randomAppName))
				Expect(courier.UnmapRouteCall.Received.Domain).To(Equal(randomDomain))
				Expect(courier.UnmapRouteCall.Received.Hostname).To(Equal(randomAppName))

				Eventually(logBuffer).Should(Say(fmt.Sprintf("unmapped route %s", randomAppName)))
			})

			It("deletes the original application ", func() {
				Expect(pusher.Success()).To(Succeed())

				Expect(courier.DeleteCall.Received.AppName).To(Equal(randomAppName))

				Eventually(logBuffer).Should(Say(fmt.Sprintf("deleted %s", randomAppName)))
			})

			Context("when domain is not provided", func() {
				It("does not call unmap route", func() {
					deploymentInfo.Domain = ""

					pusher = Pusher{
						Courier:        courier,
						DeploymentInfo: deploymentInfo,
						EventManager:   eventManager,
						Response:       response,
						Log:            interfaces.DeploymentLogger{Log: interfaces.DefaultLogger(logBuffer, logging.DEBUG, "pusher_test")},
					}

					pusher.Success()

					Expect(courier.UnmapRouteCall.Received.AppName).To(BeEmpty())
					Expect(courier.UnmapRouteCall.Received.Domain).To(BeEmpty())
					Expect(courier.UnmapRouteCall.Received.Hostname).To(BeEmpty())
				})
			})

			Context("when unmapping the route fails", func() {
				It("only logs an error", func() {
					courier.UnmapRouteCall.Returns.Output = []byte("unmap output")
					courier.UnmapRouteCall.Returns.Error = errors.New("Unmap Error")

					err := pusher.Success()
					Expect(err).To(MatchError(state.UnmapRouteError{randomAppName, []byte("unmap output")}))

					Eventually(logBuffer).Should(Say(fmt.Sprintf("could not unmap %s", randomAppName)))
				})
			})

			Context("when deleting the original app fails", func() {
				It("returns an error", func() {
					courier.ExistsCall.Returns.Bool = true
					courier.DeleteCall.Returns.Output = []byte("delete output")
					courier.DeleteCall.Returns.Error = errors.New("delete error")

					err := pusher.Success()
					Expect(err).To(MatchError(state.DeleteApplicationError{randomAppName, []byte("delete output")}))

					Eventually(logBuffer).Should(Say(fmt.Sprintf("could not delete %s", randomAppName)))
				})
			})
		})

		Context("when the application does not exist", func() {
			It("does not delete the non-existant original application", func() {
				courier.ExistsCall.Returns.Bool = false

				err := pusher.Success()
				Expect(err).ToNot(HaveOccurred())

				Expect(courier.DeleteCall.Received.AppName).To(BeEmpty())

				Eventually(logBuffer).ShouldNot(Say("delete"))
			})
		})

		It("should write the foundation URL to the log", func() {
			courier.ExistsCall.Returns.Bool = true
			pusher.Success()
			Eventually(logBuffer).Should(Say(randomFoundationURL + ": unmapping route"))
			Eventually(logBuffer).Should(Say(randomFoundationURL + ": unmapped route"))

		})
		Context("when UnmapRoute error returns an error", func() {
			It("should write the error message to the log with the foundation URL", func() {
				courier.ExistsCall.Returns.Bool = true
				courier.UnmapRouteCall.Returns.Error = errors.New("an error")
				pusher.Success()

				Eventually(logBuffer).Should(Say(randomFoundationURL + ": could not unmap"))
			})
		})
		Context("When deleteApplication is called", func() {
			It("should write the foundation URL to the log", func() {
				courier.ExistsCall.Returns.Bool = true
				pusher.Success()
				Eventually(logBuffer).Should(Say(randomFoundationURL + ": deleting"))
				Eventually(logBuffer).Should(Say(randomFoundationURL + ": deleted"))

			})
		})

		Context("When DeleteCall returns an error", func() {
			It("should write the error message with foundation URL to the log", func() {
				courier.ExistsCall.Returns.Bool = true
				courier.DeleteCall.Returns.Error = errors.New("an error")
				pusher.Success()
				Eventually(logBuffer).Should(Say(randomFoundationURL + ": could not delete"))
				Eventually(logBuffer).Should(Say(randomFoundationURL + ": deletion error"))
				Eventually(logBuffer).Should(Say(randomFoundationURL + ": deletion output"))

			})
		})

		Context("When renameNewBuildToOriginalAppName is called", func() {
			It("should write the foundation URL to the log", func() {
				courier.ExistsCall.Returns.Bool = true
				pusher.Success()
				Eventually(logBuffer).Should(Say(randomFoundationURL + ": renaming"))
				Eventually(logBuffer).Should(Say(randomFoundationURL + ": renamed"))

			})
		})

		Context("When RenameCall returns an error", func() {
			It("should write the error message with foundation URL to the log", func() {
				courier.ExistsCall.Returns.Bool = true
				courier.RenameCall.Returns.Error = errors.New("an error")
				pusher.Success()
				Eventually(logBuffer).Should(Say(randomFoundationURL + ": could not rename"))

			})
		})
	})

	Describe("Undo", func() {
		Context("when the app exists", func() {
			BeforeEach(func() {
				courier.ExistsCall.Returns.Bool = true
			})

			It("check that the app exists", func() {
				Expect(pusher.Undo()).To(Succeed())
				Expect(courier.ExistsCall.Received.AppName).To(Equal(randomAppName))
			})

			It("deletes the app that was pushed", func() {

				Expect(pusher.Undo()).To(Succeed())

				Expect(courier.DeleteCall.Received.AppName).To(Equal(randomAppName + TemporaryNameSuffix + randomUUID))

				Eventually(logBuffer).Should(Say(fmt.Sprintf("rolling back deploy of %s", randomAppName)))
				Eventually(logBuffer).Should(Say(fmt.Sprintf("deleted %s", randomAppName)))
			})

			Context("when deleting fails", func() {
				It("returns an error and writes a message to the info log", func() {
					courier.DeleteCall.Returns.Output = []byte("delete call output")
					courier.DeleteCall.Returns.Error = errors.New("delete error")

					err := pusher.Undo()
					Expect(err).To(MatchError(state.DeleteApplicationError{tempAppWithUUID, []byte("delete call output")}))

					Eventually(logBuffer).Should(Say(fmt.Sprintf("could not delete %s", tempAppWithUUID)))
				})
			})
		})

		Context("when the app does not exist", func() {
			It("renames the newly built app to the intended application name", func() {
				Expect(pusher.Undo()).To(Succeed())

				Expect(courier.RenameCall.Received.AppName).To(Equal(randomAppName + TemporaryNameSuffix + randomUUID))
				Expect(courier.RenameCall.Received.AppNameVenerable).To(Equal(randomAppName))

				Eventually(logBuffer).Should(Say("renamed %s to %s", tempAppWithUUID, randomAppName))
			})

			Context("when renaming fails", func() {
				It("returns an error and writes a message to the info log", func() {
					courier.RenameCall.Returns.Error = errors.New("rename error")
					courier.RenameCall.Returns.Output = []byte("rename error")

					err := pusher.Undo()
					Expect(err).To(MatchError(state.RenameError{tempAppWithUUID, []byte("rename error")}))

					Eventually(logBuffer).Should(Say(fmt.Sprintf("could not rename %s to %s", tempAppWithUUID, randomAppName)))
				})
			})
		})

		Context("when DisableRollback is true", func() {
			It("should write message with Foundation url", func() {
				pusher.Environment.DisableRollback = true
				pusher.Undo()

				Eventually(logBuffer).Should(Say(randomFoundationURL + ": Failed to deploy"))
			})
		})

		Context("when DisableRollback is false", func() {
			Context("when courier ExistCall returns true", func() {
				It("should write the message to the log", func() {
					pusher.Environment.DisableRollback = false
					courier.ExistsCall.Returns.Bool = true
					pusher.Undo()

					Eventually(logBuffer).Should(Say(randomFoundationURL + ": rolling back deploy of"))
				})
			})

			Context("when courier ExistCall returns false", func() {
				It("should write the message to the log", func() {
					pusher.Environment.DisableRollback = false
					courier.ExistsCall.Returns.Bool = false
					pusher.Undo()

					Eventually(logBuffer).Should(Say(randomFoundationURL + ": app"))
				})
			})
		})
	})

	Describe("Finally", func() {
		It("is successful", func() {
			courier.CleanUpCall.Returns.Error = nil

			Expect(pusher.Finally()).To(Succeed())
		})
	})

	Describe("Verify", func() {
		It("returns nil", func() {
			Expect(pusher.Verify()).To(BeNil())
		})
	})
})
