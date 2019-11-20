package push

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"

	"encoding/base64"
	"errors"
	"io"
	"path"
	"reflect"

	"github.com/compozed/deployadactyl/creator"
	"github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/mocks"
	"github.com/compozed/deployadactyl/randomizer"
	"github.com/compozed/deployadactyl/state/push"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

var _ = Describe("Service", func() {

	const (
		CONFIGPATH      = "./test_config.yml"
		ENVIRONMENTNAME = "test"
		TESTCONFIG      = `---
environments:
- name: Test
  domain: example.com
  enable_rollback: false
  foundations:
  - api1.example.com
  - api2.example.com
  - api3.example.com
  - api4.example.com
`
		manifest = `---
applications:
- name: healthy-timmy
  memory: 10MB
  disk_quota: 5MB
  instances: 1`
	)

	var (
		deployadactylServer *httptest.Server
		prechecker          *mocks.Prechecker
		fetcher             *mocks.Fetcher
		eventManager        *mocks.EventManager
		provider            creator.CreatorModuleProvider

		couriers     []*mocks.Courier
		responseBody []byte
		response     *http.Response
		org          = randomizer.StringRunes(10)
		space        = os.Getenv("SILENT_DEPLOY_ENVIRONMENT")
		appName      = randomizer.StringRunes(10)
	)

	BeforeEach(func() {
		os.Setenv("CF_USERNAME", randomizer.StringRunes(10))
		os.Setenv("CF_PASSWORD", randomizer.StringRunes(10))

		Expect(ioutil.WriteFile(CONFIGPATH, []byte(TESTCONFIG), 0644)).To(Succeed())

		prechecker = &mocks.Prechecker{}
		fetcher = &mocks.Fetcher{}
		eventManager = &mocks.EventManager{}
		couriers = make([]*mocks.Courier, 0)

		provider = creator.CreatorModuleProvider{
			NewPrechecker: func(eventManager interfaces.EventManager) interfaces.Prechecker {
				prechecker.AssertAllFoundationsUpCall.Returns.Error = errors.New("a precheck error")
				return prechecker
			},
			NewCourier: func(executor interfaces.Executor) interfaces.Courier {
				courier := &mocks.Courier{}
				couriers = append(couriers, courier)
				courier.ExistsCall.Returns.Bool = true
				if len(couriers) == 2 {
					courier.PushCall.Returns.Error = errors.New("failed to push")
				}

				return courier
			},
			NewFetcher: func(fs *afero.Afero, ex interfaces.Extractor, log interfaces.DeploymentLogger) interfaces.Fetcher {
				wd, _ := os.Getwd()

				dstf, _ := fs.TempDir("", "service-failure-test-")
				dst, _ := fs.OpenFile(path.Join(dstf, "manifest.yml"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
				src, _ := fs.Open(wd + "/../../artifetcher/fixtures/deployadactyl-fixture-unzipped/manifest.yml")

				io.Copy(dst, src)

				fetcher.FetchCall.Returns.AppPath = dst.Name()
				return fetcher
			},
			NewEventManager: func(log interfaces.DeploymentLogger, bindings []interfaces.Binding) interfaces.EventManager {
				return eventManager
			},
		}

		creator, err := creator.Custom("DEBUG", CONFIGPATH, provider)

		Expect(err).ToNot(HaveOccurred())

		controller := creator.CreateController()
		deployadactylHandler := creator.CreateControllerHandler(controller)

		deployadactylServer = httptest.NewServer(deployadactylHandler)

		j, err := json.Marshal(gin.H{
			"artifact_url":          "the artifact url",
			"health_check_endpoint": "/health",
			"manifest":              base64.StdEncoding.EncodeToString([]byte(manifest)),
		})
		Expect(err).ToNot(HaveOccurred())
		jsonBuffer := bytes.NewBuffer(j)

		requestURL := fmt.Sprintf("%s/v3/apps/%s/%s/%s/%s", deployadactylServer.URL, ENVIRONMENTNAME, org, space, appName)
		req, err := http.NewRequest("POST", requestURL, jsonBuffer)
		Expect(err).ToNot(HaveOccurred())

		req.Header.Add("Content-Type", "application/json")

		client := &http.Client{}

		response, err = client.Do(req)
		Expect(err).ToNot(HaveOccurred())

		responseBody, err = ioutil.ReadAll(response.Body)
		response.Body.Close()
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		deployadactylServer.Close()
	})

	It("returns correct status code", func() {
		Expect(response.StatusCode).To(Equal(http.StatusInternalServerError), string(responseBody))
	})
	It("calls prechecker with all foundation urls", func() {
		fs := prechecker.AssertAllFoundationsUpCall.Received.Environment.Foundations
		Expect(fs).To(Equal([]string{"api1.example.com", "api2.example.com", "api3.example.com", "api4.example.com"}))
	})
	It("calls Emit the correct number of times", func() {
		Expect(len(eventManager.EmitCall.Received.Events)).To(Equal(3))
	})
	It("emits a deploy.start event", func() {
		Expect(eventManager.EmitCall.Received.Events[0].Type).To(Equal("deploy.start"))
	})
	It("emits a deploy.failure event", func() {
		Expect(eventManager.EmitCall.Received.Events[1].Type).To(Equal("deploy.failure"))
	})
	It("emits a deploy.finish event", func() {
		Expect(eventManager.EmitCall.Received.Events[2].Type).To(Equal("deploy.finish"))
	})
	It("calls EmitEvent the correct number of times", func() {
		Expect(len(eventManager.EmitEventCall.Received.Events)).To(Equal(3))
	})
	It("emits a DeployStartedEvent", func() {
		Expect(reflect.TypeOf(eventManager.EmitEventCall.Received.Events[0])).To(Equal(reflect.TypeOf(push.DeployStartedEvent{})))
	})
	It("emits a DeployFailureEvent", func() {
		Expect(reflect.TypeOf(eventManager.EmitEventCall.Received.Events[1])).To(Equal(reflect.TypeOf(push.DeployFailureEvent{})))
	})
	It("emits a DeployFinishedEvent", func() {
		Expect(reflect.TypeOf(eventManager.EmitEventCall.Received.Events[2])).To(Equal(reflect.TypeOf(push.DeployFinishedEvent{})))
	})
})
