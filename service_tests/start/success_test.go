package start

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"

	"reflect"
	"strings"

	"github.com/compozed/deployadactyl/creator"
	"github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/mocks"
	"github.com/compozed/deployadactyl/randomizer"
	"github.com/compozed/deployadactyl/state/start"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service", func() {

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
	)

	var (
		deployadactylServer *httptest.Server
		prechecker          *mocks.Prechecker
		eventManager        *mocks.EventManager
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
		eventManager = &mocks.EventManager{}
		couriers = make([]*mocks.Courier, 0)

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
			"state": "started",
		})
		Expect(err).ToNot(HaveOccurred())
		jsonBuffer := bytes.NewBuffer(j)

		requestURL := fmt.Sprintf("%s/v3/apps/%s/%s/%s/%s", deployadactylServer.URL, ENVIRONMENTNAME, org, space, appName)
		req, err := http.NewRequest("PUT", requestURL, jsonBuffer)
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
	It("calls prechecker with all foundation urls", func() {
		fs := prechecker.AssertAllFoundationsUpCall.Received.Environment.Foundations
		Expect(fs).To(Equal([]string{"api1.example.com", "api2.example.com", "api3.example.com", "api4.example.com"}))
	})
	It("creates correct number of courier objects", func() {
		Expect(len(couriers)).To(Equal(4))
	})
	It("calls courier login with correct info", func() {
		for _, c := range couriers {
			Expect(c.LoginCall.Received.Username).To(Equal(os.Getenv("CF_USERNAME")))
			Expect(c.LoginCall.Received.Password).To(Equal(os.Getenv("CF_PASSWORD")))
			Expect(c.LoginCall.Received.Org).To(Equal(org))
			Expect(c.LoginCall.Received.Space).To(Equal(space))
			Expect(c.LoginCall.Received.SkipSSL).To(Equal(false))
		}
	})
	It("calls courier login with correct foundation url", func() {
		furls := []string{"api1.example.com", "api2.example.com", "api3.example.com", "api4.example.com"}
		for i, c := range couriers {
			Expect(c.LoginCall.Received.FoundationURL).To(Equal(furls[i]))
		}
	})
	It("checks for prior existence of the app", func() {
		for _, c := range couriers {
			Expect(c.ExistsCall.Received.AppName).To(Equal(appName))
		}
	})
	It("calls courier stop with correct info", func() {
		for _, c := range couriers {
			Expect(c.StartCall.Received.AppName).To(ContainSubstring(appName))
		}
	})
	It("calls EmitEvent the correct number of times", func() {
		Expect(len(eventManager.EmitEventCall.Received.Events)).To(Equal(3))
	})
	It("emits a StartStartedEvent", func() {
		Expect(reflect.TypeOf(eventManager.EmitEventCall.Received.Events[0])).To(Equal(reflect.TypeOf(start.StartStartedEvent{})))
	})
	It("emits a StartSuccessEvent", func() {
		Expect(reflect.TypeOf(eventManager.EmitEventCall.Received.Events[1])).To(Equal(reflect.TypeOf(start.StartSuccessEvent{})))
	})
	It("emits a StartFinishedEvent", func() {
		Expect(reflect.TypeOf(eventManager.EmitEventCall.Received.Events[2])).To(Equal(reflect.TypeOf(start.StartFinishedEvent{})))
	})
})
