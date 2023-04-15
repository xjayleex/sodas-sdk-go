package devops

import "github.com/xjayleex/sodas-sdk-go/rest"

type DevelopInterface interface {
	EnvironmentGetter
	RESTClient() rest.Interface
}

type DevelopClient struct {
	path   string
	client rest.Interface
}

func newDevelop(c DevopsInterface, prefix string) DevelopInterface {
	return &DevelopClient{
		path:   prefix + "/development",
		client: c.RESTClient(),
	}
}

func (c *DevelopClient) Environment() EnvironmentInterface {
	return newEnvironment(c, c.path)
}

func (c *DevelopClient) RESTClient() rest.Interface {
	return c.client
}
