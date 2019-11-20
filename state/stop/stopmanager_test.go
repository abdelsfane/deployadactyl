package stop_test

import (
	"github.com/compozed/deployadactyl/state/stop"
	"github.com/compozed/deployadactyl/structs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/compozed/deployadactyl/controller/deployer/bluegreen"
	"github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/mocks"
	"github.com/compozed/deployadactyl/randomizer"
	"github.com/go-errors/errors"
	"github.com/onsi/gomega/gbytes"
	"github.com/op/go-logging"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
)

type courierCreator struct {
	CourierCreatorFn func() (interfaces.Courier, error)
}

func (c courierCreator) CreateCourier() (interfaces.Courier, error) {
	if c.CourierCreatorFn != nil {
		return c.CourierCreatorFn()
	}

	courier := &mocks.Courier{}

	courier.LoginCall.Returns.Output = []byte("logged in\t")
	courier.DeleteCall.Returns.Output = []byte("deleted app\t")
	courier.PushCall.Returns.Output = []byte("pushed app\t")
	courier.RenameCall.Returns.Output = []byte("renamed app\t")
	courier.MapRouteCall.Returns.Output = append(courier.MapRouteCall.Returns.Output, []byte("mapped route\t"))
	courier.ExistsCall.Returns.Bool = true

	return courier, nil
}

var _ = Describe("Stopmanager", func() {
	var (
		response    io.ReadWriter
		stopManager interfaces.ActionCreator
		creator     *courierCreator
		logBuffer   *gbytes.Buffer
	)
	BeforeEach(func() {
		logBuffer = gbytes.NewBuffer()
		log := interfaces.DefaultLogger(logBuffer, logging.DEBUG, "deployer tests")
		response = gbytes.NewBuffer()
		creator = &courierCreator{}
		stopManager = stop.StopManager{
			CourierCreator: creator,
			Log:            interfaces.DeploymentLogger{log, randomizer.StringRunes(10)},
			DeployEventData: structs.DeployEventData{
				DeploymentInfo: &structs.DeploymentInfo{},
				Response:       response,
			},
		}
	})
	Describe("Create", func() {
		Context("when courier build succeeds", func() {
			It("should return a Stopper object", func() {
				env := structs.Environment{}
				foundationURL := "foundation url"
				stopper, _ := stopManager.Create(env, response, foundationURL)

				Expect(reflect.TypeOf(stopper)).Should(Equal(reflect.TypeOf(&stop.Stopper{})))

			})
			It("should return a Stopper object with correct data", func() {
				env := structs.Environment{
					Name: "myEnv",
				}
				foundationURL := "foundation url"
				deploymentInfo := structs.DeploymentInfo{
					AppName:  "myApp",
					Username: "bob",
					Password: "password",
				}
				*stopManager.(stop.StopManager).DeployEventData.DeploymentInfo = deploymentInfo
				stopper, _ := stopManager.Create(env, response, foundationURL)

				stopperData := stopper.(*stop.Stopper)
				Expect(stopperData.CFContext.Application).Should(Equal("myApp"))
				Expect(stopperData.CFContext.Environment).Should(Equal("myEnv"))
				Expect(stopperData.Authorization.Username).Should(Equal("bob"))
				Expect(stopperData.Authorization.Password).Should(Equal("password"))
				Expect(stopperData.FoundationURL).Should(Equal(foundationURL))

			})
		})

		Context("when courier build failed", func() {
			It("should return an error", func() {
				creator.CourierCreatorFn = func() (interfaces.Courier, error) {
					return nil, errors.New("a test error")
				}

				env := structs.Environment{}
				foundationURL := "foundation url"
				_, err := stopManager.Create(env, response, foundationURL)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).Should(ContainSubstring("a test error"))

			})
		})
	})
	Describe("OnFinish", func() {
		Context("when no error occurs", func() {
			It("returns http status OK", func() {
				deployResponse := stopManager.OnFinish(structs.Environment{}, response, nil)

				Expect(deployResponse.StatusCode).To(Equal(http.StatusOK))
			})
			It("logs successful stop", func() {
				stopManager.(stop.StopManager).DeployEventData.DeploymentInfo.AppName = "Conveyor"
				stopManager.OnFinish(structs.Environment{}, response, nil)

				Eventually(logBuffer).Should(gbytes.Say("successfully stopped application %s", "Conveyor"))
			})
			It("records success in the response", func() {
				stopManager.OnFinish(structs.Environment{}, response, nil)

				bytes, _ := ioutil.ReadAll(response)
				Eventually(string(bytes)).Should(ContainSubstring("Your stop was successful!"))
			})
		})

		Context("when an error occurs", func() {
			Context("and it is a log in error", func() {
				It("returns a http status bad request", func() {
					deployResponse := stopManager.OnFinish(structs.Environment{}, response, errors.New("login failed"))

					Expect(deployResponse.StatusCode).To(Equal(http.StatusBadRequest))
				})
			})
			It("returns a internal server error", func() {
				deployResponse := stopManager.OnFinish(structs.Environment{}, response, errors.New("a test error"))

				Expect(deployResponse.StatusCode).To(Equal(http.StatusInternalServerError))
			})
		})
	})
	Describe("InitiallyError", func() {
		It("should return LoginErrors", func() {
			errors := []error{errors.New("first error")}
			err := stopManager.InitiallyError(errors)

			Expect(reflect.TypeOf(err)).Should(Equal(reflect.TypeOf(bluegreen.LoginError{})))
		})
	})
	Describe("ExecuteError", func() {
		It("should return StopError", func() {
			errs := []error{errors.New("first error")}
			err := stopManager.ExecuteError(errs)

			Expect(reflect.TypeOf(err)).Should(Equal(reflect.TypeOf(bluegreen.StopError{})))
		})
	})
	Describe("UndoError", func() {
		It("should return RollbackStopError", func() {
			errs := []error{errors.New("first error")}
			executeErrors := []error{errors.New("execute error")}

			err := stopManager.UndoError(executeErrors, errs)

			Expect(reflect.TypeOf(err)).Should(Equal(reflect.TypeOf(bluegreen.RollbackStopError{})))
		})
	})
	Describe("SuccessError", func() {
		It("should return FinishStopError", func() {
			errors := []error{errors.New("first error")}
			err := stopManager.SuccessError(errors)

			Expect(reflect.TypeOf(err)).Should(Equal(reflect.TypeOf(bluegreen.FinishStopError{})))
		})
	})
})
