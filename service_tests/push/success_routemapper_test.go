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
	"github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/mocks"
	"github.com/compozed/deployadactyl/randomizer"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
  instances: 1
  custom-routes:
  - route: custom-route-1.example.com
  - route: custom-route-2.example.com`
	)

	var (
		deployadactylServer *httptest.Server
		prechecker          *mocks.Prechecker
		fetcher             *mocks.Fetcher
		eventManager        *mocks.EventManager
		courier             *mocks.Courier
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
		courier = &mocks.Courier{}
		couriers = make([]*mocks.Courier, 0)

		provider = creator.CreatorModuleProvider{
			NewPrechecker: func(eventManager interfaces.EventManager) interfaces.Prechecker {
				return prechecker
			},
			NewCourier: func(executor interfaces.Executor) interfaces.Courier {
				couriers = append(couriers, courier)
				courier.ExistsCall.Returns.Bool = true
				courier.DomainsCall.Returns.Domains = []string{"example.com", "notValid.com"}

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
		}

		creator, err := creator.Custom("DEBUG", CONFIGPATH, provider)

		Expect(err).ToNot(HaveOccurred())

		controller := creator.CreateController()
		deployadactylHandler := creator.CreateControllerHandler(controller)

		deployadactylServer = httptest.NewServer(deployadactylHandler)

		j, err := json.Marshal(gin.H{
			"artifact_url": "the artifact url",
			"manifest":     base64.StdEncoding.EncodeToString([]byte(manifest)),
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
		Expect(response.StatusCode).To(Equal(http.StatusOK), string(responseBody))
	})

	It("MapRoute is called with correct input", func() {
		Expect(courier.MapRouteCall.Received.Domain[0]).To(Equal("example.com"))
		Expect(courier.MapRouteCall.Received.Hostname[0]).To(ContainSubstring("custom-route"))
		Expect(courier.MapRouteCall.Received.AppName[0]).To(ContainSubstring(appName))

		Expect(courier.MapRouteCall.Received.Domain[1]).To(Equal("example.com"))
		Expect(courier.MapRouteCall.Received.Hostname[1]).To(ContainSubstring("custom-route"))
		Expect(courier.MapRouteCall.Received.AppName[1]).To(ContainSubstring(appName))
	})
})
