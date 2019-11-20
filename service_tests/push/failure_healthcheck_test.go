package push

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"

	"github.com/compozed/deployadactyl/creator"
	"github.com/compozed/deployadactyl/eventmanager/handlers/healthchecker"
	"github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/mocks"
	"github.com/compozed/deployadactyl/randomizer"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/spf13/afero"
)

var _ = Describe("", func() {
	const (
		CONFIGPATH      = "./test_config.yml"
		ENVIRONMENTNAME = "test"
		TESTCONFIG      = `---
environments:
- name: Test
  domain: example.com
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
		client              *mocks.Client
		provider            creator.CreatorModuleProvider

		couriers     []*mocks.Courier
		responseBody []byte
		response     *http.Response
		org          = strings.ToLower(randomizer.StringRunes(10))
		space        = strings.ToLower(os.Getenv("SILENT_DEPLOY_ENVIRONMENT"))
		appName      = strings.ToLower(randomizer.StringRunes(10))
	)

	BeforeEach(func() {
		os.Setenv("CF_USERNAME", randomizer.StringRunes(10))
		os.Setenv("CF_PASSWORD", randomizer.StringRunes(10))

		Expect(ioutil.WriteFile(CONFIGPATH, []byte(TESTCONFIG), 0644)).To(Succeed())

		prechecker = &mocks.Prechecker{}
		fetcher = &mocks.Fetcher{}
		eventManager = &mocks.EventManager{}
		couriers = make([]*mocks.Courier, 0)
		client = &mocks.Client{}

		client.GetCall.Returns.Response = http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       NewBuffer(),
		}

		healthChecker := healthchecker.HealthChecker{
			Client: client,
			OldURL: "api.cf",
			NewURL: "apps",
		}

		provider = creator.CreatorModuleProvider{
			NewPrechecker: func(eventManager interfaces.EventManager) interfaces.Prechecker {
				return prechecker
			},
			NewCourier: func(executor interfaces.Executor) interfaces.Courier {
				courier := &mocks.Courier{}
				couriers = append(couriers, courier)
				courier.ExistsCall.Returns.Bool = true

				return courier
			},
			NewFetcher: func(fs *afero.Afero, ex interfaces.Extractor, log interfaces.DeploymentLogger) interfaces.Fetcher {
				wd, _ := os.Getwd()

				dstf, _ := fs.TempDir("", "service-success-test-")
				dst, _ := fs.OpenFile(path.Join(dstf, "manifest.yml"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
				src, _ := fs.Open(wd + "/../../artifetcher/fixtures/deployadactyl-fixture-unzipped/manifest.yml")

				io.Copy(dst, src)

				fetcher.FetchCall.Returns.AppPath = dst.Name()
				return fetcher
			},
			NewEventManager: func(log interfaces.DeploymentLogger, bindings []interfaces.Binding) interfaces.EventManager {
				return eventManager
			},
			NewHealthChecker: func(oldURL, newURL, silentDeployURL, silentDeployEnvironment string, client interfaces.Client) healthchecker.HealthChecker {
				return healthChecker
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

	It("returns failed healthcheck response", func() {
		Expect(string(responseBody)).To(ContainSubstring("cannot deploy application: push failed"))
		Expect(string(responseBody)).To(ContainSubstring("health check failed"))
		Expect(string(responseBody)).To(ContainSubstring("status code: 500"))
		Expect(string(responseBody)).To(ContainSubstring("endpoint: /health"))
	})

	It("client receives correct input", func() {
		Expect(client.GetCall.Received.URL).To(ContainSubstring(".example.com/health"))
	})
})
