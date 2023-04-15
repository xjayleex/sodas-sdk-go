package gateway

import (
	"net/http"

	"github.com/xjayleex/sodas-sdk-go/rest"
)

type GatewayInterface interface {
	RESTClient() rest.Interface
	Auth() AuthInterface
}

func NewGatewayClient(c *rest.Config, hc *http.Client) (*GatewayClient, error) {
	configShallowCopy := *c
	client, err := rest.RESTClientForConfigAndClient(&configShallowCopy, hc)
	if err != nil {
		return nil, err
	}
	return &GatewayClient{
		path:   "/gateway",
		client: client,
	}, nil
}

type GatewayClient struct {
	path   string
	client *rest.RESTClient
}

func (c *GatewayClient) Auth() AuthInterface {
	return newAuth(c, c.path)
}

func (c *GatewayClient) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.client
}
