package eventmanager

type InvalidArgumentError struct{}

func (e InvalidArgumentError) Error() string {
	return "invalid argument: error handler does not exist"
}
