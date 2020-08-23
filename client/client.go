package client

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/AlexKLWS/youtube-audio-stream/exerrors"
)

// ClientWrapper offers basic http request methods
type ClientWrapper struct {
	Silent     bool
	httpClient http.Client
}

type Client interface {
	HTTPGet(ctx context.Context, url string) (resp *http.Response, err error)
	HTTPGetBodyBytes(ctx context.Context, url string) ([]byte, error)
}

var instance Client

// NewSilent creates new HTTP client with silent set to true
func NewSilent(transport *http.Transport) Client {
	instance = &ClientWrapper{Silent: true, httpClient: http.Client{Transport: transport}}
	return instance
}

// New creates new HTTP client
func New(transport *http.Transport) Client {
	instance = &ClientWrapper{Silent: false, httpClient: http.Client{Transport: transport}}
	return instance
}

// Get returns current singleton instance of Client
func Get() Client {
	return instance
}

// HTTPGet does a HTTP GET request, checks the response to be a 200 OK and returns it
func (c *ClientWrapper) HTTPGet(ctx context.Context, url string) (resp *http.Response, err error) {
	if !c.Silent {
		log.Println("GET", url)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err = c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, exerrors.ErrUnexpectedStatusCode(resp.StatusCode)
	}

	return
}

// HTTPGetBodyBytes reads the whole HTTP body and returns it
func (c *ClientWrapper) HTTPGetBodyBytes(ctx context.Context, url string) ([]byte, error) {
	resp, err := c.HTTPGet(ctx, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
