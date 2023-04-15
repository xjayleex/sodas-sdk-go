package devops

import (
	"net/http"

	"github.com/xjayleex/sodas-sdk-go/rest"
)

type DevopsInterface interface {
	RESTClient() rest.Interface
	Develop() DevelopInterface
}

func NewDevopsClient(c *rest.Config, hc *http.Client) (*DevopsClient, error) {
	configShallowCopy := *c
	client, err := rest.RESTClientForConfigAndClient(&configShallowCopy, hc)
	if err != nil {
		return nil, err
	}

	return &DevopsClient{
		path:   "/devops",
		client: client,
	}, nil
}

type DevopsClient struct {
	path   string
	client *rest.RESTClient
}

func (c *DevopsClient) Develop() DevelopInterface {
	return newDevelop(c, c.path)
}

func (c *DevopsClient) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.client
}
