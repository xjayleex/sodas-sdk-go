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
	AuthLoginUrl    = "/api/v1/gateway/authentication/user/login"
	RefreshTokenKey = "refreshToken"
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
	url := fmt.Sprintf("%s%s%s", HttpScheme, template.baseURL, AuthLoginUrl)
	marshalled, err := json.Marshal(template)
	if err != nil {
		return "", errors.New("unable to marshal request body")
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(marshalled))
	if err != nil {
		return "", fmt.Errorf("unable to make http request with given request template : %v", err)
	}
	// TODO : required http header must be added here.
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New("unable to get refresh token")
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

	token, ok := data["refreshToken"]

	if !ok {
		return "", errors.New("failed on get refreshToken from response body")
	}
	return RefreshToken(fmt.Sprint(token)), nil

}
