/*
Modify the http package to allow easy calls
*/
package client

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type HTTP struct {
	client *http.Client
}

var (
	defaultTimeout = 30 * time.Second
	activeClient   = &http.Client{Timeout: defaultTimeout}
	clientMu       sync.RWMutex
)

func buildHTTPClient(timeout time.Duration, proxyURL string) (*http.Client, error) {
	if timeout <= 0 {
		timeout = defaultTimeout
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()

	if proxyURL != "" {
		proxy, err := url.Parse(proxyURL)
		if err != nil {
			return nil, err
		}
		transport.Proxy = http.ProxyURL(proxy)
	}

	return &http.Client{
		Timeout:   timeout,
		Transport: transport,
	}, nil
}

func Configure(proxyURL string) error {
	client, err := buildHTTPClient(defaultTimeout, proxyURL)
	if err != nil {
		return err
	}

	clientMu.Lock()
	activeClient = client
	clientMu.Unlock()

	return nil
}

func currentClient() *http.Client {
	clientMu.RLock()
	defer clientMu.RUnlock()
	return activeClient
}

// timeout the http request to a level and sets a proxy is passed
func New(timeout time.Duration, proxy string) (*HTTP, error) {
	client, err := buildHTTPClient(timeout, proxy)
	if err != nil {
		return nil, err
	}

	return &HTTP{client: client}, nil
}

// make a GET request and return some info about the request
func Get(url string) (int, []byte, error) {
	resp, err := GetResponse(url)

	if err != nil {
		return -1, nil, err
	}

	// read the response body
	bodyBytes, err := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()

	if err != nil {
		return resp.StatusCode, nil, err
	}

	return resp.StatusCode, bodyBytes, nil
}

func GetResponse(url string) (*http.Response, error) {
	return currentClient().Get(url)
}
