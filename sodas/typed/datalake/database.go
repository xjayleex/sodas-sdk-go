package datalake

import "github.com/xjayleex/sodas-sdk-go/rest"

type DatabaseGetter interface {
	Database() DatabaseInterface
}

type DatabaseInterface interface {
}

type database struct {
	path   string
	client rest.Interface
}

func newDatabase(c *PostgresqlClient) *database {
	return &database{
		path:   c.path + "/database",
		client: c.client,
	}
}
