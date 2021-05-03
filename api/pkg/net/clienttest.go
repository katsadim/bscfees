package net

import (
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/url"
)

type MockClient struct {
	mock.Mock
}

func (c *MockClient) SendRequest(url *url.URL, h map[string]string) (*http.Response, error) {
	args := c.Called(url, h)
	firstArg := args.Get(0)

	var resp *http.Response
	switch firstArg := firstArg.(type) {
	case nil:
		resp = nil
	case *http.Response:
		resp = firstArg
	}

	return resp, args.Error(1)
}
