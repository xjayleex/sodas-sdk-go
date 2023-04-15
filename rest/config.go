package rest

import (
	"net/http"
	"time"
)

type Config struct {
	Host         string
	APIPath      string
	APIVersion   string
	Username     string
	Password     string
	ID           string // used in BearerTokenRoundTripper
	RefreshToken string // used in BearerTokenRoundTripper
	TLSClientConfig
}

func (c *Config) HasTokenAuth() bool {
	return len(c.RefreshToken) != 0 && len(c.ID) != 0
}

type TLSClientConfig struct {
	Insecure bool

	ServerName string

	CertFile string
	KeyFile  string
	CAFile   string

	CertData []byte
	KeyData  []byte
	CAData   []byte

	Timeout time.Duration
}

func RESTClientFor(config *Config) (*RESTClient, error) {
	// Validate config.Host
	_, _, err := defaultServerUrlFor(config)
	if err != nil {
		return nil, err
	}

	httpClient, err := HTTPClientFor(config)
	if err != nil {
		return nil, err
	}
	return RESTClientForConfigAndClient(config, httpClient)

}

func RESTClientForConfigAndClient(config *Config, httpClient *http.Client) (*RESTClient, error) {
	baseURL, versionedAPIPath, err := defaultServerUrlFor(config)
	if err != nil {
		return nil, err
	}
	restClient := NewRESTClient(baseURL, versionedAPIPath, config.RefreshToken, httpClient)
	return restClient, nil
}
