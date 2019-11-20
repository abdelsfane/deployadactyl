package interfaces

import (
	"io"
)

// Fetcher interface.
type Fetcher interface {
	Fetch(url, manifest string) (string, error)
	FetchArtifactFromRequest(body io.Reader, contentType string) (string, string, error)
}
