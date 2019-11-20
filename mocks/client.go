package mocks

import "net/http"

// Client handmade mock for tests.
type Client struct {
	GetCall struct {
		Received struct {
			URL string
		}
		Returns struct {
			Response http.Response
			Error    error
		}
	}
}

// Get mock method.
func (c *Client) Get(url string) (*http.Response, error) {
	c.GetCall.Received.URL = url

	return &c.GetCall.Returns.Response, c.GetCall.Returns.Error
}
