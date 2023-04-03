package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	ErrUnauthorized = fmt.Errorf("Unauthorized(401)")
)

const (
	HttpScheme = "http://"
)

const (
	AuthLoginEndpoint       = "/api/v1/gateway/authentication/user/login"
	AuthAccessTokenEndpoint = "/api/v1/gateway/authentication/user/refreshUser"
	RefreshTokenKey         = "refreshToken"
	AccessTokenKey          = "accessToken"
)

type AccessTokenOpts struct {
	// refreshEndpoint string
	url string

	BaseUrl string
	Uid     string
	Refresh string
}

func (o *AccessTokenOpts) apiUrl() string {
	if o.url == "" {
		o.url = fmt.Sprintf("%s%s%s", HttpScheme, o.BaseUrl, AuthAccessTokenEndpoint)
		return o.url
	}
	return o.url
}

type TokenManager interface {
	Token() (string, error)
	UserId() string
}

// tokenManager must be initialized with NewTokenManager function.
type tokenManager struct {
	//	sync.Mutex
	httpC http.Client
	a     string // representing access token
	r     string // representing refresh token
	opts  AccessTokenOpts
}

func NewTokenManager(opts AccessTokenOpts) *tokenManager {
	return &tokenManager{
		httpC: http.Client{},
		opts:  opts,
		a:     "",
		r:     opts.Refresh,
	}
}

func (t *tokenManager) UserId() string {
	return t.opts.Uid
}

func (t *tokenManager) Token() (string, error) {
	if t.a == "" {
		// try get access token with refresh token =-> AuthAccessTokenEndpoint
		resp, err := t.requestAccessToken()
		if err != nil /* && errors.Is(err, ErrUnauthorized) */ {
			// we can try to get new refresh token only if we have user password.
			return "", err
		}
		t.a, t.r = resp.access, resp.refresh
	}
	return t.a, nil
}

// let it can occur 401 error, ...
func (t *tokenManager) requestAccessToken() (*tokenResponse, error) {
	reqBody := struct {
		Uid          string `json:"id"`
		RefreshToken string `json:"refreshToken"`
	}{
		Uid:          t.opts.Uid,
		RefreshToken: t.r,
	}

	endpoint := t.opts.apiUrl()
	marshalled, err := json.Marshal(reqBody)
	if err != nil {
		return nil, errors.New("unable to marshal request body")
	}
	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(marshalled))
	if err != nil {
		return nil, fmt.Errorf("unable to make http request with given request template : %v", err)
	}
	// TODO : required http header must be added here.
	req.Header.Set("Content-Type", "application/json")
	httpresp, err := t.httpC.Do(req)
	if err != nil {
		return nil, errors.New("unable to get access token")
	}

	if httpresp.StatusCode != http.StatusCreated {
		if httpresp.StatusCode == http.StatusUnauthorized {
			return nil, ErrUnauthorized
		}
		return nil, fmt.Errorf("failed on create access token with status code: %d", httpresp.StatusCode)
	}

	defer httpresp.Body.Close()
	body, err := io.ReadAll(httpresp.Body)

	if err != nil {
		return nil, fmt.Errorf("unable to read response body : %v", err)
	}

	var data map[string]interface{}

	if err = json.Unmarshal(body, &data); err != nil {
		return nil, errors.New("unable to unmarshal response body")
	}
	resp := &tokenResponse{}
	if access, ok := data[AccessTokenKey]; ok {
		if refresh, ok := data[RefreshTokenKey]; ok {
			resp.access, resp.refresh = fmt.Sprint(access), fmt.Sprint(refresh)
		} else {
			resp.access, resp.refresh = fmt.Sprint(access), t.r
		}
	} else {
		return nil, errors.New("either the key for a refresh or access token does not exist in response body")
	}
	return resp, nil
}

type tokenResponse struct {
	access  string
	refresh string
}

func (t *tokenResponse) Access() string {
	return t.access
}

func (t *tokenResponse) Refresh() string {
	return t.refresh
}

type RefreshTokenOpts struct {
	url string

	BaseUrl  string
	Uid      string
	Password string
	Offline  string
}

func (o *RefreshTokenOpts) apiUrl() string {
	if o.url == "" {
		o.url = fmt.Sprintf("%s%s%s", HttpScheme, o.BaseUrl, AuthLoginEndpoint)
		return o.url
	}
	return o.url
}

func RefreshToken(opts RefreshTokenOpts) (*tokenResponse, error) {
	httpC := http.Client{}
	reqBody := struct {
		Uid      string `json:"id"`
		Password string `json:"password"`
		Offline  string `json:"offline"`
	}{
		Uid:      opts.Uid,
		Password: opts.Password,
		Offline:  opts.Offline,
	}

	endpoint := opts.apiUrl()
	marshalled, err := json.Marshal(reqBody)
	if err != nil {
		return nil, errors.New("unable to marshal request body")
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(marshalled))
	if err != nil {
		return nil, fmt.Errorf("unable to make http request with given request template : %v", err)
	}
	// TODO : required http header must be added here.
	req.Header.Set("Content-Type", "application/json")

	httpresp, err := httpC.Do(req)
	if err != nil {
		return nil, errors.New("unable to get refresh token")
	}

	if httpresp.StatusCode != http.StatusCreated {
		return nil, errors.New("failed on create refresh token")
	}

	defer httpresp.Body.Close()
	body, err := io.ReadAll(httpresp.Body)

	if err != nil {
		return nil, fmt.Errorf("unable to read response body : %v", err)
	}

	var data map[string]interface{}

	if err = json.Unmarshal(body, &data); err != nil {
		return nil, errors.New("unable to unmarshal response body")
	}

	resp := &tokenResponse{}
	if refresh, ok := data[RefreshTokenKey]; !ok {
		return nil, errors.New("either the key for a refresh token does not exist in response body")
	} else {
		resp.refresh = fmt.Sprint(refresh)
	}

	if access, ok := data[AccessTokenKey]; ok {
		resp.access = fmt.Sprint(access)
	}

	return resp, nil
}
