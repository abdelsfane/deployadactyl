package mocks

import (
	"io"
)

// Fetcher handmade mock for tests.
type Fetcher struct {
	FetchCall struct {
		Received struct {
			ArtifactURL string
			Manifest    string
		}
		Returns struct {
			AppPath string
			Error   error
		}
	}

	FetchArtifactFromRequestCall struct {
		Received struct {
			Request     io.Reader
			ContentType string
		}
		Returns struct {
			AppPath  string
			Manifest string
			Error    error
		}
	}
}

// Fetch mock method.
func (f *Fetcher) Fetch(url, manifest string) (string, error) {
	f.FetchCall.Received.ArtifactURL = url
	f.FetchCall.Received.Manifest = manifest

	return f.FetchCall.Returns.AppPath, f.FetchCall.Returns.Error
}

// FetchZipFromRequest mock method.
func (f *Fetcher) FetchArtifactFromRequest(body io.Reader, contentType string) (string, string, error) {
	f.FetchArtifactFromRequestCall.Received.Request = body
	f.FetchArtifactFromRequestCall.Received.ContentType = contentType

	return f.FetchArtifactFromRequestCall.Returns.AppPath, f.FetchArtifactFromRequestCall.Returns.Manifest, f.FetchArtifactFromRequestCall.Returns.Error
}
