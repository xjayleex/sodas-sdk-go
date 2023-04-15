package gateway

import "github.com/xjayleex/sodas-sdk-go/rest"

type AuthInterface interface {
	UserGetter
}

type AuthClient struct {
	path   string
	client rest.Interface
}

func newAuth(c GatewayInterface, prefix string) AuthInterface {
	return &AuthClient{
		path:   prefix + "/authentication",
		client: c.RESTClient(),
	}
}
func (c *AuthClient) User() UserInterface {
	return newUser(c)
}
