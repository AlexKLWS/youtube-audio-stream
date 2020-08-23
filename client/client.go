package client

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/AlexKLWS/youtube-audio-stream/exerrors"
)

// Client offers basic http request methods
type Client struct {
	Silent     bool
	httpClient http.Client
}

// NewSilent returns new HTTP client with silent set to true
func NewSilent(transport *http.Transport) *Client {
	return &Client{Silent: true, httpClient: http.Client{Transport: transport}}
}

// New returns new HTTP client
func New(transport *http.Transport) *Client {
	return &Client{Silent: false, httpClient: http.Client{Transport: transport}}
}

// HTTPGet does a HTTP GET request, checks the response to be a 200 OK and returns it
func (c *Client) HTTPGet(ctx context.Context, url string) (resp *http.Response, err error) {
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
func (c *Client) HTTPGetBodyBytes(ctx context.Context, url string) ([]byte, error) {
	resp, err := c.HTTPGet(ctx, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
