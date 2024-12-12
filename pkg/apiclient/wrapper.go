package apiclient

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	Domain    string
	Schema    string
	Header    *http.Header
	Retry     int
	RetryTime int

	http *resty.Client
}

type Option func(*Client)

func New(options ...Option) *Client {
	c := &Client{
		http: resty.New(),
	}

	for _, opt := range options {
		opt(c)
	}

	return c
}

type Request struct {
	*resty.Request

	domain string
	schema string
}

type RequestOption func(*Request)

func (c *Client) NewRequest(options ...RequestOption) *Request {
	r := &Request{
		Request: c.http.R(),
		domain:  c.Domain,
		schema:  c.Schema,
	}

	if c.Retry > 0 {
		c.http.RetryCount = c.Retry
	}

	if c.RetryTime > 0 {
		c.http.RetryMaxWaitTime = time.Duration(c.RetryTime) * time.Millisecond
	}

	r.SetHeaderMultiValues(*c.Header)

	for _, opt := range options {
		opt(r)
	}

	return r
}

func (r *Request) Execute(method, url string) error {
	if r.domain != "" {
		url = strings.TrimRight(r.domain, "/") + "/" + strings.TrimLeft(url, "/")
	}

	resp, err := r.Request.Execute(method, r.schema+url)
	if err != nil {
		return err
	}

	if resp.IsError() {
		if resp.Error() == nil {
			return ApiError{
				Code: resp.StatusCode(),
				Body: resp.Body(),
			}
		}

		return resp.Error().(error)
	}

	return nil
}
