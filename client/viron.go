// Code generated by goagen v1.3.1, DO NOT EDIT.
//
// API "Modoki API": viron Resource Client
//
// Command:
// $ goagen
// --design=github.com/cs3238-tsuzu/modoki/design
// --out=$(GOPATH)/src/github.com/cs3238-tsuzu/modoki
// --version=v1.3.1

package client

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// AuthtypeVironPath computes a request path to the authtype action of viron.
func AuthtypeVironPath() string {

	return fmt.Sprintf("/api/v1/viron_authtype")
}

// Get viron authtype
func (c *Client) AuthtypeViron(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewAuthtypeVironRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewAuthtypeVironRequest create the request corresponding to the authtype action endpoint of the viron resource.
func (c *Client) NewAuthtypeVironRequest(ctx context.Context, path string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "https"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// GetVironPath computes a request path to the get action of viron.
func GetVironPath() string {

	return fmt.Sprintf("/api/v1/viron")
}

// Get viron menu
func (c *Client) GetViron(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewGetVironRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewGetVironRequest create the request corresponding to the get action endpoint of the viron resource.
func (c *Client) NewGetVironRequest(ctx context.Context, path string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "https"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.JWTSigner != nil {
		if err := c.JWTSigner.Sign(req); err != nil {
			return nil, err
		}
	}
	return req, nil
}

// SigninVironPath computes a request path to the signin action of viron.
func SigninVironPath() string {

	return fmt.Sprintf("/api/v1/signin")
}

// Creates a valid JWT
func (c *Client) SigninViron(ctx context.Context, path string, payload *SigninPayload, contentType string) (*http.Response, error) {
	req, err := c.NewSigninVironRequest(ctx, path, payload, contentType)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewSigninVironRequest create the request corresponding to the signin action endpoint of the viron resource.
func (c *Client) NewSigninVironRequest(ctx context.Context, path string, payload *SigninPayload, contentType string) (*http.Request, error) {
	var body bytes.Buffer
	if contentType == "" {
		contentType = "*/*" // Use default encoder
	}
	err := c.Encoder.Encode(payload, &body, contentType)
	if err != nil {
		return nil, fmt.Errorf("failed to encode body: %s", err)
	}
	scheme := c.Scheme
	if scheme == "" {
		scheme = "https"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("POST", u.String(), &body)
	if err != nil {
		return nil, err
	}
	header := req.Header
	if contentType == "*/*" {
		header.Set("Content-Type", "application/json")
	} else {
		header.Set("Content-Type", contentType)
	}
	return req, nil
}
