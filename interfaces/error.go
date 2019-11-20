package interfaces

type DeploymentError interface {
	Code() string
	Error() string
}

type LogMatchedError interface {
	Code() string
	Error() string
	Details() []string
	Solution() string
}
