package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	ListCredentialEndpoint = "/api/v1/data-lake/object-storage/credential/list"
	// TODO : Where to get base url??
	HardCodedBaseUrl = "api.221.154.134.31.traefik.me:10011"
)

// minio keycloack credential이 계속해서 expire 된다면, 로직을 Token 받아오는 로직처럼 바꾸어야함..
type Credential struct {
	AccessKeyId     string
	SecretAccessKey string
}

type CredsListOpts struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

func GetCreds(tokenman TokenManager) ([]*Credential, error) {
	access, err := tokenman.AccessToken()
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("%s%s%s", HttpScheme, HardCodedBaseUrl, ListCredentialEndpoint)
	return requestCredentials(endpoint, *access, nil)
}

func requestCredentials(endpoint, accessToken string, listOpts *CredsListOpts) ([]*Credential, error) {
	if listOpts == nil {
		listOpts = &CredsListOpts{
			Offset: 0,
			Limit:  1,
		}
	}
	reqParam := listOpts
	marshalled, err := json.Marshal(reqParam)
	if err != nil {
		return nil, errors.New("unable to marshal request param")
	}

	req, err := http.NewRequest("GET", endpoint, bytes.NewReader(marshalled))
	if err != nil {
		return nil, fmt.Errorf("unable to make request with given request template : %v", err)
	}
	httpC := &http.Client{}
	req.Header.Set("Content-Type", "application/json")
	// Note that auth type bearer is hard coded here
	req.Header.Set("Authorization", fmt.Sprintf("%s%s", "Bearer ", accessToken))
	httpresp, err := httpC.Do(req)
	if err != nil {
		return nil, errors.New("unable to make http request do")
	}

	if httpresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("faield on get creds : status %d", httpresp.StatusCode)
	}

	defer httpresp.Body.Close()
	body, err := io.ReadAll(httpresp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body : %v", err)
	}

	//var data map[string]interface{}
	var data = struct {
		Results []struct {
			AccessKey string `json:"accessKey"`
			SecretKey string `json:"secretKey"`
		} `json:"results"`
		Total int `json:"total"`
	}{}

	if err = json.Unmarshal(body, &data); err != nil {
		return nil, errors.New("unable to unmarshal response body")
	}

	credentials := make([]*Credential, data.Total)
	for i, result := range data.Results {
		credentials[i] = &Credential{
			AccessKeyId:     result.AccessKey,
			SecretAccessKey: result.SecretKey,
		}
	}
	return credentials, nil
}
