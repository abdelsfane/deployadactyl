package healthchecker_test

import (
	"errors"
	"fmt"
	"net/http"

	. "github.com/compozed/deployadactyl/eventmanager/handlers/healthchecker"
	"github.com/compozed/deployadactyl/mocks"
	"github.com/compozed/deployadactyl/randomizer"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

	I "github.com/compozed/deployadactyl/interfaces"
	"github.com/op/go-logging"
)

var _ = Describe("Healthchecker", func() {

	var (
		randomAppName       string
		randomFoundationURL string
		randomDomain        string
		randomEndpoint      string
		randomUsername      string
		randomPassword      string
		randomOrg           string
		randomSpace         string
		randomHostname      string
		randomEnvironment   string
		randomUUID          string
		tempAppWithUUID     string

		healthchecker      HealthChecker
		healthCheckRequest HealthCheckRequest
		client             *mocks.Client
		courier            *mocks.Courier
		logBuffer          *Buffer
	)

	BeforeEach(func() {
		randomAppName = "randomAppName-" + randomizer.StringRunes(10)

		s := "random-" + randomizer.StringRunes(10)

		randomFoundationURL = fmt.Sprintf("https://api.cf.%s.com", s)
		randomDomain = fmt.Sprintf("apps.%s.com", s)
		randomHostname = randomAppName

		randomEndpoint = "/" + randomizer.StringRunes(10)

		randomUsername = "randomUsername" + randomizer.StringRunes(10)
		randomPassword = "randomPassword" + randomizer.StringRunes(10)
		randomOrg = "randomOrg" + randomizer.StringRunes(10)
		randomSpace = "randomSpace" + randomizer.StringRunes(10)
		randomEnvironment = "randomEnvironment" + randomizer.StringRunes(10)
		randomUUID = "randomUUID" + randomizer.StringRunes(10)
		tempAppWithUUID = randomAppName + "-new-build-" + randomUUID

		courier = &mocks.Courier{}
		client = &mocks.Client{}

		logBuffer = NewBuffer()

		healthchecker = HealthChecker{
			OldURL:                  "api.cf",
			NewURL:                  "apps",
			SilentDeployURL:         "silentapps",
			SilentDeployEnvironment: "silentenvironment",
			Client:                  client,
		}

		healthCheckRequest = HealthCheckRequest{
			HealthCheckEndpoint: randomEndpoint,
			Courier:             courier,
			Logger:              I.DeploymentLogger{Log: I.DefaultLogger(logBuffer, logging.DEBUG, "healthchecker_test")},
			Environment:         randomEnvironment,
			FoundationUrl:       randomFoundationURL,
			TempAppWithUUID:     tempAppWithUUID,
		}

		client.GetCall.Returns.Response = http.Response{
			StatusCode: http.StatusOK,
			Body:       NewBuffer(),
		}
	})

	Describe("HealthChecker", func() {
		Context("the new build application is healthy", func() {
			Context("the endpoint provided is valid", func() {
				It("does not return an error", func() {
					client.GetCall.Returns.Response = http.Response{StatusCode: http.StatusOK}

					err := healthchecker.HealthChecker(healthCheckRequest)

					Expect(err).ToNot(HaveOccurred())
				})

				It("maps a new temporary route", func() {
					healthchecker.HealthChecker(healthCheckRequest)

					Expect(courier.MapRouteCall.Received.AppName[0]).To(Equal(tempAppWithUUID))
					Expect(courier.MapRouteCall.Received.Domain[0]).To(Equal(randomDomain))
					Expect(courier.MapRouteCall.Received.Hostname[0]).To(Equal(tempAppWithUUID))
				})

				It("formats the foundation url", func() {
					healthchecker.HealthChecker(healthCheckRequest)

					Expect(client.GetCall.Received.URL).To(Equal(fmt.Sprintf("https://%s.%s%s", tempAppWithUUID, randomDomain, randomEndpoint)))
				})

				It("unmaps the temporary route", func() {
					healthchecker.HealthChecker(healthCheckRequest)

					Expect(courier.UnmapRouteCall.Received.AppName).To(Equal(tempAppWithUUID))
					Expect(courier.UnmapRouteCall.Received.Domain).To(Equal(randomDomain))
					Expect(courier.UnmapRouteCall.Received.Hostname).To(Equal(tempAppWithUUID))
				})

				It("deletes the temporary routes", func() {
					healthchecker.HealthChecker(healthCheckRequest)

					Expect(courier.DeleteRouteCall.Received.Domain).To(Equal(randomDomain))
					Expect(courier.DeleteRouteCall.Received.Hostname).To(Equal(tempAppWithUUID))
				})

				It("unmaps the temporary route before deleting it", func() {
					healthchecker.HealthChecker(healthCheckRequest)

					Expect(courier.UnmapRouteCall.OrderCalled < courier.DeleteRouteCall.OrderCalled).To(Equal(true))
				})

				It("prints success logs to the console", func() {
					client.GetCall.Returns.Response = http.Response{StatusCode: http.StatusOK}

					healthchecker.HealthChecker(healthCheckRequest)

					Eventually(logBuffer).Should(Say("starting health check"))
					Eventually(logBuffer).Should(Say("mapping temporary route %s.%s", tempAppWithUUID, randomDomain))
					Eventually(logBuffer).Should(Say("mapped temporary route %s.%s", tempAppWithUUID, randomDomain))
					Eventually(logBuffer).Should(Say("checking route https://%s.%s%s", tempAppWithUUID, randomDomain, randomEndpoint))
					Eventually(logBuffer).Should(Say("health check successful for https://%s.%s%s", tempAppWithUUID, randomDomain, randomEndpoint))
					Eventually(logBuffer).Should(Say("unmapping temporary route %s.%s", tempAppWithUUID, randomDomain))
					Eventually(logBuffer).Should(Say("unmapped temporary route %s.%s", tempAppWithUUID, randomDomain))
					Eventually(logBuffer).Should(Say("finished health check"))
				})

				It("maps route for silent deploy environment", func() {
					healthchecker = HealthChecker{
						OldURL:                  "api.cf",
						NewURL:                  "apps",
						SilentDeployURL:         "silentapps",
						SilentDeployEnvironment: randomEnvironment,
						Client:                  client,
					}

					healthchecker.HealthChecker(healthCheckRequest)

					Expect(courier.MapRouteCall.Received.Domain[0]).To(ContainSubstring("silentapps"))
					Eventually(logBuffer).Should(Say("finished health check"))
				})
			})

			Context("the endpoint provided is not valid", func() {
				BeforeEach(func() {
					client.GetCall.Returns.Response = http.Response{
						StatusCode: http.StatusNotFound,
						Body:       NewBuffer(),
					}
				})

				It("returns an error", func() {
					body := []byte("Could not find page")

					buf := NewBuffer()
					buf.Write(body)

					client.GetCall.Returns.Response = http.Response{
						StatusCode: http.StatusNotFound,
						Body:       buf,
					}

					err := healthchecker.HealthChecker(healthCheckRequest)

					Expect(err).To(MatchError(HealthCheckError{http.StatusNotFound, randomEndpoint, body}))
				})

				It("prints the endpoint error to the console", func() {
					healthchecker.HealthChecker(healthCheckRequest)

					Eventually(logBuffer).Should(Say("starting health check"))
					Eventually(logBuffer).Should(Say("mapping temporary route %s.%s", tempAppWithUUID, randomDomain))
					Eventually(logBuffer).Should(Say("mapped temporary route %s.%s", tempAppWithUUID, randomDomain))
					Eventually(logBuffer).Should(Say("checking route https://%s.%s%s", tempAppWithUUID, randomDomain, randomEndpoint))
					Eventually(logBuffer).Should(Say("health check failed"))
					Eventually(logBuffer).Should(Say(randomEndpoint))
				})
			})
		})

		Context("the new build application is not healthy", func() {
			It("returns an error", func() {
				client.GetCall.Returns.Response = http.Response{
					StatusCode: http.StatusNotFound,
					Body:       NewBuffer(),
				}

				err := healthchecker.HealthChecker(healthCheckRequest)
				Expect(err).To(MatchError(HealthCheckError{http.StatusNotFound, randomEndpoint, []byte{}}))
			})
		})

		Context("when mapping the temporary route fails", func() {
			It("returns an error", func() {
				courier.MapRouteCall.Returns.Output = append(courier.MapRouteCall.Returns.Output, []byte("map route output"))
				courier.MapRouteCall.Returns.Error = append(courier.MapRouteCall.Returns.Error, errors.New("map route error"))

				healthchecker.HealthChecker(healthCheckRequest)

				Eventually(logBuffer).Should(Say("mapping temporary route"))
				Eventually(logBuffer).Should(Say("failed to map temporary route"))
				Eventually(logBuffer).Should(Say("map route output"))
			})
		})

		Context("when the client fails to send the GET", func() {
			It("returns an error", func() {
				client.GetCall.Returns.Error = errors.New("client GET error")

				err := healthchecker.HealthChecker(healthCheckRequest)

				Expect(err).To(MatchError(ClientError{errors.New("client GET error")}))
			})

			It("prints the error to the logs", func() {
				client.GetCall.Returns.Error = errors.New("client GET error")

				healthchecker.HealthChecker(healthCheckRequest)

				Eventually(logBuffer).Should(Say("checking route"))
				Eventually(logBuffer).Should(Say("client GET error"))
			})
		})

		Context("when a health check endpoint is not provided", func() {
			It("returns nil", func() {
				err := healthchecker.HealthChecker(healthCheckRequest)

				Expect(err).To(BeNil())
			})
		})

		Context("when unmapping the temporary route fails", func() {
			It("prints output to the logs", func() {
				courier.UnmapRouteCall.Returns.Output = []byte("unmap route output")
				courier.UnmapRouteCall.Returns.Error = errors.New("unmap route error")

				healthchecker.HealthChecker(healthCheckRequest)

				Eventually(logBuffer).Should(Say("unmapping temporary route"))
				Eventually(logBuffer).Should(Say("failed to unmap temporary route"))
				Eventually(logBuffer).Should(Say("unmap route output"))
				Eventually(logBuffer).Should(Say("finished"))
			})
		})
	})

	Describe("format of endpoint parameter", func() {
		Context("when the endpoint does not include a '/'", func() {
			It("adds the leading '/'", func() {
				endpoint := "health"

				healthchecker.Check(randomFoundationURL, endpoint, I.DeploymentLogger{Log: I.DefaultLogger(logBuffer, logging.DEBUG, "healthchecker_test")}, randomUUID, randomFoundationURL)

				Expect(client.GetCall.Received.URL).To(Equal(fmt.Sprintf("%s/%s", randomFoundationURL, endpoint)))
			})
		})

		Context("when the endpoint does include a '/'", func() {
			It("does not add the leading '/'", func() {
				endpoint := "/health"

				healthchecker.Check(randomFoundationURL, endpoint, I.DeploymentLogger{Log: I.DefaultLogger(logBuffer, logging.DEBUG, "healthchecker_test")}, randomUUID, randomFoundationURL)

				Expect(client.GetCall.Received.URL).To(Equal(fmt.Sprintf("%s%s", randomFoundationURL, endpoint)))
			})
		})
	})
})
