package deployer_test

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/op/go-logging"
	"github.com/spf13/afero"

	"github.com/compozed/deployadactyl/config"
	. "github.com/compozed/deployadactyl/controller/deployer"
	"github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/mocks"
	"github.com/compozed/deployadactyl/randomizer"
	"github.com/compozed/deployadactyl/state/stop"
	S "github.com/compozed/deployadactyl/structs"
)

const (
	testManifest = `---
applications:
- name: deployadactyl
  memory: 256M
  disk_quota: 256M
`
)

var _ = Describe("Deployer", func() {
	var (
		deployer Deployer

		c              config.Config
		blueGreener    *mocks.BlueGreener
		prechecker     *mocks.Prechecker
		eventManager   *mocks.EventManager
		randomizerMock *mocks.Randomizer

		requestBody                  *bytes.Buffer
		appName                      string
		appPath                      string
		artifactURL                  string
		authorization                interfaces.Authorization
		domain                       string
		environment                  string
		org                          string
		space                        string
		username                     string
		uuid                         string
		manifest                     string
		base64Manifest               string
		instances                    uint16
		password                     string
		testManifestLocation         string
		response                     *bytes.Buffer
		logBuffer                    *Buffer
		log                          interfaces.DeploymentLogger
		deploymentInfo               S.DeploymentInfo
		deploymentInfoNoCustomParams S.DeploymentInfo
		foundations                  []string
		disabledRollback             bool
		environments                 = map[string]S.Environment{}
		environmentsNoCustomParams   = map[string]S.Environment{}
		af                           *afero.Afero
		pusherCreator                *mocks.PushManager
		stopperCreator               interfaces.ActionCreator
		contentType                  string
	)

	BeforeEach(func() {
		blueGreener = &mocks.BlueGreener{}
		prechecker = &mocks.Prechecker{}
		eventManager = &mocks.EventManager{}
		randomizerMock = &mocks.Randomizer{}

		appName = "appName-" + randomizer.StringRunes(10)
		appPath = "appPath-" + randomizer.StringRunes(10)
		artifactURL = "artifactURL-" + randomizer.StringRunes(10)
		domain = "domain-" + randomizer.StringRunes(10)
		environment = "environment-" + randomizer.StringRunes(10)
		org = "org-" + randomizer.StringRunes(10)
		password = "password-" + randomizer.StringRunes(10)
		space = "space-" + randomizer.StringRunes(10)
		username = "username-" + randomizer.StringRunes(10)
		uuid = "uuid-" + randomizer.StringRunes(10)
		instances = uint16(rand.Uint32())
		manifest = "manifest-" + randomizer.StringRunes(10)
		contentType = randomizer.StringRunes(10)

		disabledRollback = false

		base64Manifest = base64.StdEncoding.EncodeToString([]byte(manifest))

		randomizerMock.RandomizeCall.Returns.Runes = uuid
		eventManager.EmitCall.Returns.Error = append(eventManager.EmitCall.Returns.Error, nil)
		eventManager.EmitCall.Returns.Error = append(eventManager.EmitCall.Returns.Error, nil)
		eventManager.EmitCall.Returns.Error = append(eventManager.EmitCall.Returns.Error, nil)
		eventManager.EmitCall.Returns.Error = append(eventManager.EmitCall.Returns.Error, nil)

		requestBody = bytes.NewBufferString(fmt.Sprintf(`{
				"artifact_url": "%s",
				"manifest": "%s"
			}`,
			ccccccjffvthrflerdcilhnhlricrihvcfgvddjddrlt
			URL,
			base64Manifest,
		))

		customParams := make(map[string]interface{})
		customParams["service_now_column_name"] = "u_change"
		customParams["service_now_table_name"] = "u_table"

		deploymentInfo = S.DeploymentInfo{
			ArtifactURL:  artifactURL,
			Username:     username,
			Password:     password,
			Environment:  environment,
			Org:          org,
			Space:        space,
			AppName:      appName,
			UUID:         uuid,
			Instances:    instances,
			Manifest:     manifest,
			Domain:       domain,
			AppPath:      appPath,
			CustomParams: customParams,
			ContentType:  contentType,
			Body:         requestBody,
		}

		deploymentInfoNoCustomParams = S.DeploymentInfo{
			ArtifactURL: artifactURL,
			Username:    username,
			Password:    password,
			Environment: environment,
			Org:         org,
			Space:       space,
			AppName:     appName,
			UUID:        uuid,
			Instances:   instances,
			Manifest:    manifest,
			Domain:      domain,
			AppPath:     appPath,
		}

		foundations = []string{randomizer.StringRunes(10)}
		response = &bytes.Buffer{}

		environments[environment] = S.Environment{
			Name:            environment,
			Domain:          domain,
			Foundations:     foundations,
			Instances:       instances,
			CustomParams:    customParams,
			DisableRollback: disabledRollback,
		}
		authorization.Username = deploymentInfo.Username
		authorization.Password = deploymentInfo.Password

		c = config.Config{
			Username:     username,
			Password:     password,
			Environments: environments,
		}
		logBuffer = NewBuffer()
		log = interfaces.DeploymentLogger{Log: interfaces.DefaultLogger(logBuffer, logging.DEBUG, "deployer tests")}
		pusherCreator = &mocks.PushManager{}
		stopperCreator = stop.StopManager{}

		af = &afero.Afero{Fs: afero.NewMemMapFs()}

		testManifestLocation, _ = af.TempDir("", "")

		deployer = Deployer{
			c,
			blueGreener,
			prechecker,
			eventManager,
			randomizerMock,
			nil,
			log,
		}
	})

	Describe("prechecking the environments", func() {
		Context("when Prechecker fails", func() {
			It("rejects the request with a http.StatusInternalServerError", func() {
				prechecker.AssertAllFoundationsUpCall.Returns.Error = errors.New("prechecker failed")

				deployResponse := deployer.Deploy(&deploymentInfo, deployer.Config.Environments[environment], pusherCreator, response)
				Expect(deployResponse.Error).To(MatchError("prechecker failed"))

				Expect(deployResponse.StatusCode).To(Equal(http.StatusInternalServerError))
				Expect(prechecker.AssertAllFoundationsUpCall.Received.Environment).To(Equal(environments[environment]))
			})
		})
	})

	Describe("authentication", func() {
		Context("a username and password are not provided", func() {
			Context("when authenticate in the config is not true", func() {
				It("uses the config username and password and accepts the request with a http.StatusOK", func() {
					By("setting authenticate to false")
					env := S.Environment{Authenticate: false}
					deployer.Config.Environments[environment] = env
					pusherCreator.OnFinishCall.Returns.DeployResponse = interfaces.DeployResponse{
						StatusCode: http.StatusOK,
					}

					By("not setting basic auth")

					deployResponse := deployer.Deploy(&deploymentInfo, env, pusherCreator, response)

					Expect(deployResponse.Error).ToNot(HaveOccurred())
					Expect(deployResponse.StatusCode).To(Equal(http.StatusOK))
				})
			})
		})
	})

	Describe("deploying with JSON in the request body", func() {
		Context("when manifest is given in the request body", func() {
			Context("if the provided manifest is base64 encoded", func() {
				It("decodes the manifest, does not return an error and returns http.StatusOK", func() {
					manifest = "manifest-" + randomizer.StringRunes(10)
					deploymentInfo.Manifest = fmt.Sprintf(manifest, randomizer.StringRunes(10))

					By("base64 encoding the manifest")
					base64Manifest := base64.StdEncoding.EncodeToString([]byte(deploymentInfo.Manifest))

					By("including the manifest in the request body")
					requestBody = bytes.NewBufferString(fmt.Sprintf(`{"artifact_url": "%s", "manifest": "%s"}`,
						artifactURL,
						base64Manifest,
					))

					pusherCreator.OnFinishCall.Returns.DeployResponse = interfaces.DeployResponse{
						StatusCode: http.StatusOK,
					}

					deployResponse := deployer.Deploy(&deploymentInfo, S.Environment{}, pusherCreator, response)

					Expect(deployResponse.Error).ToNot(HaveOccurred())

					Expect(deployResponse.StatusCode).To(Equal(http.StatusOK))
				})

				It("will emit ArtifactRetrievalStart and ArtifactRetrievalSuccess", func() {
					manifest = "manifest-" + randomizer.StringRunes(10)
					deploymentInfo.Manifest = fmt.Sprintf(manifest, randomizer.StringRunes(10))

					By("base64 encoding the manifest")
					base64Manifest := base64.StdEncoding.EncodeToString([]byte(deploymentInfo.Manifest))

					By("including the manifest in the request body")
					requestBody = bytes.NewBufferString(fmt.Sprintf(`{"artifact_url": "%s", "manifest": "%s"}`,
						artifactURL,
						base64Manifest,
					))

					pusherCreator.OnFinishCall.Returns.DeployResponse = interfaces.DeployResponse{
						StatusCode: http.StatusOK,
					}

					deployResponse := deployer.Deploy(&deploymentInfo, S.Environment{}, pusherCreator, response)

					Expect(deployResponse.Error).ToNot(HaveOccurred())

					Expect(deployResponse.StatusCode).To(Equal(http.StatusOK))
				})
				It("returns statusInternalServerError when actionCreator setup fails", func() {
					manifest = "manifest-" + randomizer.StringRunes(10)

					deploymentInfo.Manifest = fmt.Sprintf(manifest, randomizer.StringRunes(10))
					pusherCreator.SetUpCall.Returns.Err = errors.New("a test error")
					By("base64 encoding the manifest")
					base64Manifest := base64.StdEncoding.EncodeToString([]byte(deploymentInfo.Manifest))

					By("including the manifest in the request body")
					requestBody = bytes.NewBufferString(fmt.Sprintf(`{"artifact_url": "%s", "manifest": "%s"}`,
						artifactURL,
						base64Manifest,
					))

					deployResponse := deployer.Deploy(&deploymentInfo, S.Environment{}, pusherCreator, response)

					Expect(deployResponse.StatusCode).To(Equal(http.StatusInternalServerError))

				})
			})
		})

		Context("when a UUID is provided", func() {
			It("does not create a new UUID", func() {
				manifest = "manifest-" + randomizer.StringRunes(10)
				deploymentInfo.Manifest = fmt.Sprintf(manifest, randomizer.StringRunes(10))
				base64Manifest := base64.StdEncoding.EncodeToString([]byte(deploymentInfo.Manifest))

				pusherCreator.OnFinishCall.Returns.DeployResponse = interfaces.DeployResponse{
					StatusCode: http.StatusOK,
				}

				requestBody = bytes.NewBufferString(fmt.Sprintf(`{"artifact_url": "%s", "manifest": "%s"}`,
					artifactURL,
					base64Manifest,
				))

				deployResponse := deployer.Deploy(&deploymentInfo, S.Environment{}, pusherCreator, response)

				Expect(deployResponse.Error).ToNot(HaveOccurred())

				Expect(deployResponse.StatusCode).To(Equal(http.StatusOK))
				Expect(deployResponse.DeploymentInfo.UUID).To(Equal(uuid))

			})
		})

		Context("when no UUID is provided", func() {
			It("creates a new UUID", func() {
				manifest = "manifest-" + randomizer.StringRunes(10)
				deploymentInfo.Manifest = fmt.Sprintf(manifest, randomizer.StringRunes(10))
				base64Manifest := base64.StdEncoding.EncodeToString([]byte(deploymentInfo.Manifest))

				pusherCreator.OnFinishCall.Returns.DeployResponse = interfaces.DeployResponse{
					StatusCode: http.StatusOK,
				}

				requestBody = bytes.NewBufferString(fmt.Sprintf(`{"artifact_url": "%s", "manifest": "%s"}`,
					artifactURL,
					base64Manifest,
				))

				uuid = ""
				deployResponse := deployer.Deploy(&deploymentInfo, S.Environment{}, pusherCreator, response)

				Expect(deployResponse.Error).ToNot(HaveOccurred())

				Expect(deployResponse.StatusCode).To(Equal(http.StatusOK))
				Expect(deployResponse.DeploymentInfo.UUID).ToNot(Equal(uuid))

			})
		})
	})

	Describe("deploying with a zip file in the request body", func() {
		Context("fetching an artifact from the request body", func() {
			Context("when Fetcher fails", func() {
				It("returns an error and http.StatusInternalServerError", func() {
					pusherCreator.SetUpCall.Returns.Err = errors.New("a test error")

					deployResponse := deployer.Deploy(&deploymentInfo, S.Environment{}, pusherCreator, response)

					Expect(deployResponse.Error.Error()).To(ContainSubstring("a test error"))

					Expect(deployResponse.StatusCode).To(Equal(http.StatusInternalServerError))
				})
			})
		})
	})

	Describe("emitting events during a deployment", func() {
		BeforeEach(func() {
			eventManager.EmitCall.Returns.Error = nil
		})

		Context("when blue greener succeeds", func() {

			It("returns correct deployment info", func() {
				eventManager.EmitCall.Returns.Error = append(eventManager.EmitCall.Returns.Error, nil)
				eventManager.EmitCall.Returns.Error = append(eventManager.EmitCall.Returns.Error, nil)
				eventManager.EmitCall.Returns.Error = append(eventManager.EmitCall.Returns.Error, nil)

				deployResponse := deployer.Deploy(&deploymentInfo, S.Environment{}, pusherCreator, response)

				Expect(deployResponse.DeploymentInfo.UUID).ToNot(Equal(""))
				manifest := deployResponse.DeploymentInfo.Manifest
				Expect(manifest).To(ContainSubstring("manifest-"))
			})
		})
	})

	Describe("happy path deploying with json in the request body", func() {
		Context("when no errors occur", func() {
			It("accepts the request and returns http.StatusOK", func() {
				deploymentInfo.ContentType = "JSON"
				pusherCreator.OnFinishCall.Returns.DeployResponse = interfaces.DeployResponse{
					StatusCode: http.StatusOK,
				}

				deployResponse := deployer.Deploy(&deploymentInfo, environments[environment], pusherCreator, response)

				Expect(deployResponse.Error).To(BeNil())

				Expect(deployResponse.StatusCode).To(Equal(http.StatusOK))

				Eventually(logBuffer).Should(Say("prechecking the foundations"))
				Expect(prechecker.AssertAllFoundationsUpCall.Received.Environment).To(Equal(environments[environment]))
				Expect(blueGreener.ExecuteCall.Received.Environment).To(Equal(environments[environment]))
			})
		})
	})

	Describe("happy path deploying with a zip file in the request body", func() {
		Context("when no errors occur", func() {
			It("accepts the request and returns http.StatusOK", func() {
				Expect(af.WriteFile(testManifestLocation+"/manifest.yml", []byte(testManifest), 0644)).To(Succeed())

				pusherCreator.OnFinishCall.Returns.DeployResponse = interfaces.DeployResponse{
					StatusCode: http.StatusOK,
				}

				deployResponse := deployer.Deploy(&deploymentInfo, environments[environment], pusherCreator, response)
				Expect(deployResponse.Error).To(BeNil())

				Expect(deployResponse.StatusCode).To(Equal(http.StatusOK))

				Eventually(logBuffer).Should(Say("prechecking the foundations"))

				Expect(prechecker.AssertAllFoundationsUpCall.Received.Environment).To(Equal(environments[environment]))
				Expect(blueGreener.ExecuteCall.Received.Environment).To(Equal(environments[environment]))
			})
		})
	})

	Describe("extract custom params from yaml", func() {

		Context("when no custom params are provided", func() {
			BeforeEach(func() {
				environmentsNoCustomParams[environment] = S.Environment{
					Name:        environment,
					Domain:      domain,
					Foundations: foundations,
					Instances:   instances,
				}

				c := config.Config{
					Username:     username,
					Password:     password,
					Environments: environmentsNoCustomParams,
				}

				deployer = Deployer{
					c,
					blueGreener,
					prechecker,
					eventManager,
					randomizerMock,
					nil,
					log,
				}
			})

			It("doesn't return an error", func() {

				deployResponse := deployer.Deploy(&deploymentInfoNoCustomParams, environmentsNoCustomParams[environment], pusherCreator, response)

				Expect(deployResponse.Error).ToNot(HaveOccurred())
				Expect(blueGreener.ExecuteCall.Received.Environment).To(Equal(environmentsNoCustomParams[environment]))
			})
		})
	})
	Describe("Deploy", func() {
		var (
			deployer          interfaces.Deployer
			pusherCreatorMock *mocks.PushManager
		)
		BeforeEach(func() {
			pusherCreatorMock = &mocks.PushManager{}
			deployer = Deployer{
				c,
				blueGreener,
				prechecker,
				eventManager,
				randomizerMock,
				nil,
				log,
			}
		})
		Context("when no initialization errors occur", func() {
			It("it calls setup on the provided action creator", func() {

				deployer.Deploy(&deploymentInfo, S.Environment{}, pusherCreatorMock, response)

				Expect(pusherCreatorMock.SetUpCall.Called).To(Equal(true))
			})
		})

		It("calls Start on the provided action creator", func() {

			deployer.Deploy(&deploymentInfo, S.Environment{}, pusherCreatorMock, response)

			Expect(pusherCreatorMock.OnStartCall.Called).To(Equal(true))
		})

		Context("when action creator OnStart fails", func() {
			It("returns an error", func() {
				pusherCreatorMock.OnStartCall.Returns.Err = errors.New("a test error")

				deployResponse := deployer.Deploy(&deploymentInfo, S.Environment{}, pusherCreatorMock, response)

				Expect(deployResponse.Error).To(Equal(pusherCreatorMock.OnStartCall.Returns.Err))
			})
		})

		It("calls CleanUp on the provided action creator", func() {
			deployer.Deploy(&deploymentInfo, S.Environment{}, pusherCreatorMock, response)

			Expect(pusherCreatorMock.CleanUpCall.Called).To(Equal(true))
		})

		It("calls OnFinish on the provided action creator", func() {
			deployer.Deploy(&deploymentInfo, S.Environment{}, pusherCreatorMock, response)

			Expect(pusherCreatorMock.OnFinishCall.Called).To(Equal(true))
		})
	})
})
