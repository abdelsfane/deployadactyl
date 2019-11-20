package stop_test

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/compozed/deployadactyl/config"
	D "github.com/compozed/deployadactyl/controller/deployer"
	"github.com/compozed/deployadactyl/controller/deployer/error_finder"
	I "github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/mocks"
	"github.com/compozed/deployadactyl/randomizer"
	"github.com/compozed/deployadactyl/request"
	"github.com/compozed/deployadactyl/state"
	. "github.com/compozed/deployadactyl/state/stop"
	"github.com/compozed/deployadactyl/structs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/op/go-logging"
	"net/http"
	"reflect"
)

var _ = Describe("StopDeployment", func() {
	var (
		deployer           *mocks.Deployer
		pushManagerFactory *mocks.PushManagerFactory
		stopManagerFactory *mocks.StopManagerFactory
		eventManager       *mocks.EventManager
		errorFinder        *mocks.ErrorFinder
		controller         *StopController
		deployment         I.Deployment
		authResolver       *state.AuthResolver
		envResolver        *state.EnvResolver
		logBuffer          *Buffer

		appName     string
		environment string
		org         string
		space       string
		response    *bytes.Buffer
	)

	BeforeEach(func() {
		logBuffer = NewBuffer()
		appName = "appName-" + randomizer.StringRunes(10)
		environment = "environment-" + randomizer.StringRunes(10)
		org = "org-" + randomizer.StringRunes(10)
		space = "non-prod"

		eventManager = &mocks.EventManager{}
		deployer = &mocks.Deployer{}
		pushManagerFactory = &mocks.PushManagerFactory{}

		authResolver = &state.AuthResolver{Config: config.Config{}}
		envResolver = &state.EnvResolver{Config: config.Config{}}

		stopManagerFactory = &mocks.StopManagerFactory{}
		errorFinder = &mocks.ErrorFinder{}
		controller = &StopController{
			Deployer:           deployer,
			Log:                I.DeploymentLogger{Log: I.DefaultLogger(logBuffer, logging.DEBUG, "api_test"), UUID: randomizer.StringRunes(10)},
			StopManagerFactory: stopManagerFactory,
			EventManager:       eventManager,
			AuthResolver:       authResolver,
			ErrorFinder:        errorFinder,
			EnvResolver:        envResolver,
		}
		environments := map[string]structs.Environment{}
		environments[environment] = structs.Environment{}
		envResolver.Config.Environments = environments
		bodyByte := []byte("{}")
		response = &bytes.Buffer{}

		deployment = I.Deployment{
			Body:          &bodyByte,
			Type:          "application/json",
			CFContext:     I.CFContext{},
			Authorization: I.Authorization{},
		}

	})

	Context("When UUID is not provided", func() {
		It("Should populate UUID", func() {

			deployment := &I.Deployment{
				CFContext: I.CFContext{
					Environment: environment,
				}}
			response := bytes.NewBuffer([]byte{})
			putDeploymentRequest := request.PutDeploymentRequest{
				Deployment: *deployment,
				Request:    request.PutRequest{Data: nil},
			}

			deploymentResponse := controller.StopDeployment(putDeploymentRequest, response)

			Expect(deploymentResponse.DeploymentInfo.UUID).ShouldNot(BeEmpty())
		})
	})
	It("Should return org, space, appname, and environment when provided", func() {

		deployment := &I.Deployment{
			CFContext: I.CFContext{
				Organization: "myOrg",
				Space:        "mySpace",
				Application:  "myApp",
				Environment:  environment,
			},
		}
		response := bytes.NewBuffer([]byte{})
		putDeploymentRequest := request.PutDeploymentRequest{
			Deployment: *deployment,
			Request:    request.PutRequest{Data: nil},
		}

		deploymentResponse := controller.StopDeployment(putDeploymentRequest, response)

		Expect(deploymentResponse.DeploymentInfo.Org).Should(Equal("myOrg"))
		Expect(deploymentResponse.DeploymentInfo.Environment).Should(Equal(environment))
		Expect(deploymentResponse.DeploymentInfo.Space).Should(Equal("mySpace"))
		Expect(deploymentResponse.DeploymentInfo.AppName).Should(Equal("myApp"))

	})
	It("Should log start of process", func() {

		deployment := &I.Deployment{
			CFContext: I.CFContext{
				Application: "myApp",
				Environment: environment,
			},
		}
		response := bytes.NewBuffer([]byte{})
		putDeploymentRequest := request.PutDeploymentRequest{
			Deployment: *deployment,
			Request:    request.PutRequest{Data: nil},
		}

		deploymentResponse := controller.StopDeployment(putDeploymentRequest, response)

		Expect(logBuffer).Should(Say(fmt.Sprintf("Preparing to stop %s with UUID %s", "myApp", deploymentResponse.DeploymentInfo.UUID)))

	})

	Context("When StopStartEvent succeeds", func() {
		It("should emit a StopStarteEvent", func() {
			data := make(map[string]interface{})
			data["mykey"] = "first value"
			deployment := &I.Deployment{
				CFContext: I.CFContext{
					Organization: "myOrg",
					Space:        "mySpace",
					Application:  "myApp",
					Environment:  environment,
				},
			}

			putDeploymentRequest := request.PutDeploymentRequest{
				Deployment: *deployment,
				Request:    request.PutRequest{Data: data},
			}

			controller.StopDeployment(putDeploymentRequest, response)

			Expect(reflect.TypeOf(eventManager.EmitEventCall.Received.Events[0])).Should(Equal(reflect.TypeOf(StopStartedEvent{})))
			stopEvent := eventManager.EmitEventCall.Received.Events[0].(StopStartedEvent)
			Expect(stopEvent.CFContext.Space).Should(Equal("mySpace"))
			Expect(stopEvent.CFContext.Application).Should(Equal("myApp"))
			Expect(stopEvent.CFContext.Environment).Should(Equal(environment))
			Expect(stopEvent.CFContext.Organization).Should(Equal("myOrg"))
			Expect(stopEvent.Data).Should(Equal(data))

		})
	})

	Context("When StopStartEvent fails", func() {
		It("should return error", func() {
			eventManager.EmitEventCall.Returns.Error = []error{errors.New("anything")}

			deployment := &I.Deployment{
				CFContext: I.CFContext{
					Environment: environment,
				},
			}
			putDeploymentRequest := request.PutDeploymentRequest{
				Deployment: *deployment,
				Request:    request.PutRequest{Data: nil},
			}

			deployResponse := controller.StopDeployment(putDeploymentRequest, response)

			Expect(deployResponse.StatusCode).Should(Equal(http.StatusInternalServerError))
			Expect(reflect.TypeOf(deployResponse.Error)).Should(Equal(reflect.TypeOf(D.EventError{})))

		})
	})

	Context("When environment does not exist", func() {
		It("Should return error", func() {

			deployment := &I.Deployment{
				CFContext: I.CFContext{
					Environment: "bad environment",
				}}
			response := bytes.NewBuffer([]byte{})
			putDeploymentRequest := request.PutDeploymentRequest{
				Deployment: *deployment,
				Request:    request.PutRequest{Data: nil},
			}

			deploymentResponse := controller.StopDeployment(putDeploymentRequest, response)

			Expect(reflect.TypeOf(deploymentResponse.Error)).Should(Equal(reflect.TypeOf(D.EnvironmentNotFoundError{})))
		})
	})

	Context("When environment exists", func() {
		It("Should return SkipSSL, CustomParams, and Domain", func() {

			envResolver.Config.Environments[environment] = structs.Environment{
				SkipSSL:      true,
				Domain:       "myDomain",
				CustomParams: make(map[string]interface{}),
			}
			envResolver.Config.Environments[environment].CustomParams["customName"] = "customParams"

			deployment := &I.Deployment{
				CFContext: I.CFContext{
					Environment: environment,
				}}
			response := bytes.NewBuffer([]byte{})
			putDeploymentRequest := request.PutDeploymentRequest{
				Deployment: *deployment,
				Request:    request.PutRequest{Data: nil},
			}

			deploymentResponse := controller.StopDeployment(putDeploymentRequest, response)
			Expect(deploymentResponse.DeploymentInfo.Domain).Should(Equal("myDomain"))
			Expect(deploymentResponse.DeploymentInfo.SkipSSL).Should(Equal(true))
			Expect(deploymentResponse.DeploymentInfo.CustomParams["customName"]).Should(Equal("customParams"))
		})
	})

	Context("When auth does not exist", func() {
		Context("When environment authenticate is true", func() {
			It("Should return error", func() {
				envResolver.Config.Environments[environment] = structs.Environment{
					Authenticate: true,
				}
				deployment := &I.Deployment{
					CFContext: I.CFContext{
						Environment: environment,
					}}
				response := bytes.NewBuffer([]byte{})
				putDeploymentRequest := request.PutDeploymentRequest{
					Deployment: *deployment,
					Request:    request.PutRequest{Data: nil},
				}

				deploymentResponse := controller.StopDeployment(putDeploymentRequest, response)

				Expect(reflect.TypeOf(deploymentResponse.Error)).Should(Equal(reflect.TypeOf(D.BasicAuthError{})))
			})
		})

		Context("When environment authenticate is false", func() {
			It("Should username and password using the config", func() {
				authResolver.Config.Username = "username"
				authResolver.Config.Password = "password"
				envResolver.Config.Environments[environment] = structs.Environment{
					Authenticate: false,
				}
				deployment := &I.Deployment{
					CFContext: I.CFContext{
						Environment: environment,
					}}
				response := bytes.NewBuffer([]byte{})
				putDeploymentRequest := request.PutDeploymentRequest{
					Deployment: *deployment,
					Request:    request.PutRequest{Data: nil},
				}

				deploymentResponse := controller.StopDeployment(putDeploymentRequest, response)

				Expect(deploymentResponse.DeploymentInfo.Username).Should(Equal("username"))
				Expect(deploymentResponse.DeploymentInfo.Password).Should(Equal("password"))
			})
		})
	})

	Context("When auth is provided", func() {
		It("Should populate the deploymentInfo with the username and password", func() {
			deployment := &I.Deployment{
				Authorization: I.Authorization{
					Username: "myUser",
					Password: "myPassword",
				},
				CFContext: I.CFContext{
					Environment: environment,
				},
			}
			response := bytes.NewBuffer([]byte{})
			putDeploymentRequest := request.PutDeploymentRequest{
				Deployment: *deployment,
				Request:    request.PutRequest{Data: nil},
			}

			deploymentResponse := controller.StopDeployment(putDeploymentRequest, response)
			Expect(deploymentResponse.DeploymentInfo.Username).Should(Equal("myUser"))
			Expect(deploymentResponse.DeploymentInfo.Password).Should(Equal("myPassword"))
		})
	})

	Context("When auth is provided", func() {
		It("Should populate the deploymentInfo with the username and password", func() {
			deployment := &I.Deployment{
				Authorization: I.Authorization{
					Username: "myUser",
					Password: "myPassword",
				},
				CFContext: I.CFContext{
					Environment: environment,
				},
			}
			response := bytes.NewBuffer([]byte{})
			putDeploymentRequest := request.PutDeploymentRequest{
				Deployment: *deployment,
				Request:    request.PutRequest{Data: nil},
			}

			deploymentResponse := controller.StopDeployment(putDeploymentRequest, response)
			Expect(deploymentResponse.DeploymentInfo.Username).Should(Equal("myUser"))
			Expect(deploymentResponse.DeploymentInfo.Password).Should(Equal("myPassword"))
		})
	})

	Context("When data is provided", func() {
		It("should return deployment info with proper data", func() {
			data := map[string]interface{}{
				"user_id": "myuserid",
				"group":   "mygroup",
			}
			deployment := &I.Deployment{
				CFContext: I.CFContext{
					Environment: environment,
				},
			}
			response := bytes.NewBuffer([]byte{})
			putDeploymentRequest := request.PutDeploymentRequest{
				Deployment: *deployment,
				Request:    request.PutRequest{Data: data},
			}

			deploymentResponse := controller.StopDeployment(putDeploymentRequest, response)
			Expect(deploymentResponse.DeploymentInfo.Data["user_id"]).Should(Equal("myuserid"))
			Expect(deploymentResponse.DeploymentInfo.Data["group"]).Should(Equal("mygroup"))

		})
	})
	It("should create stop manager", func() {

		deployment := &I.Deployment{
			Authorization: I.Authorization{
				Username: "myUser",
			},
			CFContext: I.CFContext{
				Environment: environment,
			},
		}
		response := bytes.NewBuffer([]byte{})
		putDeploymentRequest := request.PutDeploymentRequest{
			Deployment: *deployment,
			Request:    request.PutRequest{Data: nil},
		}

		controller.StopDeployment(putDeploymentRequest, response)
		Expect(stopManagerFactory.StopManagerCall.Called).Should(Equal(true))
		Expect(stopManagerFactory.StopManagerCall.Received.DeployEventData.DeploymentInfo.Username).Should(Equal("myUser"))
	})
	It("should call deploy with the stop manager ", func() {
		manager := &mocks.StopManager{}
		stopManagerFactory.StopManagerCall.Returns.ActionCreater = manager
		deployment := &I.Deployment{
			CFContext: I.CFContext{
				Environment: environment,
			},
		}
		putDeploymentRequest := request.PutDeploymentRequest{
			Deployment: *deployment,
			Request:    request.PutRequest{Data: nil},
		}

		response := bytes.NewBuffer([]byte{})
		controller.StopDeployment(putDeploymentRequest, response)
		Expect(deployer.DeployCall.Received.ActionCreator).Should(Equal(manager))
	})
	It("should call deploy with the stop manager ", func() {
		deployer.DeployCall.Returns.Error = errors.New("test error")
		deployer.DeployCall.Returns.StatusCode = http.StatusOK

		deployment := &I.Deployment{
			CFContext: I.CFContext{
				Environment: environment,
			},
		}
		response := bytes.NewBuffer([]byte{})
		putDeploymentRequest := request.PutDeploymentRequest{
			Deployment: *deployment,
			Request:    request.PutRequest{Data: nil},
		}

		deploymentResponse := controller.StopDeployment(putDeploymentRequest, response)

		Expect(deploymentResponse.Error.Error()).Should(Equal("test error"))
		Expect(deploymentResponse.StatusCode).Should(Equal(http.StatusOK))

	})

	Context("when stop succeeds", func() {
		Context("if StopSuccessEvent succeeds", func() {
			It("should emit StopSuccessEvent", func() {
				data := make(map[string]interface{})
				data["mykey"] = "first value"

				deployment := &I.Deployment{
					CFContext: I.CFContext{
						Organization: "myOrg",
						Space:        "mySpace",
						Application:  "myApp",
						Environment:  environment,
					},
					Authorization: I.Authorization{
						Username: "myUser",
						Password: "myPassword",
					},
				}
				response := bytes.NewBuffer([]byte{})

				envResolver.Config.Environments[environment] = structs.Environment{
					Name:         environment,
					Authenticate: true,
				}
				putDeploymentRequest := request.PutDeploymentRequest{
					Deployment: *deployment,
					Request:    request.PutRequest{Data: data},
				}

				controller.StopDeployment(putDeploymentRequest, response)

				Expect(reflect.TypeOf(eventManager.EmitEventCall.Received.Events[1])).To(Equal(reflect.TypeOf(StopSuccessEvent{})))
				stopSuccessEvent := eventManager.EmitEventCall.Received.Events[1].(StopSuccessEvent)

				Expect(stopSuccessEvent.CFContext.Space).Should(Equal("mySpace"))
				Expect(stopSuccessEvent.CFContext.Application).Should(Equal("myApp"))
				Expect(stopSuccessEvent.CFContext.Environment).Should(Equal(environment))
				Expect(stopSuccessEvent.CFContext.Organization).Should(Equal("myOrg"))
				Expect(stopSuccessEvent.Authorization.Username).Should(Equal("myUser"))
				Expect(stopSuccessEvent.Authorization.Password).Should(Equal("myPassword"))
				Expect(stopSuccessEvent.Environment.Name).Should(Equal(environment))
				Expect(stopSuccessEvent.Data).Should(Equal(data))

			})
			It("should emit a StopStartedEvent", func() {
				data := make(map[string]interface{})
				data["mykey"] = "first value"

				deployment := &I.Deployment{
					CFContext: I.CFContext{
						Organization: "myOrg",
						Space:        "mySpace",
						Application:  "myApp",
						Environment:  environment,
					},
				}
				putDeploymentRequest := request.PutDeploymentRequest{
					Deployment: *deployment,
					Request:    request.PutRequest{Data: data},
				}

				controller.StopDeployment(putDeploymentRequest, response)

				Expect(reflect.TypeOf(eventManager.EmitEventCall.Received.Events[0])).Should(Equal(reflect.TypeOf(StopStartedEvent{})))
				stopEvent := eventManager.EmitEventCall.Received.Events[0].(StopStartedEvent)
				Expect(stopEvent.CFContext.Space).Should(Equal("mySpace"))
				Expect(stopEvent.CFContext.Application).Should(Equal("myApp"))
				Expect(stopEvent.CFContext.Environment).Should(Equal(environment))
				Expect(stopEvent.CFContext.Organization).Should(Equal("myOrg"))
				Expect(stopEvent.Data).Should(Equal(data))

			})
		})
		Context("if StopSuccessEvent fails", func() {
			It("should log the error", func() {
				eventManager.EmitEventCall.Returns.Error = []error{nil, errors.New("errors")}
				deployment := &I.Deployment{
					CFContext: I.CFContext{
						Environment: environment,
					},
				}
				response := bytes.NewBuffer([]byte{})
				putDeploymentRequest := request.PutDeploymentRequest{
					Deployment: *deployment,
					Request:    request.PutRequest{Data: nil},
				}

				controller.StopDeployment(putDeploymentRequest, response)

				Eventually(logBuffer).Should(Say("an error occurred when emitting a StopSuccessEvent event: errors"))
			})
		})

	})

	Context("when stop fails", func() {
		It("print errors", func() {
			deployment := &I.Deployment{
				CFContext: I.CFContext{
					Environment: environment,
				},
			}
			deployer.DeployCall.Returns.Error = errors.New("deploy error")
			errorFinder.FindErrorsCall.Returns.Errors = []I.LogMatchedError{error_finder.CreateLogMatchedError("a test error", []string{"error 1", "error 2", "error 3"}, "error solution", "test code")}
			response := bytes.NewBuffer([]byte{})
			putDeploymentRequest := request.PutDeploymentRequest{
				Deployment: *deployment,
				Request:    request.PutRequest{Data: nil},
			}

			controller.StopDeployment(putDeploymentRequest, response)
			Eventually(response).Should(ContainSubstring("Potential solution"))
		})
		It("should emit StopFailureEvent", func() {
			data := make(map[string]interface{})
			data["mykey"] = "first value"

			deployment := &I.Deployment{
				CFContext: I.CFContext{
					Organization: "myOrg",
					Space:        "mySpace",
					Application:  "myApp",
					Environment:  environment,
				},
				Authorization: I.Authorization{
					Username: "myUser",
					Password: "myPassword",
				},
			}
			response := bytes.NewBuffer([]byte{})

			envResolver.Config.Environments[environment] = structs.Environment{
				Name:         environment,
				Authenticate: true,
			}
			deployer.DeployCall.Returns.Error = errors.New("deploy error")
			putDeploymentRequest := request.PutDeploymentRequest{
				Deployment: *deployment,
				Request:    request.PutRequest{Data: data},
			}

			controller.StopDeployment(putDeploymentRequest, response)

			Expect(reflect.TypeOf(eventManager.EmitEventCall.Received.Events[1])).To(Equal(reflect.TypeOf(StopFailureEvent{})))
			event := eventManager.EmitEventCall.Received.Events[1].(StopFailureEvent)

			Expect(event.CFContext.Space).Should(Equal("mySpace"))
			Expect(event.CFContext.Application).Should(Equal("myApp"))
			Expect(event.CFContext.Environment).Should(Equal(environment))
			Expect(event.CFContext.Organization).Should(Equal("myOrg"))
			Expect(event.Authorization.Username).Should(Equal("myUser"))
			Expect(event.Authorization.Password).Should(Equal("myPassword"))
			Expect(event.Environment.Name).Should(Equal(environment))
			Expect(event.Data).Should(Equal(data))
			Expect(event.Error.Error()).Should(Equal("deploy error"))

		})
		Context("if StopFailureEvent fails", func() {
			It("should log the error", func() {
				eventManager.EmitEventCall.Returns.Error = []error{nil, errors.New("errors")}
				deployment := &I.Deployment{
					CFContext: I.CFContext{
						Environment: environment,
					},
				}
				deployer.DeployCall.Returns.Error = errors.New("deploy error")

				response := bytes.NewBuffer([]byte{})
				putDeploymentRequest := request.PutDeploymentRequest{
					Deployment: *deployment,
					Request:    request.PutRequest{Data: nil},
				}

				controller.StopDeployment(putDeploymentRequest, response)

				Eventually(logBuffer).Should(Say("an error occurred when emitting a StopFailureEvent event: errors"))
			})
		})

	})

	Context("when stop finishes", func() {
		It("should log an emit StopFinish event", func() {
			deployment := &I.Deployment{
				CFContext: I.CFContext{
					Environment: environment,
				},
			}
			response := bytes.NewBuffer([]byte{})
			putDeploymentRequest := request.PutDeploymentRequest{
				Deployment: *deployment,
				Request:    request.PutRequest{Data: nil},
			}

			controller.StopDeployment(putDeploymentRequest, response)

			Eventually(logBuffer).Should(Say("emitting a StopFinishedEvent"))
		})
		It("should emit StopFinishedEvent", func() {
			data := make(map[string]interface{})
			data["mykey"] = "first value"

			deployment := &I.Deployment{
				CFContext: I.CFContext{
					Organization: "myOrg",
					Space:        "mySpace",
					Application:  "myApp",
					Environment:  environment,
				},
				Authorization: I.Authorization{
					Username: "myUser",
					Password: "myPassword",
				},
			}
			response := bytes.NewBuffer([]byte{})

			envResolver.Config.Environments[environment] = structs.Environment{
				Name:         environment,
				Authenticate: true,
			}
			putDeploymentRequest := request.PutDeploymentRequest{
				Deployment: *deployment,
				Request:    request.PutRequest{Data: data},
			}

			controller.StopDeployment(putDeploymentRequest, response)

			Expect(reflect.TypeOf(eventManager.EmitEventCall.Received.Events[2])).To(Equal(reflect.TypeOf(StopFinishedEvent{})))
			event := eventManager.EmitEventCall.Received.Events[2].(StopFinishedEvent)

			Expect(event.CFContext.Space).Should(Equal("mySpace"))
			Expect(event.CFContext.Application).Should(Equal("myApp"))
			Expect(event.CFContext.Environment).Should(Equal(environment))
			Expect(event.CFContext.Organization).Should(Equal("myOrg"))
			Expect(event.Authorization.Username).Should(Equal("myUser"))
			Expect(event.Authorization.Password).Should(Equal("myPassword"))
			Expect(event.Environment.Name).Should(Equal(environment))
			Expect(event.Data).Should(Equal(data))

		})
	})
})
