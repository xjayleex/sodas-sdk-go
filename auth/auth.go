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

type tokenOpts struct {
	refreshEndpoint string
	accessEndpoint  string

	baseUrl  string
	uid      string
	password string
	offline  string
}

func NewTokenOpts(baseUrl, id, password, offline string) *tokenOpts {
	opts := &tokenOpts{
		baseUrl:  baseUrl,
		uid:      id,
		password: password,
		offline:  offline,
	}
	opts.refreshEndpoint = fmt.Sprintf("%s%s%s", HttpScheme, baseUrl, AuthLoginEndpoint)
	opts.accessEndpoint = fmt.Sprintf("%s%s%s", HttpScheme, baseUrl, AuthAccessTokenEndpoint)
	return opts
}

type TokenManager interface {
	AccessToken() (*string, error)
	Refresh() error
	UserId() string
}

// tokenManager must be initialized with NewTokenManager function.
type tokenManager struct {
	//	sync.Mutex
	httpC *http.Client
	a     *string // representing access token
	r     *string // representing refresh token
	opts  *tokenOpts
}

func NewTokenManager(opts *tokenOpts) *tokenManager {
	return &tokenManager{
		httpC: &http.Client{},
		opts:  opts,
	}
}

func (t *tokenManager) UserId() string {
	return t.opts.uid
}

func (t *tokenManager) AccessToken() (*string, error) {
	if t.a == nil && t.r == nil {
		if err := t.refresh(); err != nil {
			return nil, err
		} else {
			return t.a, nil
		}
	} else if t.a == nil && t.r != nil {
		err := t.Refresh()
		return t.a, err
	} else {
		return t.a, nil
	}
}

func (t *tokenManager) Refresh() error {
	if t.r != nil {
		na, err := t.getAccessToken()
		if err != nil {
			if errors.Is(err, ErrUnauthorized) {
				return t.refresh()
			} else {
				return errors.New("unknown error")
			}
		}
		t.a = &na
	}
	return t.refresh()
}

// let it can occur 401 error, ...
func (t *tokenManager) getAccessToken() (string, error) {
	reqBody := struct {
		Uid          string `json:"id"`
		RefreshToken string `json:"refreshToken"`
	}{
		Uid:          t.opts.uid,
		RefreshToken: *t.r,
	}

	endpoint := t.opts.accessEndpoint
	marshalled, err := json.Marshal(reqBody)
	if err != nil {
		return "", errors.New("unable to marshal request body")
	}
	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(marshalled))
	if err != nil {
		return "", fmt.Errorf("unable to make http request with given request template : %v", err)
	}
	// TODO : required http header must be added here.
	req.Header.Set("Content-Type", "application/json")
	resp, err := t.httpC.Do(req)
	if err != nil {
		return "", errors.New("unable to get access token")
	}

	if resp.StatusCode != http.StatusCreated {
		if resp.StatusCode == http.StatusUnauthorized {
			return "", ErrUnauthorized
		}
		return "", fmt.Errorf("failed on create access token with status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", fmt.Errorf("unable to read response body : %v", err)
	}

	var data map[string]interface{}

	if err = json.Unmarshal(body, &data); err != nil {
		return "", errors.New("unable to unmarshal response body")
	}

	if access, ok := data[AccessTokenKey]; ok {
		if refresh, ok := data[RefreshTokenKey]; ok {
			nr, na := fmt.Sprint(refresh), fmt.Sprint(access)
			t.r, t.a = &nr, &na
		}
	} else {
		return "", errors.New("either the key for a refresh or access token does not exist in response body")
	}
	return *t.a, nil
}

func (t *tokenManager) refresh() error {
	resp, err := t.requestRefresh()
	if err != nil {
		return err
	}
	t.a, t.r = &resp.access, &resp.refresh
	return nil
}

func (t *tokenManager) requestRefresh() (*tokenResponse, error) {
	reqBody := struct {
		Uid      string `json:"id"`
		Password string `json:"password"`
		Offline  string `json:"offline"`
	}{
		Uid:      t.opts.uid,
		Password: t.opts.password,
		Offline:  t.opts.offline,
	}
	// TODO : implement above bodyData

	endpoint := t.opts.refreshEndpoint
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

	var resp tokenResponse
	if refresh, ok := data[RefreshTokenKey]; ok {
		if access, ok := data[AccessTokenKey]; ok {
			resp.refresh, resp.access = fmt.Sprint(refresh), fmt.Sprint(access)
			t.r, t.a = &resp.refresh, &resp.access
		}
	} else {
		return nil, errors.New("either the key for a refresh or access token does not exist in response body")
	}
	return &resp, nil
}

type tokenResponse struct {
	access  string
	refresh string
}

/*
func (t RefreshToken) MarshalJSON() ([]byte, error) {
	//return []byte(fmt.Sprintf(`"refreshToken": "%s"`, t)), nil
	//return []byte(fmt.Sprintf(`{"refreshToken": "%s"}`, t)), nil
	return []byte(t), nil
}*/
