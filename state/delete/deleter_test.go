package delete_test

import (
	"errors"
	//"fmt"
	"math/rand"

	"github.com/compozed/deployadactyl/mocks"
	"github.com/compozed/deployadactyl/randomizer"
	. "github.com/compozed/deployadactyl/state/delete"
	. "github.com/compozed/deployadactyl/state/push"
	S "github.com/compozed/deployadactyl/structs"
	"github.com/op/go-logging"

	"fmt"

	"github.com/compozed/deployadactyl/state"

	"github.com/compozed/deployadactyl/interfaces"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
)

var _ = Describe("Deleter", func() {
	var (
		deleter      Deleter
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

		deleter = Deleter{
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

				Expect(deleter.Initially()).To(Succeed())

				Expect(courier.LoginCall.Received.FoundationURL).To(Equal(randomFoundationURL))
				Expect(courier.LoginCall.Received.Username).To(Equal(randomUsername))
				Expect(courier.LoginCall.Received.Password).To(Equal(randomPassword))
				Expect(courier.LoginCall.Received.Org).To(Equal(randomOrg))
				Expect(courier.LoginCall.Received.Space).To(Equal(randomSpace))
				Expect(courier.LoginCall.Received.SkipSSL).To(Equal(skipSSL))
			})

			It("writes the output of the courier to the response", func() {
				courier.LoginCall.Returns.Output = []byte("login succeeded")

				Expect(deleter.Initially()).To(Succeed())

				Eventually(response).Should(Say("login succeeded"))
			})
		})

		Context("when login fails", func() {
			It("returns an error", func() {
				courier.LoginCall.Returns.Output = []byte("login output")
				courier.LoginCall.Returns.Error = errors.New("login error")

				err := deleter.Initially()
				Expect(err).To(MatchError(state.LoginError{randomFoundationURL, []byte("login output")}))
			})

			It("writes the output of the courier to the response", func() {
				courier.LoginCall.Returns.Output = []byte("login output")
				courier.LoginCall.Returns.Error = errors.New("login error")

				err := deleter.Initially()
				Expect(err).To(HaveOccurred())

				Eventually(response).Should(Say("login output"))
			})

			It("logs an error", func() {
				courier.LoginCall.Returns.Error = errors.New("login error")

				err := deleter.Initially()
				Expect(err).To(HaveOccurred())

				Eventually(logBuffer).Should(Say(fmt.Sprintf("could not login to %s", randomFoundationURL)))
			})
		})
	})

	Describe("Execute", func() {
		Context("when the delete succeeds", func() {
			It("returns with success", func() {
				courier.ExistsCall.Returns.Bool = true
				courier.DeleteCall.Returns.Output = []byte("delete succeeded")

				Expect(deleter.Execute()).To(Succeed())

				Expect(courier.DeleteCall.Received.AppName).To(Equal(randomAppName))

				Eventually(response).Should(Say("delete succeeded"))

				Eventually(logBuffer).Should(Say(fmt.Sprintf("%s: deleting app %s", randomFoundationURL, randomAppName)))
				Eventually(logBuffer).Should(Say(fmt.Sprintf("%s: successfully deleted app %s", randomFoundationURL, randomAppName)))
			})
		})

		Context("when the delete fails", func() {
			It("returns an error", func() {
				courier.ExistsCall.Returns.Bool = true
				courier.DeleteCall.Returns.Output = []byte("this is some output")
				courier.DeleteCall.Returns.Error = errors.New("")

				err := deleter.Execute()

				Expect(err).To(MatchError(state.DeleteError{ApplicationName: randomAppName, Out: []byte("this is some output")}))
			})
		})

		Context("when the app does not exist", func() {
			It("returns an error", func() {
				courier.ExistsCall.Returns.Bool = false

				err := deleter.Execute()

				Expect(err).To(MatchError(state.ExistsError{ApplicationName: randomAppName}))
			})
		})
	})

	Describe("Undo", func() {
		Context("when successful", func() {
			It("returns with success", func() {
				Expect(deleter.Undo()).To(Succeed())

				Eventually(response).Should(Say("delete feature is unable to rollback"))
			})

			It("should write the message with the foundation URL to the log", func() {
				Expect(deleter.Undo()).To(Succeed())

				Eventually(logBuffer).Should(Say(randomFoundationURL + ": delete feature is unable to rollback"))
			})
		})
	})

	Describe("Verify", func() {
		It("returns nil", func() {
			Expect(deleter.Verify()).To(BeNil())
		})
	})

	Describe("Success", func() {
		It("returns nil", func() {
			Expect(deleter.Success()).To(BeNil())
		})
	})

	Describe("Finally", func() {
		It("returns nil", func() {
			Expect(deleter.Finally()).To(BeNil())
		})
	})
})
