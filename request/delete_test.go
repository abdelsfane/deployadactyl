package request

import (
	"github.com/compozed/deployadactyl/interfaces"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"reflect"
)

var _ = Describe("DeleteDeploymentRequest", func() {
	Describe("GetId", func() {
		It("should return the UUID", func() {
			request := DeleteDeploymentRequest{
				Request: DeleteRequest{
					UUID: "abcd123",
				},
			}

			Expect(request.GetId()).To(Equal("abcd123"))
		})
	})

	Describe("GetData", func() {
		It("should return the data object", func() {
			request := DeleteDeploymentRequest{
				Request: DeleteRequest{
					Data: make(map[string]interface{}),
				},
			}
			request.Request.Data["foo"] = "bar"

			Expect(request.GetData()).To(Equal(request.Request.Data))
		})

		It("should allow modification of the data", func() {
			request := DeleteDeploymentRequest{
				Request: DeleteRequest{
					Data: make(map[string]interface{}),
				},
			}

			data := request.GetData()
			data["foo"] = "bar"

			Expect(request.GetData()["foo"]).To(Equal("bar"))
		})
	})

	Describe("GetContext", func() {
		It("should return the context", func() {
			request := DeleteDeploymentRequest{
				Deployment: interfaces.Deployment{
					CFContext: interfaces.CFContext{
						Environment: "the environment",
					},
				},
			}

			Expect(request.GetContext()).To(Equal(request.CFContext))
		})
	})

	Describe("SetContext", func() {
		It("should return a new modified DeleteDeploymentRequest", func() {
			origRequest := DeleteDeploymentRequest{
				Deployment: interfaces.Deployment{
					CFContext: interfaces.CFContext{
						Environment: "the environment",
					},
				},
			}

			expected := interfaces.CFContext{
				Application: "the application",
			}
			newRequest := origRequest.SetContext(expected)

			Expect(origRequest.GetContext().Environment).To(Equal("the environment"))
			Expect(newRequest.GetContext()).To(Equal(expected))
		})
	})

	Describe("GetAuthorization", func() {
		It("should return the authorization", func() {
			request := DeleteDeploymentRequest{
				Deployment: interfaces.Deployment{
					Authorization: interfaces.Authorization{
						Username: "the user",
					},
				},
			}

			Expect(request.GetAuthorization()).To(Equal(request.Deployment.Authorization))
		})
	})

	Describe("SetAuthorization", func() {
		It("should return a new modified DeleteDeploymentRequest", func() {
			origRequest := DeleteDeploymentRequest{
				Deployment: interfaces.Deployment{
					Authorization: interfaces.Authorization{
						Username: "username test",
					},
				},
			}

			expected := interfaces.Authorization{
				Password: "password test",
			}

			newRequest := origRequest.SetAuthorization(expected)

			Expect(origRequest.GetAuthorization().Username).To(Equal("username test"))
			Expect(newRequest.GetAuthorization()).To(Equal(expected))
		})
	})

	Describe("GetRequest", func() {
		It("should return the request object", func() {
			request := DeleteDeploymentRequest{
				Request: DeleteRequest{
					State: "test state data",
				},
			}

			Expect(request.GetRequest()).To(Equal(request.Request))
		})
	})

	Describe("SetRequest", func() {
		It("should return a new modified DeleteDeploymentRequest", func() {
			request := DeleteDeploymentRequest{
				Request: DeleteRequest{
					UUID: "UUID test data",
				},
			}

			expected := DeleteRequest{
				State: "this is the state test",
			}
			newRequest, _ := request.SetRequest(expected)

			Expect(request.GetRequest().(DeleteRequest).UUID).To(Equal("UUID test data"))
			Expect(newRequest.GetRequest()).To(Equal(expected))
		})

		Context("if provided object is not a DeleteRequest", func() {
			It("should return an error", func() {
				request := DeleteDeploymentRequest{
					Request: DeleteRequest{
						State: "test state data",
					},
				}

				_, err := request.SetRequest("")

				Expect(err).To(HaveOccurred())
				Expect(reflect.TypeOf(err)).To(Equal(reflect.TypeOf(InvalidArgumentError{})))
			})
		})
	})
})
