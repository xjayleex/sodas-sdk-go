package clientgo_test

import (
	"log"
	"testing"

	clientgo "github.com/xjayleex/sodas-sdk-go/client-go"
	"github.com/xjayleex/sodas-sdk-go/testutil"
)

var (
	authArgs            = testutil.DataLakeAuthDefaultArgs
	refreshTokenRequest *clientgo.RefreshTokenRequest
)

func TestMain(m *testing.M) {
	refreshTokenRequest = clientgo.NewRefreshTokenRequest()
}

func TestGetRefreshToken(t *testing.T) {
	_, err := clientgo.GetRefreshToken(refreshTokenRequest)
	if err != nil {
		log.Println(err)
	}

}
