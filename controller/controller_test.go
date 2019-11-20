package controller_test

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"

	"io/ioutil"

	"os"

	"strings"

	"github.com/compozed/deployadactyl/config"
	. "github.com/compozed/deployadactyl/controller"
	I "github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/mocks"
	"github.com/compozed/deployadactyl/randomizer"
	"github.com/compozed/deployadactyl/request"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/op/go-logging"
)

var _ = Describe("Controller", func() {

	var (
		eventManager     *mocks.EventManager
		errorFinder      *mocks.ErrorFinder
		requestProcessor *mocks.RequestProcessor

		receivedBuffer  *bytes.Buffer
		receivedUuid    string
		receivedRequest interface{}

		controller *Controller
		logBuffer  *Buffer

		appName     string
		environment string
		org         string
		space       string
		byteBody    []byte
		server      *httptest.Server
		testUuid    string
	)

	BeforeEach(func() {
		logBuffer = NewBuffer()
		appName = strings.ToLower("appName-" + randomizer.StringRunes(10))
		environment = strings.ToLower("environment-" + randomizer.StringRunes(10))
		org = strings.ToLower("org-" + randomizer.StringRunes(10))
		space = strings.ToLower("non-prod")
		testUuid = "uuid1234"

		eventManager = &mocks.EventManager{}
		requestProcessor = &mocks.RequestProcessor{}

		requestFactory := func(uuid string, request interface{}, output *bytes.Buffer) I.RequestProcessor {
			receivedUuid = uuid
			receivedBuffer = output
			receivedRequest = request

			requestProcessor.Response = output
			return requestProcessor
		}

		errorFinder = &mocks.ErrorFinder{}
		controller = &Controller{
			Log: I.DefaultLogger(logBuffer, logging.DEBUG, "api_test"),
			RequestProcessorFactory: requestFactory,
			Config:                  config.Config{},
			ErrorFinder:             errorFinder,
		}
	})

	Describe("PostRequestHandler", func() {
		var (
			router        *gin.Engine
			resp          *httptest.ResponseRecorder
			jsonBuffer    *bytes.Buffer
			foundationURL string
		)
		BeforeEach(func() {
			router = gin.New()
			resp = httptest.NewRecorder()
			jsonBuffer = &bytes.Buffer{}

			router.POST("/v3/apps/:environment/:org/:space/:appName", controller.PostRequestHandler)

			server = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				byteBody, _ = ioutil.ReadAll(req.Body)
				req.Body.Close()
			}))

			silentDeployUrl := server.URL + "/v1/apps/" + os.Getenv("SILENT_DEPLOY_ENVIRONMENT")
			os.Setenv("SILENT_DEPLOY_URL", silentDeployUrl)
		})
		AfterEach(func() {
			server.Close()
		})

		Context("When the deserializer is unable to process the request", func() {

			It("provides an error response", func() {
				foundationURL = fmt.Sprintf("/v3/apps/%s/%s/%s/%s", environment, org, space, appName)

				jsonBuffer = bytes.NewBufferString(`{`)

				req, _ := http.NewRequest("POST", foundationURL, jsonBuffer)
				req.Header.Set("Content-Type", "application/json")

				router.ServeHTTP(resp, req)

				Eventually(resp.Code).Should(Equal(http.StatusBadRequest))
			})
		})

		It("creates the RequestProcessor", func() {
			foundationURL = fmt.Sprintf("/v3/apps/%s/%s/%s/%s", environment, org, space, appName)

			body := []byte(`{"artifact_url": "the url",
"environment_variables": {"foo": "bar"},
"health_check_endpoint": "the healthcheck",
"manifest": "the manifest",
"data": {"puppy": "dachshund"},
"uuid": "uuid1234"}`)

			jsonBuffer = bytes.NewBuffer(body)

			req, _ := http.NewRequest("POST", foundationURL, jsonBuffer)
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(resp, req)

			expectedEnv := make(map[string]string)
			expectedEnv["foo"] = "bar"

			expectedData := make(map[string]interface{})
			expectedData["puppy"] = "dachshund"

			Expect(*receivedBuffer).ToNot(BeNil())
			Expect(receivedRequest).To(Equal(request.PostDeploymentRequest{
				Deployment: I.Deployment{
					CFContext: I.CFContext{
						Environment:  environment,
						Organization: org,
						Space:        space,
						Application:  appName,
					},
					Body: &body,
					Type: "application/json",
				},
				Request: request.PostRequest{
					HealthCheckEndpoint:  "the healthcheck",
					ArtifactUrl:          "the url",
					EnvironmentVariables: expectedEnv,
					Manifest:             "the manifest",
					Data:                 expectedData,
					UUID:                 testUuid,
				},
			}))
			Expect(receivedUuid).To(Equal(testUuid))
		})

		It("calls Process on the RequestProcessor", func() {
			foundationURL = fmt.Sprintf("/v3/apps/%s/%s/%s/%s", environment, org, space, appName)

			jsonBuffer = bytes.NewBufferString("{}")

			req, _ := http.NewRequest("POST", foundationURL, jsonBuffer)
			req.Header.Set("Content-Type", "application/zip")

			router.ServeHTTP(resp, req)

			Expect(requestProcessor.ProcessCall.TimesCalled).To(Equal(1))
		})

		Context("when no uuid is supplied", func() {
			It("generates the uuid", func() {
				foundationURL = fmt.Sprintf("/v3/apps/%s/%s/%s/%s", environment, org, space, appName)

				body := []byte(`{"artifact_url": "the url",
"environment_variables": {"foo": "bar"},
"health_check_endpoint": "the healthcheck",
"manifest": "the manifest",
"data": {"puppy": "dachshund"}}`)

				jsonBuffer = bytes.NewBuffer(body)

				req, _ := http.NewRequest("POST", foundationURL, jsonBuffer)
				req.Header.Set("Content-Type", "application/json")

				router.ServeHTTP(resp, req)

				Expect(receivedUuid).ToNot(Equal(""))
			})
		})

		Context("when Process fails", func() {
			It("doesn't deploy and gives http.StatusInternalServerError", func() {
				foundationURL = fmt.Sprintf("/v3/apps/%s/%s/%s/%s", environment, org, space, appName)

				jsonBuffer = bytes.NewBufferString("{}")

				req, err := http.NewRequest("POST", foundationURL, jsonBuffer)
				req.Header.Set("Content-Type", "application/zip")

				Expect(err).ToNot(HaveOccurred())

				requestProcessor.ProcessCall.Returns.Response = I.DeployResponse{
					Error:      errors.New("bork"),
					StatusCode: http.StatusInternalServerError,
				}

				router.ServeHTTP(resp, req)

				Eventually(resp.Code).Should(Equal(http.StatusInternalServerError))
				Eventually(resp.Body).Should(ContainSubstring("bork"))
			})
		})

		Context("when parameters are added to the url", func() {
			It("does not return an error", func() {
				foundationURL = fmt.Sprintf("/v3/apps/%s/%s/%s/%s?broken=false", environment, org, space, appName)

				jsonBuffer = bytes.NewBufferString("{}")

				req, err := http.NewRequest("POST", foundationURL, jsonBuffer)
				req.Header.Set("Content-Type", "application/zip")

				Expect(err).ToNot(HaveOccurred())

				requestProcessor.ProcessCall.Returns.Response = I.DeployResponse{
					StatusCode: http.StatusOK,
				}

				router.ServeHTTP(resp, req)

				Eventually(resp.Code).Should(Equal(http.StatusOK))
				Expect(receivedRequest.(request.PostDeploymentRequest).CFContext).To(Equal(I.CFContext{
					Environment:  environment,
					Organization: org,
					Space:        space,
					Application:  appName,
				}))
			})
		})

		Context("when Process succeeds", func() {
			It("returns StatusOK", func() {
				foundationURL = fmt.Sprintf("/v3/apps/%s/%s/%s/%s", environment, org, space, appName)

				jsonBuffer = bytes.NewBufferString("{}")

				req, _ := http.NewRequest("POST", foundationURL, jsonBuffer)
				req.Header.Set("Content-Type", "application/json")

				requestProcessor.ProcessCall.Returns.Response = I.DeployResponse{
					StatusCode: http.StatusOK,
				}
				requestProcessor.ProcessCall.Writes = "deploy success"

				router.ServeHTTP(resp, req)

				Eventually(resp.Code).Should(Equal(http.StatusOK))
				Eventually(resp.Body).Should(ContainSubstring("deploy success"))
			})
		})
	})

	Describe("PutRequestHandler", func() {
		var (
			router     *gin.Engine
			resp       *httptest.ResponseRecorder
			jsonBuffer *bytes.Buffer
		)

		BeforeEach(func() {
			router = gin.New()
			resp = httptest.NewRecorder()
			jsonBuffer = &bytes.Buffer{}

			router.PUT("/v3/apps/:environment/:org/:space/:appName", controller.PutRequestHandler)

			server = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				byteBody, _ = ioutil.ReadAll(req.Body)
				req.Body.Close()
			}))
		})

		AfterEach(func() {
			server.Close()
		})

		Context("when stop succeeds", func() {
			It("returns http status.OK", func() {
				foundationURL := fmt.Sprintf("/v3/apps/%s/%s/%s/%s", environment, org, space, appName)
				jsonBuffer = bytes.NewBufferString(`{"state": "stopped"}`)

				req, err := http.NewRequest("PUT", foundationURL, jsonBuffer)
				req.Header.Set("Content-Type", "application/json")

				Expect(err).ToNot(HaveOccurred())

				router.ServeHTTP(resp, req)

				Eventually(resp.Code).Should(Equal(http.StatusOK))
			})
		})

		It("creates the RequestProcessor", func() {
			foundationURL := fmt.Sprintf("/v3/apps/%s/%s/%s/%s", environment, org, space, appName)

			body := []byte(`{"state": "started", "data": {"puppy": "dachshund"}, "uuid": "uuid1234"}`)

			jsonBuffer = bytes.NewBuffer(body)

			expectedData := make(map[string]interface{})
			expectedData["puppy"] = "dachshund"

			req, err := http.NewRequest("PUT", foundationURL, jsonBuffer)
			req.Header.Set("Content-Type", "application/json")

			Expect(err).ToNot(HaveOccurred())

			router.ServeHTTP(resp, req)

			Expect(*receivedBuffer).ToNot(BeNil())
			Expect(receivedRequest).To(Equal(request.PutDeploymentRequest{
				Deployment: I.Deployment{
					CFContext: I.CFContext{
						Environment:  environment,
						Organization: org,
						Space:        space,
						Application:  appName,
					},
					Body: &body,
					Type: "application/json",
				},
				Request: request.PutRequest{
					State: "started",
					Data:  expectedData,
					UUID:  testUuid,
				},
			}))
			Expect(receivedUuid).To(Equal(testUuid))
		})

		It("calls Process on the RequestProcessor", func() {
			foundationURL := fmt.Sprintf("/v3/apps/%s/%s/%s/%s", environment, org, space, appName)

			body := []byte(`{"state": "stopped"}`)

			jsonBuffer = bytes.NewBuffer(body)

			req, _ := http.NewRequest("PUT", foundationURL, jsonBuffer)
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(resp, req)

			Expect(requestProcessor.ProcessCall.TimesCalled).To(Equal(1))
		})

		It("logs request origination address", func() {
			foundationURL := fmt.Sprintf("/v3/apps/%s/%s/%s/%s", environment, org, space, appName)
			jsonBuffer = bytes.NewBufferString(`{"state": "stopped"}`)

			req, _ := http.NewRequest("PUT", foundationURL, jsonBuffer)
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(resp, req)

			Eventually(logBuffer).Should(Say("PUT Request originated from"))
		})

		Context("when no uuid is supplied", func() {
			It("generates the uuid", func() {
				foundationURL := fmt.Sprintf("/v3/apps/%s/%s/%s/%s", environment, org, space, appName)

				body := []byte(`{"artifact_url": "the url",
"environment_variables": {"foo": "bar"},
"health_check_endpoint": "the healthcheck",
"manifest": "the manifest",
"data": {"puppy": "dachshund"}}`)

				jsonBuffer = bytes.NewBuffer(body)

				req, _ := http.NewRequest("PUT", foundationURL, jsonBuffer)
				req.Header.Set("Content-Type", "application/json")

				router.ServeHTTP(resp, req)

				Expect(receivedUuid).ToNot(Equal(""))
			})
		})

		Context("when Process succeeds", func() {
			It("returns StatusOK", func() {
				foundationURL := fmt.Sprintf("/v3/apps/%s/%s/%s/%s", environment, org, space, appName)

				jsonBuffer = bytes.NewBufferString(`{"state": "started"}`)

				req, _ := http.NewRequest("PUT", foundationURL, jsonBuffer)
				req.Header.Set("Content-Type", "application/json")

				requestProcessor.ProcessCall.Returns.Response = I.DeployResponse{
					StatusCode: http.StatusOK,
				}
				requestProcessor.ProcessCall.Writes = "this is the process output"

				router.ServeHTTP(resp, req)

				Eventually(resp.Code).Should(Equal(http.StatusOK))
				Eventually(resp.Body).Should(ContainSubstring("this is the process output"))
			})
		})

		Context("when bad request body", func() {
			It("returns a Bad Request error", func() {
				foundationURL := fmt.Sprintf("/v3/apps/%s/%s/%s/%s", environment, org, space, appName)
				jsonBuffer := bytes.NewBufferString(`{`)

				req, err := http.NewRequest("PUT", foundationURL, jsonBuffer)
				req.Header.Set("Content-Type", "application/json")

				Expect(err).ToNot(HaveOccurred())

				router.ServeHTTP(resp, req)

				Expect(resp.Code).To(Equal(400))
				Expect(resp.Body.String()).To(Equal("Invalid request body."))
			})
		})
	})
})
