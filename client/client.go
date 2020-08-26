package client

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/AlexKLWS/youtube-audio-stream/consts"
	"github.com/AlexKLWS/youtube-audio-stream/exerrors"
	"github.com/spf13/viper"
	"golang.org/x/net/proxy"
)

// ClientWrapper offers basic http request methods
type ClientWrapper struct {
	httpClient http.Client
}

type Client interface {
	HTTPGet(ctx context.Context, url string) (resp *http.Response, err error)
	HTTPGetBodyBytes(ctx context.Context, url string) ([]byte, error)
}

var instance Client

func GetHTTPTransport() *http.Transport {
	httpTransport := &http.Transport{
		IdleConnTimeout:       60 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	socks5Proxy := viper.GetString(consts.SocksProxy)

	if socks5Proxy != "" {
		if viper.GetBool(consts.Debug) {
			log.Println("Using SOCKS5 proxy", socks5Proxy)
		}
		dialer, err := proxy.SOCKS5("tcp", socks5Proxy, nil, proxy.Direct)
		if err != nil {
			fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
			os.Exit(1)
		}

		// set our socks5 as the dialer
		dc := dialer.(interface {
			DialContext(ctx context.Context, network, addr string) (net.Conn, error)
		})
		httpTransport.DialContext = dc.DialContext
	} else {
		httpTransport.DialContext = (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext
	}

	return httpTransport
}

// New creates new HTTP client
func New(transport *http.Transport) Client {
	instance = &ClientWrapper{httpClient: http.Client{Transport: transport}}
	return instance
}

// Get returns current singleton instance of Client
func Get() Client {
	return instance
}

// HTTPGet does a HTTP GET request, checks the response to be a 200 OK and returns it
func (c *ClientWrapper) HTTPGet(ctx context.Context, url string) (resp *http.Response, err error) {
	if viper.GetBool(consts.Debug) {
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
