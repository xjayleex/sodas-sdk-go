package datalake

import "github.com/xjayleex/sodas-sdk-go/rest"

type ObjectStorageInterface interface {
	CredentialGetter
}

type ObjectStorageClient struct {
	path   string
	client rest.Interface
}

func newObjectStorage(c DatalakeInterface, prefix string) ObjectStorageInterface {
	return &ObjectStorageClient{
		path:   prefix + "/object-storage",
		client: c.RESTClient(),
	}
}

func (c *ObjectStorageClient) Credential() CredentialInterface {
	return newCredential(c)
}
