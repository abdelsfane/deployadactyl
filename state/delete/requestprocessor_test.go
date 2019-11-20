package delete

import (
	"bytes"

	"github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/mocks"

	"github.com/compozed/deployadactyl/request"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeleteRequestProcessor", func() {

	Describe("Process", func() {
		It("calls DeleteDeployment with the Request", func() {
			deleteController := &mocks.DeleteController{}

			processor := DeleteRequestProcessor{
				DeleteController: deleteController,
				Request: request.DeleteDeploymentRequest{
					Deployment: interfaces.Deployment{
						CFContext: interfaces.CFContext{
							Environment:  "the environment",
							Space:        "the space",
							Organization: "the org",
							Application:  "the app",
						},
						Authorization: interfaces.Authorization{
							Username: "the user",
							Password: "the password",
						},
					},
				},
			}

			processor.Process()

			Eventually(deleteController.DeleteDeploymentCall.Received.Deployment).Should(Equal(processor.Request))
		})

		It("calls DeleteDeployment with the Response", func() {
			deleteController := &mocks.DeleteController{}

			processor := DeleteRequestProcessor{
				DeleteController: deleteController,
				Response:         bytes.NewBuffer([]byte("foobar")),
			}

			processor.Process()

			Eventually(deleteController.DeleteDeploymentCall.Received.Response).Should(Equal(processor.Response))
		})

	})
})
