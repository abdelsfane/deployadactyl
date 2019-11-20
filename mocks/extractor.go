package mocks

// Extractor handmade mock for tests.
type Extractor struct {
	UnzipCall struct {
		Received struct {
			Source      string
			Destination string
			Manifest    string
		}
		Returns struct {
			Error error
		}
	}

	UntarCall struct {
		Received struct {
			Source      string
			Destination string
			Manifest    string
		}
		Returns struct {
			Error error
		}
	}
}

// Unzip mock method.
func (e *Extractor) Unzip(source, destination, manifest string) error {
	e.UnzipCall.Received.Source = source
	e.UnzipCall.Received.Destination = destination
	e.UnzipCall.Received.Manifest = manifest

	return e.UnzipCall.Returns.Error
}

func (e *Extractor) Untar(source, destination, manifest string) error {
	e.UntarCall.Received.Source = source
	e.UntarCall.Received.Destination = destination
	e.UntarCall.Received.Manifest = manifest

	return e.UntarCall.Returns.Error
}
