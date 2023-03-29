package clientgo_test

import (
	"os"
	"testing"

	clientgo "github.com/xjayleex/sodas-sdk-go/client-go"
	"github.com/xjayleex/sodas-sdk-go/testutil"
)

var (
	authArgs            = testutil.DataLakeAuthDefaultConfig
	refreshTokenRequest *clientgo.RefreshTokenRequest
	refreshToken        clientgo.RefreshToken
	accessToken         clientgo.AccessToken
)

func TestMain(m *testing.M) {
	refreshTokenRequest = clientgo.NewRefreshTokenRequest(authArgs["base_url"], authArgs["user_name"], authArgs["password"])
	ret := m.Run()
	os.Exit(ret)
}

func TestGetRefreshToken(t *testing.T) {
	var err error
	refreshToken, err = clientgo.GetRefreshToken(refreshTokenRequest)
	if err != nil {
		t.Log(err)
	} else {
		t.Log(refreshToken)
	}
}

func TestGetAccessToken(t *testing.T) {
	t.Log(refreshToken)
	var err error
	template := clientgo.NewAccessTokenRequest(authArgs["base_url"], authArgs["user_name"], refreshToken)
	accessToken, err = clientgo.GetAccessToken(template)
	if err != nil {
		t.Log(err)
	} else {
		t.Logf("AccessKey : %s", accessToken)
	}
}
