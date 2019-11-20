package mocks

import (
	"fmt"
	"github.com/compozed/deployadactyl/interfaces"
)

// HealthChecker handmade mock for tests.
type HealthChecker struct {
	CheckCall struct {
		Received struct {
			Endpoint string
			URL      string
			Log      interfaces.DeploymentLogger
		}
		Returns struct {
			Error error
		}
	}
}

func (h *HealthChecker) Check(endpoint, serverURL string, log interfaces.DeploymentLogger) error {
	h.CheckCall.Received.Endpoint = endpoint
	h.CheckCall.Received.URL = fmt.Sprintf("%s/%s", serverURL, endpoint)
	h.CheckCall.Received.Log = log

	return h.CheckCall.Returns.Error
}
