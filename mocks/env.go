package mocks

// Env handmade mock for tests.
type Env struct {
	GetCall struct {
		Received struct {
			Keys []string
		}
		Returns struct {
			Values map[string]string
		}
	}
}

// Get mock method.
func (e *Env) Get(key string) string {
	e.GetCall.Received.Keys = append(e.GetCall.Received.Keys, key)

	return e.GetCall.Returns.Values[key]
}
