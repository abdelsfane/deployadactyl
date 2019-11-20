package request

import (
	"github.com/compozed/deployadactyl/interfaces"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"reflect"
)

var _ = Describe("PostDeploymentRequest", func() {
	Describe("GetId", func() {
		It("should return the UUID", func() {
			request := PostDeploymentRequest{
				Request: PostRequest{
					UUID: "abcd123",
				},
			}

			Expect(request.GetId()).To(Equal("abcd123"))
		})
	})

	Describe("GetData", func() {
		It("should return the data object", func() {
			request := PostDeploymentRequest{
				Request: PostRequest{
					Data: make(map[string]interface{}),
				},
			}
			request.Request.Data["foo"] = "bar"

			Expect(request.GetData()).To(Equal(request.Request.Data))
		})

		It("should allow modification of the data", func() {
			request := PostDeploymentRequest{
				Request: PostRequest{
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
			request := PostDeploymentRequest{
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
		It("should return a new modified PostDeploymentRequest", func() {
			origRequest := PostDeploymentRequest{
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
			request := PostDeploymentRequest{
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
		It("should return a new modified PostDeploymentRequest", func() {
			origRequest := PostDeploymentRequest{
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
			request := PostDeploymentRequest{
				Request: PostRequest{
					ArtifactUrl: "testArtifactUrl.com/.com",
				},
			}

			Expect(request.GetRequest()).To(Equal(request.Request))
		})
	})

	Describe("SetRequest", func() {
		It("should return a new modified PostDeploymentRequest", func() {
			request := PostDeploymentRequest{
				Request: PostRequest{
					ArtifactUrl: "testArtifactUrl.com/.com",
				},
			}

			expected := PostRequest{
				Manifest: "this is the manifest",
			}
			newRequest, _ := request.SetRequest(expected)

			Expect(request.GetRequest().(PostRequest).ArtifactUrl).To(Equal("testArtifactUrl.com/.com"))
			Expect(newRequest.GetRequest()).To(Equal(expected))
		})

		Context("if provided object is not a PostRequest", func() {
			It("should return an error", func() {
				request := PostDeploymentRequest{
					Request: PostRequest{
						ArtifactUrl: "testArtifactUrl.com/.com",
					},
				}

				_, err := request.SetRequest("")

				Expect(err).To(HaveOccurred())
				Expect(reflect.TypeOf(err)).To(Equal(reflect.TypeOf(InvalidArgumentError{})))
			})
		})
	})
})
