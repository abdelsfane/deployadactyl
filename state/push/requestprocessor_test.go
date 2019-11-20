package push

import (
	"bytes"

	"github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/mocks"
	"github.com/compozed/deployadactyl/request"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PushRequestProcessor", func() {

	Describe("Process", func() {
		It("calls RunDeployment with the Request", func() {
			pushController := &mocks.PushController{}

			processor := PushRequestProcessor{
				PushController: pushController,
				Request: request.PostDeploymentRequest{
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

			Eventually(pushController.RunDeploymentCall.Received.Request).Should(Equal(processor.Request))
		})

		It("calls RunDeployment with the Response", func() {
			pushController := &mocks.PushController{}

			processor := PushRequestProcessor{
				PushController: pushController,
				Response:       bytes.NewBuffer([]byte("foobar")),
			}

			processor.Process()

			Eventually(pushController.RunDeploymentCall.Received.Response).Should(Equal(processor.Response))
		})

	})
})
