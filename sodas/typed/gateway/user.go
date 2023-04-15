package gateway

import (
	"context"
	"path"

	gatewayapi "github.com/xjayleex/sodas-sdk-go/api/gateway"
	"github.com/xjayleex/sodas-sdk-go/rest"
)

type UserGetter interface {
	User() UserInterface
}

type UserInterface interface {
	Login(context.Context, gatewayapi.LoginBody) (gatewayapi.LoginResult, error)
	RefreshUser(context.Context, gatewayapi.RefreshUserBody) (gatewayapi.RefreshUserResult, error)
}

type user struct {
	path   string
	client rest.Interface
}

func newUser(c *AuthClient) *user {
	return &user{
		path:   c.path + "/user",
		client: c.client,
	}
}

func (u *user) Login(ctx context.Context, body gatewayapi.LoginBody) (result gatewayapi.LoginResult, err error) {
	endpoint := path.Join(u.path, "/login")
	err = u.client.Post(endpoint).
		DisableAuth().
		Body(body).
		Do(ctx).
		Into(&result)
	return result, err
}

func (u *user) RefreshUser(ctx context.Context, body gatewayapi.RefreshUserBody) (result gatewayapi.RefreshUserResult, err error) {
	endpoint := path.Join(u.path, "/refreshUser")
	err = u.client.Post(endpoint).
		DisableAuth().
		Body(body).
		Do(ctx).
		Into(&result)
	return result, err
}
