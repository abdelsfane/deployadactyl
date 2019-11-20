package start

import (
	"bytes"

	"github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/mocks"

	"github.com/compozed/deployadactyl/request"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StartRequestProcessor", func() {

	Describe("Process", func() {
		It("calls StartDeployment with the Request", func() {
			startController := &mocks.StartController{}

			processor := StartRequestProcessor{
				StartController: startController,
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

			Eventually(startController.StartDeploymentCall.Received.Deployment).Should(Equal(processor.Request))
		})

		It("calls StopDeployment with the Response", func() {
			startController := &mocks.StartController{}

			processor := StartRequestProcessor{
				StartController: startController,
				Response:        bytes.NewBuffer([]byte("foobar")),
			}

			processor.Process()

			Eventually(startController.StartDeploymentCall.Received.Response).Should(Equal(processor.Response))
		})

	})
})
