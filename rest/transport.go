package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"path"
	"time"

	"github.com/xjayleex/sodas-sdk-go/api/endpoints"
	"github.com/xjayleex/sodas-sdk-go/api/gateway"
)

func HTTPClientFor(config *Config) (*http.Client, error) {
	// we use http.DefaultTransport yet.
	transport, err := TransportFor(config)
	if err != nil {
		return nil, err
	}
	var httpClient *http.Client
	if transport != http.DefaultTransport || config.Timeout > 0 {
		httpClient = &http.Client{
			Transport: transport,
			Timeout:   config.Timeout,
		}
	} else {
		httpClient = http.DefaultClient
	}

	return httpClient, nil
}

func DefaultTransportFor() (http.RoundTripper, error) {
	return http.DefaultTransport, nil
}

func TransportFor(config *Config) (http.RoundTripper, error) {
	var rt http.RoundTripper
	rt = defaultTransport()
	if config.HasTokenAuth() {
		// FIXME hardcoded api url
		// TODO
		url, versionedAPIPath, err := defaultServerUrlFor(config)
		if err != nil {
			return nil, err
		}
		var p string
		//if !strings.HasSuffix(url.Path, "/") {
		//	url.Path += "/"
		//}
		p = url.String() + path.Join(versionedAPIPath, endpoints.API_V1_GATEWAY_AUTHENTICATION_USER_REFRESHUSER)

		rt = NewBearerAuthWithRefreshRoundTripper(p, config.ID, config.RefreshToken, rt)
	}

	return rt, nil
}

func defaultTransport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

func NewBearerAuthWithRefreshRoundTripper(url, id, refreshToken string, rt http.RoundTripper) http.RoundTripper {
	return &BearerAuthWithRefreshRoundTripper{
		URL:          url,
		ID:           id,
		RefreshToken: refreshToken,
		BearerToken:  "",
		rt:           rt,
		client:       http.DefaultClient,
	}
}

type BearerAuthWithRefreshRoundTripper struct {
	URL          string // Required
	ID           string // Required
	RefreshToken string // Required

	BearerToken string

	client *http.Client

	rt http.RoundTripper
}

func (rt *BearerAuthWithRefreshRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if rt.BearerToken == "" {
		if err := rt.refreshToken(); err != nil {
			return nil, err
		}
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", rt.BearerToken))
	resp, err := rt.rt.RoundTrip(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", rt.BearerToken))
		resp, err = rt.rt.RoundTrip(req)
		if err != nil {
			return nil, err
		}
	}
	return resp, nil
}

func (rt *BearerAuthWithRefreshRoundTripper) refreshToken() error {
	// use refresh token to get new access token from token URL
	// update t.AccessToken with the new access token
	return rt.roughRefreshLogic()
}

func (rt *BearerAuthWithRefreshRoundTripper) roughRefreshLogic() error {
	body := gateway.RefreshUserBody{
		Id:           rt.ID,
		RefreshToken: rt.RefreshToken,
	}

	marshalled, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", rt.URL, bytes.NewReader(marshalled))

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := rt.client.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return errors.New("could not refresh a token")
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	var tokenResponse gateway.TokenResponse

	if err = json.Unmarshal(respBody, &tokenResponse); err != nil {
		return err
	}

	rt.BearerToken = tokenResponse.AccessToken
	rt.RefreshToken = tokenResponse.RefreshToken

	return nil
}
