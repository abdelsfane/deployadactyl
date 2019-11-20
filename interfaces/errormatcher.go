package interfaces

type ErrorMatcher interface {
	Descriptor() string
	Match(matchTo []byte) LogMatchedError
}
