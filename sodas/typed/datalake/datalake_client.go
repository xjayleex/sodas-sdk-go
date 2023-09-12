package datalake

import (
	"net/http"

	"github.com/xjayleex/sodas-sdk-go/rest"
)

type DatalakeInterface interface {
	RESTClient() rest.Interface
	ObjectStorage() ObjectStorageInterface
	Postgresql() PostgresqlInterface
}

func NewDatalakeClient(c *rest.Config, hc *http.Client) (*DatalakeClient, error) {
	configShallowCopy := *c
	client, err := rest.RESTClientForConfigAndClient(&configShallowCopy, hc)
	if err != nil {
		return nil, err
	}

	return &DatalakeClient{
		path:   "/data-lake",
		client: client,
	}, nil
}

type DatalakeClient struct {
	path   string
	client *rest.RESTClient
}

func (c *DatalakeClient) ObjectStorage() ObjectStorageInterface {
	return newObjectStorage(c, c.path)
}

func (c *DatalakeClient) Postgresql() PostgresqlInterface {
	return newPostgresql(c, c.path)
}

func (c *DatalakeClient) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.client
}
