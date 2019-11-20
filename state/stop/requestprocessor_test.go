package stop

import (
	"bytes"

	"github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/mocks"

	"github.com/compozed/deployadactyl/request"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StopRequestProcessor", func() {

	Describe("Process", func() {
		It("calls StopDeployment with the Request", func() {
			stopController := &mocks.StopController{}

			processor := StopRequestProcessor{
				StopController: stopController,
				Request: request.PutDeploymentRequest{
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

			Eventually(stopController.StopDeploymentCall.Received.Deployment).Should(Equal(processor.Request))
		})

		It("calls StopDeployment with the Response", func() {
			stopController := &mocks.StopController{}

			processor := StopRequestProcessor{
				StopController: stopController,
				Response:       bytes.NewBuffer([]byte("foobar")),
			}

			processor.Process()

			Eventually(stopController.StopDeploymentCall.Received.Response).Should(Equal(processor.Response))
		})

	})
})
