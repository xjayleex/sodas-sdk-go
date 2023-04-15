package sodas

import (
	"net/http"

	"github.com/xjayleex/sodas-sdk-go/rest"
	"github.com/xjayleex/sodas-sdk-go/sodas/typed/datalake"
	"github.com/xjayleex/sodas-sdk-go/sodas/typed/devops"
	gateway "github.com/xjayleex/sodas-sdk-go/sodas/typed/gateway"
)

type Interface interface {
	Gateway() gateway.GatewayInterface
	Datalake() datalake.DatalakeInterface
	Devops() devops.DevopsInterface
}

type Clientset struct {
	gateway  *gateway.GatewayClient
	datalake *datalake.DatalakeClient
	devops   *devops.DevopsClient
}

func (c *Clientset) Gateway() gateway.GatewayInterface {
	return c.gateway
}

func (c *Clientset) Datalake() datalake.DatalakeInterface {
	return c.datalake
}

func (c *Clientset) Devops() devops.DevopsInterface {
	return c.devops
}

func NewForConfig(c *rest.Config) (*Clientset, error) {
	configShallowCopy := *c

	httpClient, err := rest.HTTPClientFor(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	return NewForConfigAndClient(&configShallowCopy, httpClient)
}

func NewForConfigAndClient(c *rest.Config, httpClient *http.Client) (*Clientset, error) {
	configShallowCopy := *c
	var cs Clientset
	var err error
	cs.gateway, err = gateway.NewGatewayClient(&configShallowCopy, httpClient)
	if err != nil {
		return nil, err
	}

	cs.datalake, err = datalake.NewDatalakeClient(&configShallowCopy, httpClient)
	if err != nil {
		return nil, err
	}

	cs.devops, err = devops.NewDevopsClient(&configShallowCopy, httpClient)
	if err != nil {
		return nil, err
	}

	return &cs, nil
}
