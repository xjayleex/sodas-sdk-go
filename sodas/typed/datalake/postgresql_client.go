package datalake

import "github.com/xjayleex/sodas-sdk-go/rest"

type PostgresqlInterface interface {
	DatabaseGetter
}

type PostgresqlClient struct {
	path   string
	client rest.Interface
}

func newPostgresql(c DatalakeInterface, prefix string) PostgresqlInterface {
	return &PostgresqlClient{
		path:   prefix + "/postgresql",
		client: c.RESTClient(),
	}
}

func (c *PostgresqlClient) Database() DatabaseInterface {
	return newDatabase(c)
}
