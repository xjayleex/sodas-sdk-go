package datalake

import (
	"context"
	"errors"
	"fmt"
	"path"

	"github.com/xjayleex/sodas-sdk-go/api/datalake"
	"github.com/xjayleex/sodas-sdk-go/rest"
)

type CredentialGetter interface {
	Credential() CredentialInterface
}

type CredentialInterface interface {
	List(context.Context, datalake.ListCredentialParams) (datalake.ListCredentialResult, error)
	Create(context.Context) (datalake.CreateCredentialResult, error)
	Remove(context.Context) (datalake.RemoveCredentialResult, error)
}

type credential struct {
	path   string
	client rest.Interface
}

func newCredential(c *ObjectStorageClient) *credential {
	return &credential{
		path:   c.path + "/credential",
		client: c.client,
	}
}

func (c *credential) List(ctx context.Context, params datalake.ListCredentialParams) (result datalake.ListCredentialResult, err error) {
	endpoint := path.Join(c.path, "/list")
	err = c.client.Get(endpoint).
		Param("limit", fmt.Sprint(params.Limit)).
		Param("offset", fmt.Sprint(params.Offset)).
		Do(ctx).
		Into(&result)

	return result, err
}

func (c *credential) Create(ctx context.Context) (result datalake.CreateCredentialResult, err error) {

	return result, errors.New("unimplemented")
}

func (c *credential) Remove(ctx context.Context) (result datalake.RemoveCredentialResult, err error) {

	return result, errors.New("unimplemented")
}
