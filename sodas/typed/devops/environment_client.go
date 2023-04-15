package devops

import (
	"context"
	"path"

	devopsapi "github.com/xjayleex/sodas-sdk-go/api/devops"
	"github.com/xjayleex/sodas-sdk-go/rest"
)

type EnvironmentGetter interface {
	Environment() EnvironmentInterface
}

type EnvironmentInterface interface {
	RESTClient() rest.Interface

	List(context.Context, devopsapi.ListEnvironmentParams) (devopsapi.ListEnvironmentResult, error)
	ListAll(context.Context)
	Get(context.Context)
	Create(context.Context)
	Update(context.Context)
	Remove(context.Context)
	Execute(context.Context)

	TemplateGetter
}

type environment struct {
	path   string
	client rest.Interface
}

func newEnvironment(c DevelopInterface, prefix string) *environment {
	return &environment{
		path:   prefix + "/environment",
		client: c.RESTClient(),
	}
}

func (e *environment) RESTClient() rest.Interface {
	return e.client
}

func (e *environment) Template() TemplateInterface {
	return newTemplate(e, e.path)
}

func (e *environment) List(ctx context.Context, params devopsapi.ListEnvironmentParams) (result devopsapi.ListEnvironmentResult, err error) {
	endpoint := path.Join(e.path, "list")
	err = e.client.Get(endpoint).Do(ctx).Into(&result)
	return result, err
}
func (e *environment) ListAll(ctx context.Context) {}
func (e *environment) Get(ctx context.Context)     {}
func (e *environment) Create(ctx context.Context)  {}
func (e *environment) Update(ctx context.Context)  {}
func (e *environment) Remove(ctx context.Context)  {}
func (e *environment) Execute(ctx context.Context) {}

type TemplateGetter interface {
	Template() TemplateInterface
}

type TemplateInterface interface {
	ListPublic(context.Context)
	Get(context.Context)
	Create(context.Context)
	Update(context.Context)
	Remove(context.Context)
}

type template struct {
	path   string
	client rest.Interface
}

func newTemplate(e EnvironmentInterface, prefix string) *template {
	return &template{
		path:   prefix + "/template",
		client: e.RESTClient(),
	}
}

func (t *template) ListPublic(ctx context.Context) {}
func (t *template) Get(ctx context.Context)        {}
func (t *template) Create(ctx context.Context)     {}
func (t *template) Update(ctx context.Context)     {}
func (t *template) Remove(ctx context.Context)     {}
