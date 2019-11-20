package interfaces

import "net/http"

// Client is an interface for http.Client.
type Client interface {
	Get(url string) (*http.Response, error)
}
