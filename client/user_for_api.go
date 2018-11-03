// Code generated by goagen v1.3.1, DO NOT EDIT.
//
// API "Modoki API": userForApi Resource Client
//
// Command:
// $ goagen
// --design=github.com/modoki-paas/modoki/design
// --out=$(GOPATH)/src/github.com/modoki-paas/modoki
// --version=v1.3.1

package client

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// AddAuthorizedKeysUserForAPIPath computes a request path to the addAuthorizedKeys action of userForApi.
func AddAuthorizedKeysUserForAPIPath() string {

	return fmt.Sprintf("/api/v2/user/config/authorizedKeys")
}

// AddAuthorizedKeysUserForAPI makes a request to the addAuthorizedKeys action endpoint of the userForApi resource
func (c *Client) AddAuthorizedKeysUserForAPI(ctx context.Context, path string, payload *UserAuthorizedKey, contentType string) (*http.Response, error) {
	req, err := c.NewAddAuthorizedKeysUserForAPIRequest(ctx, path, payload, contentType)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewAddAuthorizedKeysUserForAPIRequest create the request corresponding to the addAuthorizedKeys action endpoint of the userForApi resource.
func (c *Client) NewAddAuthorizedKeysUserForAPIRequest(ctx context.Context, path string, payload *UserAuthorizedKey, contentType string) (*http.Request, error) {
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
	req, err := http.NewRequest("PUT", u.String(), &body)
	if err != nil {
		return nil, err
	}
	header := req.Header
	if contentType == "*/*" {
		header.Set("Content-Type", "application/json")
	} else {
		header.Set("Content-Type", contentType)
	}
	if c.JWTSigner != nil {
		if err := c.JWTSigner.Sign(req); err != nil {
			return nil, err
		}
	}
	return req, nil
}

// GetConfigUserForAPIPath computes a request path to the getConfig action of userForApi.
func GetConfigUserForAPIPath() string {

	return fmt.Sprintf("/api/v2/user/config")
}

// GetConfigUserForAPI makes a request to the getConfig action endpoint of the userForApi resource
func (c *Client) GetConfigUserForAPI(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewGetConfigUserForAPIRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewGetConfigUserForAPIRequest create the request corresponding to the getConfig action endpoint of the userForApi resource.
func (c *Client) NewGetConfigUserForAPIRequest(ctx context.Context, path string) (*http.Request, error) {
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

// GetDefaultShellUserForAPIPath computes a request path to the getDefaultShell action of userForApi.
func GetDefaultShellUserForAPIPath() string {

	return fmt.Sprintf("/api/v2/user/config/defaultShell")
}

// GetDefaultShellUserForAPI makes a request to the getDefaultShell action endpoint of the userForApi resource
func (c *Client) GetDefaultShellUserForAPI(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewGetDefaultShellUserForAPIRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewGetDefaultShellUserForAPIRequest create the request corresponding to the getDefaultShell action endpoint of the userForApi resource.
func (c *Client) NewGetDefaultShellUserForAPIRequest(ctx context.Context, path string) (*http.Request, error) {
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

// ListAuthorizedKeysUserForAPIPath computes a request path to the listAuthorizedKeys action of userForApi.
func ListAuthorizedKeysUserForAPIPath() string {

	return fmt.Sprintf("/api/v2/user/config/authorizedKeys")
}

// ListAuthorizedKeysUserForAPI makes a request to the listAuthorizedKeys action endpoint of the userForApi resource
func (c *Client) ListAuthorizedKeysUserForAPI(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewListAuthorizedKeysUserForAPIRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewListAuthorizedKeysUserForAPIRequest create the request corresponding to the listAuthorizedKeys action endpoint of the userForApi resource.
func (c *Client) NewListAuthorizedKeysUserForAPIRequest(ctx context.Context, path string) (*http.Request, error) {
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

// RemoveAuthorizedKeysUserForAPIPath computes a request path to the removeAuthorizedKeys action of userForApi.
func RemoveAuthorizedKeysUserForAPIPath() string {

	return fmt.Sprintf("/api/v2/user/config/authorizedKeys")
}

// RemoveAuthorizedKeysUserForAPI makes a request to the removeAuthorizedKeys action endpoint of the userForApi resource
func (c *Client) RemoveAuthorizedKeysUserForAPI(ctx context.Context, path string, label string) (*http.Response, error) {
	req, err := c.NewRemoveAuthorizedKeysUserForAPIRequest(ctx, path, label)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewRemoveAuthorizedKeysUserForAPIRequest create the request corresponding to the removeAuthorizedKeys action endpoint of the userForApi resource.
func (c *Client) NewRemoveAuthorizedKeysUserForAPIRequest(ctx context.Context, path string, label string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "https"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	values := u.Query()
	values.Set("label", label)
	u.RawQuery = values.Encode()
	req, err := http.NewRequest("DELETE", u.String(), nil)
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

// SetAuthorizedKeysUserForAPIPayload is the userForApi setAuthorizedKeys action payload.
type SetAuthorizedKeysUserForAPIPayload []*UserAuthorizedKey

// SetAuthorizedKeysUserForAPIPath computes a request path to the setAuthorizedKeys action of userForApi.
func SetAuthorizedKeysUserForAPIPath() string {

	return fmt.Sprintf("/api/v2/user/config/authorizedKeys")
}

// SetAuthorizedKeysUserForAPI makes a request to the setAuthorizedKeys action endpoint of the userForApi resource
func (c *Client) SetAuthorizedKeysUserForAPI(ctx context.Context, path string, payload SetAuthorizedKeysUserForAPIPayload, contentType string) (*http.Response, error) {
	req, err := c.NewSetAuthorizedKeysUserForAPIRequest(ctx, path, payload, contentType)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewSetAuthorizedKeysUserForAPIRequest create the request corresponding to the setAuthorizedKeys action endpoint of the userForApi resource.
func (c *Client) NewSetAuthorizedKeysUserForAPIRequest(ctx context.Context, path string, payload SetAuthorizedKeysUserForAPIPayload, contentType string) (*http.Request, error) {
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
	if c.JWTSigner != nil {
		if err := c.JWTSigner.Sign(req); err != nil {
			return nil, err
		}
	}
	return req, nil
}

// SetDefaultShellUserForAPIPath computes a request path to the setDefaultShell action of userForApi.
func SetDefaultShellUserForAPIPath() string {

	return fmt.Sprintf("/api/v2/user/config/defaultShell")
}

// SetDefaultShellUserForAPI makes a request to the setDefaultShell action endpoint of the userForApi resource
func (c *Client) SetDefaultShellUserForAPI(ctx context.Context, path string, defaultShell string) (*http.Response, error) {
	req, err := c.NewSetDefaultShellUserForAPIRequest(ctx, path, defaultShell)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewSetDefaultShellUserForAPIRequest create the request corresponding to the setDefaultShell action endpoint of the userForApi resource.
func (c *Client) NewSetDefaultShellUserForAPIRequest(ctx context.Context, path string, defaultShell string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "https"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	values := u.Query()
	values.Set("defaultShell", defaultShell)
	u.RawQuery = values.Encode()
	req, err := http.NewRequest("POST", u.String(), nil)
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
