package interfaces

// HealthChecker interface.
type HealthChecker interface {
	Check(endpoint, serverURL string, log DeploymentLogger) error
}
