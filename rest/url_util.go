package rest

import (
	"errors"
	"fmt"
	"net/url"
	"path"
)

func defaultServerUrlFor(config *Config) (*url.URL, string, error) {
	hasCA := len(config.CAFile) != 0 || len(config.CAData) != 0
	hasCert := len(config.CertFile) != 0 || len(config.CertData) != 0
	defaultTLS := hasCA || hasCert || config.Insecure
	host := config.Host
	if host == "" {
		return nil, "", errors.New("host field in config is a prerequisite field")
	}

	return DefaultServerURL(host, config.APIPath, config.APIVersion, defaultTLS)
}

func DefaultServerURL(host, apiPath, version string, defaultTLS bool) (*url.URL, string, error) {
	if host == "" {
		return nil, "", fmt.Errorf("host field in config is a preequisite field")
	}

	base := host
	hostURL, err := url.Parse(base)
	if err != nil || hostURL.Scheme == "" || hostURL.Host == "" {
		scheme := "http://"
		if defaultTLS {
			scheme = "https://"
		}
		hostURL, err = url.Parse(scheme + base)
		if err != nil {
			return nil, "", err
		}

		if hostURL.Path != "" && hostURL.Path != "/" {
			return nil, "", fmt.Errorf("host must be a URL or a host:port pair: %q", base)
		}
	}

	versionedAPIPath := path.Join("/", apiPath, version)
	return hostURL, versionedAPIPath, nil
}
