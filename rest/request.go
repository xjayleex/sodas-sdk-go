package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"

	"golang.org/x/net/http2"
)

type Request struct {
	c   *RESTClient
	err error

	basicAuth bool

	pathPrefix string
	verb       string
	endpoint   string

	timeout time.Duration

	params  url.Values
	headers http.Header
	// body    io.Reader
	bodyBytes []byte
}

func NewRequest(c *RESTClient, endpoint string) *Request {
	var pathPrefix string
	if c.base != nil {
		pathPrefix = path.Join("/", c.base.Path, c.versionedAPIPath)
	} else {
		pathPrefix = path.Join("/", c.versionedAPIPath)
	}

	var timeout time.Duration
	if c.Client != nil {
		timeout = c.Client.Timeout
	}
	return &Request{
		c:          c,
		pathPrefix: pathPrefix,
		endpoint:   endpoint,
		timeout:    timeout,
		basicAuth:  true,
	}
}

func (r *Request) Verb(verb string) *Request {
	r.verb = verb
	return r
}

func (r *Request) Do(ctx context.Context) Result {
	var result Result
	// result = r.request_depr(ctx)
	err := r.request(ctx, func(req *http.Request, resp *http.Response) {
		result = r.transformResponse(req, resp)
	})
	if err != nil {
		return Result{
			err: err,
		}
	}
	return result
}

func (r *Request) request(ctx context.Context, fn func(*http.Request, *http.Response)) error {
	var client *http.Client
	client = r.c.Client
	if client == nil {
		client = http.DefaultClient
	}

	if !r.basicAuth {
		// FIXME transport cache에서 bearer auth roundtripper 빠진 transport 사용하도로 구현해야함.
		// TODO
	}

	if r.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, r.timeout)
		defer cancel()
	}

	req, err := r.newHTTPRequest(ctx)
	if err != nil {
		return err
	}
	// convert or not transport layer
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer readAndCloseResponseBody(resp)
	fn(req, resp)
	return nil
}

func (r *Request) transformResponse(req *http.Request, resp *http.Response) Result {
	var body []byte
	if resp.Body != nil {
		data, err := io.ReadAll(resp.Body)
		switch err.(type) {
		case nil:
			body = data
		case http2.StreamError:
			return Result{
				err: fmt.Errorf("stream error %#v when reading response body, may be caused ny close connection", err),
			}
		default:
			return Result{
				err: fmt.Errorf("unexpected error when reading response body, Original error: %w", err),
			}
		}
	}

	contentType := resp.Header.Get("Content-Type")
	if len(contentType) == 0 {
		contentType = r.c.ContentType
	}

	switch {
	case resp.StatusCode == http.StatusSwitchingProtocols:
		// no-op
	case resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusPartialContent:
		var errorResponse ErrorResponse
		err := json.Unmarshal(body, &errorResponse)
		if err != nil {
			return Result{
				err: errors.New("status error, with unknown message"),
			}
		}
		return Result{
			err: fmt.Errorf("status error, with message : %s", errorResponse.Message),
		}
	}

	return Result{
		statusCode:  resp.StatusCode,
		contentType: contentType,
		body:        body,
		err:         nil,
	}
}

func (r *Request) newHTTPRequest(ctx context.Context) (*http.Request, error) {
	body := bytes.NewReader(r.bodyBytes)

	url := r.URL().String()
	req, err := http.NewRequest(r.verb, url, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Header = r.headers
	return req, nil
}

func (r *Request) URL() *url.URL {
	p := r.pathPrefix
	p = path.Join(p, r.endpoint)
	// TODO more
	finalURL := &url.URL{}
	if r.c.base != nil {
		*finalURL = *r.c.base
	}

	finalURL.Path = p

	query := url.Values{}
	for key, values := range r.params {
		for _, value := range values {
			query.Add(key, value)
		}
	}
	finalURL.RawQuery = query.Encode()
	return finalURL
}

func (r *Request) SetHeader(key string, values ...string) *Request {
	if r.headers == nil {
		r.headers = http.Header{}
	}

	r.headers.Del(key)

	for _, value := range values {
		r.headers.Add(key, value)
	}
	return r
}

func (r *Request) Param(name, value string) *Request {
	if r.err != nil {
		return r
	}
	return r.setParam(name, value)
}

func (r *Request) setParam(name, value string) *Request {
	if r.params == nil {
		r.params = make(url.Values)
	}
	r.params[name] = append(r.params[name], value)
	return r
}

// 만약에 Body()가 호출되지 않으면(Get call) Content-Type은? Content Type은 Get요청에 필요없나 ?
func (r *Request) Body(obj interface{}) *Request {
	if r.err != nil {
		return r
	}
	data, err := json.Marshal(obj)
	if err != nil {
		r.err = err
		return r
	}
	r.bodyBytes = data
	r.SetHeader("Content-Type", r.c.ContentType)
	return r
}

func (r *Request) DisableAuth() *Request {
	r.basicAuth = false
	return r
}

type Result struct {
	statusCode  int
	contentType string
	body        []byte
	err         error
}

func (r Result) Into(obj interface{}) error {
	if r.err != nil {
		return r.err
	}
	err := json.Unmarshal(r.body, obj)
	return err
}

func readAndCloseResponseBody(resp *http.Response) {
	if resp == nil {
		return
	}

	const maxBodySlurpSize = 2 << 10
	defer resp.Body.Close()

	if resp.ContentLength <= maxBodySlurpSize {
		io.Copy(io.Discard, &io.LimitedReader{
			R: resp.Body,
			N: maxBodySlurpSize,
		})
	}
}

type ErrorResponse struct {
	StatusCode int    `json:"statusCode"`
	Code       int    `json:"code"`
	Flag       string `json:"flag"`
	Message    string `json:"message"`
}
