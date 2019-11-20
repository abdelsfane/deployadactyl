package creator

import (
	"bytes"

	"reflect"

	"strconv"

	"github.com/compozed/deployadactyl/artifetcher"
	"github.com/compozed/deployadactyl/artifetcher/extractor"
	"github.com/compozed/deployadactyl/config"
	"github.com/compozed/deployadactyl/controller/deployer"
	"github.com/compozed/deployadactyl/controller/deployer/bluegreen"
	"github.com/compozed/deployadactyl/eventmanager"
	"github.com/compozed/deployadactyl/eventmanager/handlers/healthchecker"
	"github.com/compozed/deployadactyl/eventmanager/handlers/routemapper"
	I "github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/mocks"
	"github.com/compozed/deployadactyl/request"
	"github.com/compozed/deployadactyl/state/push"
	"github.com/compozed/deployadactyl/state/start"
	"github.com/compozed/deployadactyl/state/stop"
	"github.com/compozed/deployadactyl/structs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

var _ = Describe("RequestCreator", func() {

	Describe("CreateDeployer", func() {

		Context("when mock constructor is provided", func() {
			It("should return the mock implementation", func() {
				expected := &mocks.Deployer{}
				creator := Creator{
					provider: CreatorModuleProvider{
						NewDeployer: func(cnfg config.Config, blueGreener I.BlueGreener, preChecker I.Prechecker, eventManager I.EventManager, randomizer I.Randomizer, errorFinder I.ErrorFinder, logger I.DeploymentLogger) I.Deployer {
							return expected
						},
					},
				}
				rc := RequestCreator{
					Creator: creator,
				}
				controller := rc.CreateDeployer()
				Expect(controller).To(Equal(expected))
			})
		})

		Context("when mock constructor is not provided", func() {
			It("should return the default implementation", func() {
				creator := Creator{}
				rc := RequestCreator{
					Creator:      creator,
					Log:          I.DeploymentLogger{UUID: "the uuid"},
					EventManager: &mocks.EventManager{},
				}
				actual := rc.CreateDeployer()
				Expect(reflect.TypeOf(actual)).To(Equal(reflect.TypeOf(&deployer.Deployer{})))
				concrete := actual.(*deployer.Deployer)
				Expect(concrete.Config).ToNot(BeNil())
				Expect(concrete.BlueGreener).ToNot(BeNil())
				Expect(concrete.Prechecker).ToNot(BeNil())
				Expect(concrete.EventManager).To(Equal(rc.EventManager))
				Expect(concrete.Randomizer).ToNot(BeNil())
				Expect(concrete.ErrorFinder).ToNot(BeNil())
				Expect(concrete.Log.UUID).To(Equal("the uuid"))
			})
		})
	})

	Describe("CreateBlueGreener", func() {
		Context("when mock constructor is provided", func() {
			It("should return the mock implementation", func() {
				expected := &mocks.BlueGreener{}
				creator := Creator{
					provider: CreatorModuleProvider{
						NewBlueGreen: func(logger I.DeploymentLogger) I.BlueGreener {
							return expected
						},
					},
				}
				rc := RequestCreator{
					Creator: creator,
				}
				greener := rc.CreateBlueGreener()
				Expect(greener).To(Equal(expected))
			})
		})

		Context("when mock constructor is not provided", func() {
			It("should return the default implementation", func() {
				creator := Creator{}
				rc := RequestCreator{
					Creator:      creator,
					Log:          I.DeploymentLogger{UUID: "the uuid"},
					EventManager: &mocks.EventManager{},
				}
				actual := rc.CreateBlueGreener()
				Expect(reflect.TypeOf(actual)).To(Equal(reflect.TypeOf(&bluegreen.BlueGreen{})))
				concrete := actual.(*bluegreen.BlueGreen)
				Expect(concrete.Log.UUID).To(Equal("the uuid"))
			})
		})
	})

	Describe("CreateFetcher", func() {

		Context("when mock constructor is provided", func() {
			It("should return the mock implementation", func() {
				expected := &mocks.Fetcher{}
				creator := Creator{
					provider: CreatorModuleProvider{
						NewFetcher: func(fs *afero.Afero, ex I.Extractor, log I.DeploymentLogger) I.Fetcher {
							return expected
						},
					},
				}
				rc := RequestCreator{
					Creator: creator,
				}
				controller := rc.CreateFetcher()
				Expect(controller).To(Equal(expected))
			})
		})

		Context("when mock constructor is not provided", func() {
			It("should return the default implementation", func() {
				creator := Creator{
					fileSystem: &afero.Afero{},
				}
				rc := RequestCreator{
					Creator:      creator,
					Log:          I.DeploymentLogger{UUID: "the uuid"},
					EventManager: &mocks.EventManager{},
				}
				fetcher := rc.CreateFetcher()
				Expect(reflect.TypeOf(fetcher)).To(Equal(reflect.TypeOf(&artifetcher.Artifetcher{})))
				concrete := fetcher.(*artifetcher.Artifetcher)
				Expect(concrete.Log.UUID).To(Equal(rc.Log.UUID))
				Expect(concrete.Extractor).ToNot(BeNil())
				Expect(concrete.FileSystem).To(Equal(creator.fileSystem))
			})
		})
	})

	Describe("CreateExtractor", func() {

		Context("when mock constructor is provided", func() {
			It("should return the mock implementation", func() {
				expected := &mocks.Extractor{}
				creator := Creator{
					provider: CreatorModuleProvider{
						NewExtractor: func(log I.DeploymentLogger, fs *afero.Afero) I.Extractor {
							return expected
						},
					},
				}
				rc := RequestCreator{
					Creator: creator,
				}
				controller := rc.CreateExtractor()
				Expect(controller).To(Equal(expected))
			})
		})

		Context("when mock constructor is not provided", func() {
			It("should return the default implementation", func() {
				creator := Creator{
					fileSystem: &afero.Afero{},
				}
				rc := RequestCreator{
					Creator:      creator,
					Log:          I.DeploymentLogger{UUID: "the uuid"},
					EventManager: &mocks.EventManager{},
				}
				e := rc.CreateExtractor()
				Expect(reflect.TypeOf(e)).To(Equal(reflect.TypeOf(&extractor.Extractor{})))
				concrete := e.(*extractor.Extractor)
				Expect(concrete.Log.UUID).To(Equal(rc.Log.UUID))
				Expect(concrete.FileSystem).To(Equal(creator.fileSystem))
			})
		})
	})

	Describe("newRequestCreator", func() {

		Context("when EventManager constructor is provided", func() {
			It("should return with the provided EventManager", func() {
				log, _ := NewLogger()

				expectedEventManager := eventmanager.EventManager{
					Log: log,
				}

				creator, err := New(CreatorModuleProvider{
					NewConfig: func() (config.Config, error) {
						return config.Config{}, nil
					},
					NewEventManager: func(logger I.DeploymentLogger, bindings []I.Binding) I.EventManager {
						return &expectedEventManager
					},
				})

				rc := newRequestCreator(creator, "the uuid", bytes.NewBuffer([]byte("")))

				Expect(err).ToNot(HaveOccurred())
				Expect(rc.EventManager).To(Equal(&expectedEventManager))
			})
		})

		Context("when EventManager constructor is not provided", func() {
			It("should return the default EventManager", func() {
				creator, err := New(CreatorModuleProvider{
					NewConfig: func() (config.Config, error) {
						return config.Config{}, nil
					},
				})

				rc := newRequestCreator(creator, "the uuid", bytes.NewBuffer([]byte("")))

				Expect(err).ToNot(HaveOccurred())
				Expect(reflect.TypeOf(rc.EventManager)).To(Equal(reflect.TypeOf(&eventmanager.EventManager{})))
			})
		})

		It("should populate EventManager with bindings", func() {
			creator, _ := New(CreatorModuleProvider{
				NewConfig: func() (config.Config, error) {
					return config.Config{}, nil
				},
			})

			for i := 0; i < 3; i++ {
				binding := &mocks.EventBinding{}
				binding.EmitCall.Received.Event = "event " + strconv.Itoa(i)
				creator.GetEventBindings().AddBinding(binding)
			}

			rc := newRequestCreator(creator, "the uuid", bytes.NewBuffer([]byte("")))

			Expect(rc.EventManager.(*eventmanager.EventManager).Bindings).To(Equal(creator.GetEventBindings().GetBindings()))
		})

		Context("when adding bindings to event manager", func() {
			It("doesn't modify the original bindings", func() {
				c, _ := New(CreatorModuleProvider{
					NewConfig: func() (config.Config, error) {
						return config.Config{}, nil
					},
				})
				c.GetEventBindings().AddBinding(&mocks.EventBinding{})
				c.GetEventBindings().AddBinding(&mocks.EventBinding{})

				rc := newRequestCreator(c, "the uuid", bytes.NewBuffer([]byte("")))
				rc.EventManager.AddBinding(&mocks.EventBinding{})

				Expect(len(c.GetEventBindings().GetBindings())).To(Equal(2))
			})
		})
	})

	Describe("PushRequestCreator", func() {

		Describe("CreateRequestProcessor", func() {
			Context("when mock constructor is provided", func() {
				It("should return the mock implementation", func() {

					expected := &mocks.RequestProcessor{}
					creator := Creator{
						provider: CreatorModuleProvider{
							NewPushRequestProcessor: func(log I.DeploymentLogger, pc request.PushController, request request.PostDeploymentRequest, buffer *bytes.Buffer) I.RequestProcessor {
								return expected
							},
						},
					}
					rc := PushRequestCreator{
						RequestCreator: RequestCreator{
							Creator: creator,
						},
					}
					processor := rc.CreateRequestProcessor()
					Expect(processor).To(Equal(expected))
				})
			})

			Context("when mock constructor is not provided", func() {
				It("should return the default implementation", func() {

					response := bytes.NewBuffer([]byte("the response"))
					request := request.PostDeploymentRequest{
						Deployment: I.Deployment{
							CFContext: I.CFContext{
								Organization: "the org",
							},
						},
					}

					rc := PushRequestCreator{
						RequestCreator: RequestCreator{
							Buffer: response,
							Log:    I.DeploymentLogger{UUID: "the uuid"},
						},
						Request: request,
					}
					processor := rc.CreateRequestProcessor()

					Expect(reflect.TypeOf(processor)).To(Equal(reflect.TypeOf(&push.PushRequestProcessor{})))
					concrete := processor.(*push.PushRequestProcessor)
					Expect(concrete.PushController).ToNot(BeNil())
					Expect(concrete.Response).To(Equal(response))
					Expect(concrete.Request).To(Equal(request))
					Expect(concrete.Log.UUID).To(Equal("the uuid"))
				})

			})
		})

		Describe("CreatePushController", func() {

			Context("when mock constructor is provided", func() {
				It("should return the mock implementation", func() {
					expected := &mocks.PushController{}
					creator := Creator{
						provider: CreatorModuleProvider{
							NewPushController: func(log I.DeploymentLogger, deployer, silentDeployer I.Deployer, eventManager I.EventManager, errorFinder I.ErrorFinder, pushManagerFactory I.PushManagerFactory, authResolver I.AuthResolver, resolver I.EnvResolver) request.PushController {
								return expected
							},
						},
					}
					rc := PushRequestCreator{
						RequestCreator: RequestCreator{
							Creator: creator,
						},
					}
					controller := rc.CreatePushController()
					Expect(controller).To(Equal(expected))
				})
			})

			Context("when mock constructor is not provided", func() {
				It("should return the default implementation", func() {
					creator := Creator{}
					rc := PushRequestCreator{
						RequestCreator: RequestCreator{
							Creator:      creator,
							Log:          I.DeploymentLogger{UUID: "the uuid"},
							EventManager: &mocks.EventManager{},
						},
					}
					controller := rc.CreatePushController()
					Expect(reflect.TypeOf(controller)).To(Equal(reflect.TypeOf(&push.PushController{})))
					concrete := controller.(*push.PushController)
					Expect(concrete.Deployer).ToNot(BeNil())
					Expect(concrete.SilentDeployer).ToNot(BeNil())
					Expect(concrete.Log.UUID).To(Equal("the uuid"))
					Expect(concrete.EventManager).To(Equal(rc.EventManager))
					Expect(concrete.ErrorFinder).ToNot(BeNil())
					Expect(concrete.PushManagerFactory).ToNot(BeNil())
					Expect(concrete.AuthResolver).ToNot(BeNil())
					Expect(concrete.EnvResolver).ToNot(BeNil())
				})
			})
		})

		Describe("PushManager", func() {

			Context("when mock constructor is provided", func() {
				It("should return the mock implementation", func() {
					expected := &mocks.PushManager{}
					creator := Creator{
						provider: CreatorModuleProvider{
							NewPushManager: func(courierCreator I.CourierCreator, eventManager I.EventManager, log I.DeploymentLogger, fetcher I.Fetcher, deployEventData structs.DeployEventData, fileSystemCleaner push.FileSystemCleaner, cfContext I.CFContext, auth I.Authorization, environment structs.Environment, envVars map[string]string, checker healthchecker.HealthChecker, mapper routemapper.RouteMapper) I.ActionCreator {
								return expected
							},
						},
					}
					rc := PushRequestCreator{
						RequestCreator: RequestCreator{
							Creator: creator,
						},
					}
					controller := rc.PushManager(structs.DeployEventData{}, I.Authorization{}, structs.Environment{}, make(map[string]string))
					Expect(controller).To(Equal(expected))
				})
			})

			Context("when mock constructor is not provided", func() {
				It("should return the default implementation", func() {
					creator := Creator{
						fileSystem: &afero.Afero{},
					}
					rc := PushRequestCreator{
						RequestCreator: RequestCreator{
							Creator:      creator,
							Log:          I.DeploymentLogger{UUID: "the uuid"},
							EventManager: &mocks.EventManager{},
						},
						Request: request.PostDeploymentRequest{
							Deployment: I.Deployment{
								CFContext: I.CFContext{
									Organization: "the org",
								},
							},
						},
					}
					controller := rc.PushManager(structs.DeployEventData{}, I.Authorization{}, structs.Environment{}, make(map[string]string))
					Expect(reflect.TypeOf(controller)).To(Equal(reflect.TypeOf(&push.PushManager{})))
					concrete := controller.(*push.PushManager)
					Expect(concrete.CourierCreator).To(Equal(creator))
					Expect(concrete.EventManager).To(Equal(rc.EventManager))
					Expect(concrete.Logger).To(Equal(rc.Log))
					Expect(concrete.Fetcher).ToNot(BeNil())
					Expect(concrete.DeployEventData).ToNot(BeNil())
					Expect(concrete.FileSystemCleaner).ToNot(BeNil())
					Expect(concrete.CFContext).To(Equal(rc.Request.Deployment.CFContext))
					Expect(concrete.Auth).ToNot(BeNil())
					Expect(concrete.Environment).ToNot(BeNil())
					Expect(concrete.EnvironmentVariables).ToNot(BeNil())
				})
			})
		})
	})

	Describe("StopRequestCreator", func() {

		Describe("CreateRequestProcessor", func() {
			Context("when mock constructor is provided", func() {
				It("should return the mock implementation", func() {

					expected := &mocks.RequestProcessor{}
					creator := Creator{
						provider: CreatorModuleProvider{
							NewStopRequestProcessor: func(log I.DeploymentLogger, sc request.StopController, request request.PutDeploymentRequest, buffer *bytes.Buffer) I.RequestProcessor {
								return expected
							},
						},
					}
					rc := StopRequestCreator{
						RequestCreator: RequestCreator{
							Creator: creator,
						},
					}
					processor := rc.CreateRequestProcessor()
					Expect(processor).To(Equal(expected))
				})
			})

			Context("when mock constructor is not provided", func() {
				It("should return the default implementation", func() {

					response := bytes.NewBuffer([]byte("the response"))
					request := request.PutDeploymentRequest{
						Deployment: I.Deployment{
							CFContext: I.CFContext{
								Organization: "the org",
							},
						},
					}

					rc := StopRequestCreator{
						RequestCreator: RequestCreator{
							Buffer: response,
							Log:    I.DeploymentLogger{UUID: "the uuid"},
						},
						Request: request,
					}
					processor := rc.CreateRequestProcessor()

					Expect(reflect.TypeOf(processor)).To(Equal(reflect.TypeOf(&stop.StopRequestProcessor{})))
					concrete := processor.(*stop.StopRequestProcessor)
					Expect(concrete.StopController).ToNot(BeNil())
					Expect(concrete.Response).To(Equal(response))
					Expect(concrete.Request).To(Equal(request))
					Expect(concrete.Log.UUID).To(Equal("the uuid"))
				})

			})
		})

		Describe("CreateStopController", func() {

			Context("when mock constructor is provided", func() {
				It("should return the mock implementation", func() {
					expected := &mocks.StopController{}
					creator := Creator{
						provider: CreatorModuleProvider{
							NewStopController: func(log I.DeploymentLogger, deployer I.Deployer, eventManager I.EventManager, errorFinder I.ErrorFinder, stopManagerFactory I.StopManagerFactory, authResolver I.AuthResolver, resolver I.EnvResolver) request.StopController {
								return expected
							},
						},
					}
					rc := StopRequestCreator{
						RequestCreator: RequestCreator{
							Creator: creator,
						},
					}
					controller := rc.CreateStopController()
					Expect(controller).To(Equal(expected))
				})
			})

			Context("when mock constructor is not provided", func() {
				It("should return the default implementation", func() {
					creator := Creator{}
					rc := StopRequestCreator{
						RequestCreator: RequestCreator{
							Creator:      creator,
							Log:          I.DeploymentLogger{UUID: "the uuid"},
							EventManager: &mocks.EventManager{},
						},
					}
					controller := rc.CreateStopController()
					Expect(reflect.TypeOf(controller)).To(Equal(reflect.TypeOf(&stop.StopController{})))
					concrete := controller.(*stop.StopController)
					Expect(concrete.Deployer).ToNot(BeNil())
					Expect(concrete.Log.UUID).To(Equal("the uuid"))
					Expect(concrete.EventManager).To(Equal(rc.EventManager))
					Expect(concrete.ErrorFinder).ToNot(BeNil())
					Expect(concrete.StopManagerFactory).ToNot(BeNil())
					Expect(concrete.AuthResolver).ToNot(BeNil())
					Expect(concrete.EnvResolver).ToNot(BeNil())
				})
			})
		})

		Describe("StopManager", func() {

			Context("when mock constructor is provided", func() {
				It("should return the mock implementation", func() {
					expected := &mocks.StopManager{}
					creator := Creator{
						provider: CreatorModuleProvider{
							NewStopManager: func(courierCreator I.CourierCreator, eventManager I.EventManager, log I.DeploymentLogger, deployEventData structs.DeployEventData) I.ActionCreator {
								return expected
							},
						},
					}
					rc := StopRequestCreator{
						RequestCreator: RequestCreator{
							Creator: creator,
						},
					}
					controller := rc.StopManager(structs.DeployEventData{})
					Expect(controller).To(Equal(expected))
				})
			})

			Context("when mock constructor is not provided", func() {
				It("should return the default implementation", func() {
					creator := Creator{}
					rc := StopRequestCreator{
						RequestCreator: RequestCreator{
							Creator:      creator,
							Log:          I.DeploymentLogger{UUID: "the uuid"},
							EventManager: &mocks.EventManager{},
						},
						Request: request.PutDeploymentRequest{
							Deployment: I.Deployment{
								CFContext: I.CFContext{
									Organization: "the org",
								},
							},
						},
					}
					controller := rc.StopManager(structs.DeployEventData{})
					Expect(reflect.TypeOf(controller)).To(Equal(reflect.TypeOf(&stop.StopManager{})))
					concrete := controller.(*stop.StopManager)
					Expect(concrete.CourierCreator).To(Equal(creator))
					Expect(concrete.EventManager).To(Equal(rc.EventManager))
					Expect(concrete.DeployEventData).ToNot(BeNil())
					Expect(concrete.Log).To(Equal(rc.Log))
				})
			})
		})

	})

	Describe("StartRequestCreator", func() {

		Describe("CreateRequestProcessor", func() {
			Context("when mock constructor is provided", func() {
				It("should return the mock implementation", func() {

					expected := &mocks.RequestProcessor{}
					creator := Creator{
						provider: CreatorModuleProvider{
							NewStartRequestProcessor: func(log I.DeploymentLogger, sc request.StartController, request request.PutDeploymentRequest, buffer *bytes.Buffer) I.RequestProcessor {
								return expected
							},
						},
					}
					rc := StartRequestCreator{
						RequestCreator: RequestCreator{
							Creator: creator,
						},
					}
					processor := rc.CreateRequestProcessor()
					Expect(processor).To(Equal(expected))
				})
			})

			Context("when mock constructor is not provided", func() {
				It("should return the default implementation", func() {

					response := bytes.NewBuffer([]byte("the response"))
					request := request.PutDeploymentRequest{
						Deployment: I.Deployment{
							CFContext: I.CFContext{
								Organization: "the org",
							},
						},
					}

					rc := StartRequestCreator{
						RequestCreator: RequestCreator{
							Buffer: response,
							Log:    I.DeploymentLogger{UUID: "the uuid"},
						},
						Request: request,
					}
					processor := rc.CreateRequestProcessor()

					Expect(reflect.TypeOf(processor)).To(Equal(reflect.TypeOf(&start.StartRequestProcessor{})))
					concrete := processor.(*start.StartRequestProcessor)
					Expect(concrete.StartController).ToNot(BeNil())
					Expect(concrete.Response).To(Equal(response))
					Expect(concrete.Request).To(Equal(request))
					Expect(concrete.Log.UUID).To(Equal("the uuid"))
				})

			})
		})

		Describe("CreateStartController", func() {

			Context("when mock constructor is provided", func() {
				It("should return the mock implementation", func() {
					expected := &mocks.StartController{}
					creator := Creator{
						provider: CreatorModuleProvider{
							NewStartController: func(log I.DeploymentLogger, deployer I.Deployer, eventManager I.EventManager, errorFinder I.ErrorFinder, startManagerFactory I.StartManagerFactory, authResolver I.AuthResolver, resolver I.EnvResolver) request.StartController {
								return expected
							},
						},
					}
					rc := StartRequestCreator{
						RequestCreator: RequestCreator{
							Creator: creator,
						},
					}
					controller := rc.CreateStartController()
					Expect(controller).To(Equal(expected))
				})
			})

			Context("when mock constructor is not provided", func() {
				It("should return the default implementation", func() {
					creator := Creator{}
					rc := StartRequestCreator{
						RequestCreator: RequestCreator{
							Creator:      creator,
							Log:          I.DeploymentLogger{UUID: "the uuid"},
							EventManager: &mocks.EventManager{},
						},
					}
					controller := rc.CreateStartController()
					Expect(reflect.TypeOf(controller)).To(Equal(reflect.TypeOf(&start.StartController{})))
					concrete := controller.(*start.StartController)
					Expect(concrete.Deployer).ToNot(BeNil())
					Expect(concrete.Log.UUID).To(Equal("the uuid"))
					Expect(concrete.EventManager).To(Equal(rc.EventManager))
					Expect(concrete.ErrorFinder).ToNot(BeNil())
					Expect(concrete.StartManagerFactory).ToNot(BeNil())
					Expect(concrete.AuthResolver).ToNot(BeNil())
					Expect(concrete.EnvResolver).ToNot(BeNil())
				})
			})
		})

		Describe("StartManager", func() {

			Context("when mock constructor is provided", func() {
				It("should return the mock implementation", func() {
					expected := &mocks.StartManager{}
					creator := Creator{
						provider: CreatorModuleProvider{
							NewStartManager: func(courierCreator I.CourierCreator, eventManager I.EventManager, log I.DeploymentLogger, deployEventData structs.DeployEventData) I.ActionCreator {
								return expected
							},
						},
					}
					rc := StartRequestCreator{
						RequestCreator: RequestCreator{
							Creator: creator,
						},
					}
					controller := rc.StartManager(structs.DeployEventData{})
					Expect(controller).To(Equal(expected))
				})
			})

			Context("when mock constructor is not provided", func() {
				It("should return the default implementation", func() {
					creator := Creator{}
					rc := StartRequestCreator{
						RequestCreator: RequestCreator{
							Creator:      creator,
							Log:          I.DeploymentLogger{UUID: "the uuid"},
							EventManager: &mocks.EventManager{},
						},
						Request: request.PutDeploymentRequest{
							Deployment: I.Deployment{
								CFContext: I.CFContext{
									Organization: "the org",
								},
							},
						},
					}
					controller := rc.StartManager(structs.DeployEventData{})
					Expect(reflect.TypeOf(controller)).To(Equal(reflect.TypeOf(&start.StartManager{})))
					concrete := controller.(*start.StartManager)
					Expect(concrete.CourierCreator).To(Equal(creator))
					Expect(concrete.EventManager).To(Equal(rc.EventManager))
					Expect(concrete.DeployEventData).ToNot(BeNil())
					Expect(concrete.Logger).To(Equal(rc.Log))
				})
			})
		})
	})
})
