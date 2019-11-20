package start_test

import (
	"errors"
	//"fmt"
	"math/rand"

	"github.com/compozed/deployadactyl/mocks"
	"github.com/compozed/deployadactyl/randomizer"
	. "github.com/compozed/deployadactyl/state/push"
	. "github.com/compozed/deployadactyl/state/start"
	S "github.com/compozed/deployadactyl/structs"
	"github.com/op/go-logging"

	"fmt"

	"github.com/compozed/deployadactyl/state"

	"github.com/compozed/deployadactyl/interfaces"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
)

var _ = Describe("Starter", func() {
	var (
		starter      Starter
		courier      *mocks.Courier
		eventManager *mocks.EventManager

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
		tempAppWithUUID     string
		skipSSL             bool
		deploymentInfo      S.DeploymentInfo
		cfContext           interfaces.CFContext
		auth                interfaces.Authorization
		response            *Buffer
		logBuffer           *Buffer
	)

	BeforeEach(func() {
		courier = &mocks.Courier{}
		eventManager = &mocks.EventManager{}

		randomFoundationURL = "randomFoundationURL-" + randomizer.StringRunes(10)
		randomUsername = "randomUsername-" + randomizer.StringRunes(10)
		randomPassword = "randomPassword-" + randomizer.StringRunes(10)
		randomOrg = "randomOrg-" + randomizer.StringRunes(10)
		randomSpace = "randomSpace-" + randomizer.StringRunes(10)
		randomDomain = "randomDomain-" + randomizer.StringRunes(10)
		randomAppPath = "randomAppPath-" + randomizer.StringRunes(10)
		randomAppName = "randomAppName-" + randomizer.StringRunes(10)
		randomEndpoint = "randomEndpoint-" + randomizer.StringRunes(10)
		randomUUID = randomizer.StringRunes(10)
		randomInstances = uint16(rand.Uint32())

		tempAppWithUUID = randomAppName + TemporaryNameSuffix + randomUUID

		response = NewBuffer()
		logBuffer = NewBuffer()

		eventManager.EmitCall.Returns.Error = append(eventManager.EmitCall.Returns.Error, nil)

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
		}

		cfContext = interfaces.CFContext{
			Organization: randomOrg,
			Space:        randomSpace,
			Application:  randomAppName,
		}

		auth = interfaces.Authorization{
			Username: randomUsername,
			Password: randomPassword,
		}

		starter = Starter{
			Courier:       courier,
			CFContext:     cfContext,
			Authorization: auth,
			EventManager:  eventManager,
			Response:      response,
			Log:           interfaces.DeploymentLogger{Log: interfaces.DefaultLogger(logBuffer, logging.DEBUG, "pusher_test")},
			FoundationURL: randomFoundationURL,
			AppName:       randomAppName,
		}
	})

	Describe("Initially", func() {
		Context("when login succeeds", func() {
			It("gives the correct info to the courier", func() {

				Expect(starter.Initially()).To(Succeed())

				Expect(courier.LoginCall.Received.FoundationURL).To(Equal(randomFoundationURL))
				Expect(courier.LoginCall.Received.Username).To(Equal(randomUsername))
				Expect(courier.LoginCall.Received.Password).To(Equal(randomPassword))
				Expect(courier.LoginCall.Received.Org).To(Equal(randomOrg))
				Expect(courier.LoginCall.Received.Space).To(Equal(randomSpace))
				Expect(courier.LoginCall.Received.SkipSSL).To(Equal(skipSSL))
			})

			It("writes the output of the courier to the response", func() {
				courier.LoginCall.Returns.Output = []byte("login succeeded")

				Expect(starter.Initially()).To(Succeed())

				Eventually(response).Should(Say("login succeeded"))
			})
		})

		Context("when login fails", func() {
			It("returns an error", func() {
				courier.LoginCall.Returns.Output = []byte("login output")
				courier.LoginCall.Returns.Error = errors.New("login error")

				err := starter.Initially()
				Expect(err).To(MatchError(state.LoginError{randomFoundationURL, []byte("login output")}))
			})

			It("writes the output of the courier to the response", func() {
				courier.LoginCall.Returns.Output = []byte("login output")
				courier.LoginCall.Returns.Error = errors.New("login error")

				err := starter.Initially()
				Expect(err).To(HaveOccurred())

				Eventually(response).Should(Say("login output"))
			})

			It("logs an error", func() {
				courier.LoginCall.Returns.Error = errors.New("login error")

				err := starter.Initially()
				Expect(err).To(HaveOccurred())

				Eventually(logBuffer).Should(Say(fmt.Sprintf("could not login to %s", randomFoundationURL)))
			})
		})
	})

	Describe("Execute", func() {
		Context("when the start succeeds", func() {
			It("returns with success", func() {
				courier.ExistsCall.Returns.Bool = true
				courier.StartCall.Returns.Output = []byte("start succeeded")

				Expect(starter.Execute()).To(Succeed())

				Expect(courier.StartCall.Received.AppName).To(Equal(randomAppName))

				Eventually(response).Should(Say("start succeeded"))

				Eventually(logBuffer).Should(Say(fmt.Sprintf("%s: starting app %s", randomFoundationURL, randomAppName)))
				Eventually(logBuffer).Should(Say(fmt.Sprintf("%s: successfully started app %s", randomFoundationURL, randomAppName)))
			})
		})

		Context("when the start fails", func() {
			It("returns an error", func() {
				courier.ExistsCall.Returns.Bool = true
				courier.StartCall.Returns.Output = []byte("this is some output")
				courier.StartCall.Returns.Error = errors.New("")

				err := starter.Execute()

				Expect(err).To(MatchError(state.StartError{ApplicationName: randomAppName, Out: []byte("this is some output")}))
			})
		})

		Context("when the app does not exist", func() {
			It("returns an error", func() {
				courier.ExistsCall.Returns.Bool = false

				err := starter.Execute()

				Expect(err).To(MatchError(state.ExistsError{ApplicationName: randomAppName}))
			})
		})
	})

	Describe("Undo", func() {
		Context("when the app does not exist", func() {
			It("returns an error", func() {
				courier.ExistsCall.Returns.Bool = false
				err := starter.Undo()

				Expect(err).To(MatchError(state.ExistsError{ApplicationName: randomAppName}))
			})
		})

		Context("when the stop fails", func() {
			It("returns an error", func() {
				courier.ExistsCall.Returns.Bool = true
				courier.StopCall.Returns.Output = []byte("this is some output")
				courier.StopCall.Returns.Error = errors.New("app could not be started")

				err := starter.Undo()

				Expect(err).To(MatchError(state.StopError{ApplicationName: randomAppName, Out: []byte("this is some output")}))
			})
		})

		Context("when successful", func() {
			It("returns with success", func() {
				courier.ExistsCall.Returns.Bool = true
				courier.StopCall.Returns.Output = []byte("stop succeeded")

				Expect(starter.Undo()).To(Succeed())
				Expect(courier.StopCall.Received.AppName).To(Equal(randomAppName))

				Eventually(response).Should(Say("stop succeeded"))
				Eventually(logBuffer).Should(Say(fmt.Sprintf("%s: stopping app %s", randomFoundationURL, randomAppName)))
				Eventually(logBuffer).Should(Say(fmt.Sprintf("%s: successfully restopped app %s", randomFoundationURL, randomAppName)))
			})
		})
	})

	Describe("Verify", func() {
		It("returns nil", func() {
			Expect(starter.Verify()).To(BeNil())
		})
	})

	Describe("Success", func() {
		It("returns nil", func() {
			Expect(starter.Success()).To(BeNil())
		})
	})

	Describe("Finally", func() {
		It("returns nil", func() {
			Expect(starter.Finally()).To(BeNil())
		})
	})
})
