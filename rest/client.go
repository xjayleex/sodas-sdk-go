package rest

import (
	"net/http"
	"net/url"
	"strings"
)

type Interface interface {
	Verb(verb string, endpoint string) *Request
	Post(endpoint string) *Request
	Get(endpoint string) *Request
	// Sodas Plus APIs includes only Get, and Post method getter
	// Put() *Request
	// Delete() *Request
}

type RESTClient struct {
	base             *url.URL
	versionedAPIPath string
	refreshToken     string

	ContentType string
	Client      *http.Client
}

func NewRESTClient(baseURL *url.URL, versionedAPIPath, refreshToken string, client *http.Client) *RESTClient {
	base := *baseURL
	if !strings.HasSuffix(base.Path, "/") {
		base.Path += "/"
	}
	base.RawQuery = ""
	base.Fragment = ""
	return &RESTClient{
		base:             &base,
		versionedAPIPath: versionedAPIPath,
		refreshToken:     refreshToken,
		ContentType:      "application/json",
		Client:           client,
	}
}

func (c *RESTClient) Verb(verb string, endpoint string) *Request {
	return NewRequest(c, endpoint).Verb(verb)
}

func (c *RESTClient) Get(endpoint string) *Request {
	return c.Verb("GET", endpoint)
}

func (c *RESTClient) Post(endpoint string) *Request {
	return c.Verb("POST", endpoint)
}
