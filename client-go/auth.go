package clientgo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	AuthLoginEndpoint       = "/api/v1/gateway/authentication/user/login"
	AuthAccessTokenEndpoint = "/api/v1/gateway/authentication/user/refreshUser"
	RefreshTokenKey         = "refreshToken"
	AccessTokenKey          = "accessToken"
)

type RefreshToken string

type RefreshTokenRequest struct {
	baseURL  string
	Uid      string `json:"id"`
	Password string `json:"password"`
	Offline  string `json:"offline"`
}

func NewRefreshTokenRequest(baseURL, uid, password string) *RefreshTokenRequest {
	return &RefreshTokenRequest{
		baseURL:  baseURL,
		Uid:      uid,
		Password: password,
		Offline:  "false",
	}
}

func GetRefreshToken(template *RefreshTokenRequest) (RefreshToken, error) {
	endpoint := fmt.Sprintf("%s%s%s", HttpScheme, template.baseURL, AuthLoginEndpoint)
	marshalled, err := json.Marshal(template)
	if err != nil {
		return "", errors.New("unable to marshal request body")
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(marshalled))
	if err != nil {
		return "", fmt.Errorf("unable to make http request with given request template : %v", err)
	}
	// TODO : required http header must be added here.
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New("unable to get refresh token")
	}

	if resp.StatusCode != http.StatusCreated {
		return "", errors.New("failed on create refresh token")
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

	token, ok := data[RefreshTokenKey]

	if !ok {
		return "", errors.New("key for refresh token does not exist in response body")
	}
	return RefreshToken(fmt.Sprint(token)), nil
}

type AccessToken string

type AccessTokenRequest struct {
	baseURL      string
	Uid          string       `json:"id"`
	RefreshToken RefreshToken `json:"refreshToken"`
}

func NewAccessTokenRequest(baseURL, uid string, refreshToken *RefreshToken) *AccessTokenRequest {
	return &AccessTokenRequest{
		baseURL:      baseURL,
		Uid:          uid,
		RefreshToken: RefreshTokenKey,
	}
}

func GetAccessToken(template *AccessTokenRequest) (AccessToken, error) {
	endpoint := fmt.Sprintf("%s%s%s", HttpScheme, template.baseURL, AuthAccessTokenEndpoint)
	marshalled, err := json.Marshal(template)
	if err != nil {
		return "", errors.New("unable to marshal request body")
	}
	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(marshalled))
	if err != nil {
		return "", fmt.Errorf("unable to make http request with given request template : %v", err)
	}

	// TODO : required http header must be added here.
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New("unable to get access token")
	}

	if resp.StatusCode != http.StatusCreated {
		return "", errors.New("failed on create access token")
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

	token, ok := data[AccessTokenKey]
	if !ok {
		return "", errors.New("key for access token does not exist in response body")
	}
	return AccessToken(fmt.Sprint(token)), nil
}
