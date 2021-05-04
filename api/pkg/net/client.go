package net

import (
	"bsc-fees/pkg/config"
	"net/http"
	"net/url"
)

type Client interface {
	SendRequest(url *url.URL, headers map[string]string) (*http.Response, error)
}

type client struct {
	client http.Client
	cfg    config.HTTPClient
}

// NewClient creates a new Bsc/Eth compatible client
func NewClient(cfg config.Config) Client {
	c := http.Client{
		Timeout: cfg.HttpClient.Timeout,
	}

	return &client{
		client: c,
		cfg: cfg.HttpClient,
	}
}

func (c *client) SendRequest(url *url.URL, headers map[string]string) (*http.Response, error) {

	req, err := http.NewRequest(http.MethodGet, url.String(), nil)

	if err != nil {
		return &http.Response{}, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	req.Close = c.cfg.KeepAlive

	resp, err := c.client.Do(req)
	if err != nil {
		return &http.Response{}, err
	}

	return resp, err
}
