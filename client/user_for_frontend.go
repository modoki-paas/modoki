// Code generated by goagen v1.3.1, DO NOT EDIT.
//
// API "Modoki API": userForFrontend Resource Client
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

// AddAuthorizedKeysUserForFrontendPath computes a request path to the addAuthorizedKeys action of userForFrontend.
func AddAuthorizedKeysUserForFrontendPath() string {

	return fmt.Sprintf("/frontend/v2/user/config/authorizedKeys")
}

// AddAuthorizedKeysUserForFrontend makes a request to the addAuthorizedKeys action endpoint of the userForFrontend resource
func (c *Client) AddAuthorizedKeysUserForFrontend(ctx context.Context, path string, payload *UserAuthorizedKey, contentType string) (*http.Response, error) {
	req, err := c.NewAddAuthorizedKeysUserForFrontendRequest(ctx, path, payload, contentType)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewAddAuthorizedKeysUserForFrontendRequest create the request corresponding to the addAuthorizedKeys action endpoint of the userForFrontend resource.
func (c *Client) NewAddAuthorizedKeysUserForFrontendRequest(ctx context.Context, path string, payload *UserAuthorizedKey, contentType string) (*http.Request, error) {
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

// GetAPIKeyUserForFrontendPath computes a request path to the getAPIKey action of userForFrontend.
func GetAPIKeyUserForFrontendPath() string {

	return fmt.Sprintf("/frontend/v2/user/apiKey")
}

// GetAPIKeyUserForFrontend makes a request to the getAPIKey action endpoint of the userForFrontend resource
func (c *Client) GetAPIKeyUserForFrontend(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewGetAPIKeyUserForFrontendRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewGetAPIKeyUserForFrontendRequest create the request corresponding to the getAPIKey action endpoint of the userForFrontend resource.
func (c *Client) NewGetAPIKeyUserForFrontendRequest(ctx context.Context, path string) (*http.Request, error) {
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

// GetConfigUserForFrontendPath computes a request path to the getConfig action of userForFrontend.
func GetConfigUserForFrontendPath() string {

	return fmt.Sprintf("/frontend/v2/user/config")
}

// GetConfigUserForFrontend makes a request to the getConfig action endpoint of the userForFrontend resource
func (c *Client) GetConfigUserForFrontend(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewGetConfigUserForFrontendRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewGetConfigUserForFrontendRequest create the request corresponding to the getConfig action endpoint of the userForFrontend resource.
func (c *Client) NewGetConfigUserForFrontendRequest(ctx context.Context, path string) (*http.Request, error) {
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

// GetDefaultShellUserForFrontendPath computes a request path to the getDefaultShell action of userForFrontend.
func GetDefaultShellUserForFrontendPath() string {

	return fmt.Sprintf("/frontend/v2/user/config/defaultShell")
}

// GetDefaultShellUserForFrontend makes a request to the getDefaultShell action endpoint of the userForFrontend resource
func (c *Client) GetDefaultShellUserForFrontend(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewGetDefaultShellUserForFrontendRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewGetDefaultShellUserForFrontendRequest create the request corresponding to the getDefaultShell action endpoint of the userForFrontend resource.
func (c *Client) NewGetDefaultShellUserForFrontendRequest(ctx context.Context, path string) (*http.Request, error) {
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

// ListAuthorizedKeysUserForFrontendPath computes a request path to the listAuthorizedKeys action of userForFrontend.
func ListAuthorizedKeysUserForFrontendPath() string {

	return fmt.Sprintf("/frontend/v2/user/config/authorizedKeys")
}

// ListAuthorizedKeysUserForFrontend makes a request to the listAuthorizedKeys action endpoint of the userForFrontend resource
func (c *Client) ListAuthorizedKeysUserForFrontend(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewListAuthorizedKeysUserForFrontendRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewListAuthorizedKeysUserForFrontendRequest create the request corresponding to the listAuthorizedKeys action endpoint of the userForFrontend resource.
func (c *Client) NewListAuthorizedKeysUserForFrontendRequest(ctx context.Context, path string) (*http.Request, error) {
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

// ReissueAPIKeyUserForFrontendPath computes a request path to the reissueAPIKey action of userForFrontend.
func ReissueAPIKeyUserForFrontendPath() string {

	return fmt.Sprintf("/frontend/v2/user/apiKey")
}

// ReissueAPIKeyUserForFrontend makes a request to the reissueAPIKey action endpoint of the userForFrontend resource
func (c *Client) ReissueAPIKeyUserForFrontend(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewReissueAPIKeyUserForFrontendRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewReissueAPIKeyUserForFrontendRequest create the request corresponding to the reissueAPIKey action endpoint of the userForFrontend resource.
func (c *Client) NewReissueAPIKeyUserForFrontendRequest(ctx context.Context, path string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "https"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
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

// RemoveAuthorizedKeysUserForFrontendPath computes a request path to the removeAuthorizedKeys action of userForFrontend.
func RemoveAuthorizedKeysUserForFrontendPath() string {

	return fmt.Sprintf("/frontend/v2/user/config/authorizedKeys")
}

// RemoveAuthorizedKeysUserForFrontend makes a request to the removeAuthorizedKeys action endpoint of the userForFrontend resource
func (c *Client) RemoveAuthorizedKeysUserForFrontend(ctx context.Context, path string, label string) (*http.Response, error) {
	req, err := c.NewRemoveAuthorizedKeysUserForFrontendRequest(ctx, path, label)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewRemoveAuthorizedKeysUserForFrontendRequest create the request corresponding to the removeAuthorizedKeys action endpoint of the userForFrontend resource.
func (c *Client) NewRemoveAuthorizedKeysUserForFrontendRequest(ctx context.Context, path string, label string) (*http.Request, error) {
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

// SetAuthorizedKeysUserForFrontendPayload is the userForFrontend setAuthorizedKeys action payload.
type SetAuthorizedKeysUserForFrontendPayload []*UserAuthorizedKey

// SetAuthorizedKeysUserForFrontendPath computes a request path to the setAuthorizedKeys action of userForFrontend.
func SetAuthorizedKeysUserForFrontendPath() string {

	return fmt.Sprintf("/frontend/v2/user/config/authorizedKeys")
}

// SetAuthorizedKeysUserForFrontend makes a request to the setAuthorizedKeys action endpoint of the userForFrontend resource
func (c *Client) SetAuthorizedKeysUserForFrontend(ctx context.Context, path string, payload SetAuthorizedKeysUserForFrontendPayload, contentType string) (*http.Response, error) {
	req, err := c.NewSetAuthorizedKeysUserForFrontendRequest(ctx, path, payload, contentType)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewSetAuthorizedKeysUserForFrontendRequest create the request corresponding to the setAuthorizedKeys action endpoint of the userForFrontend resource.
func (c *Client) NewSetAuthorizedKeysUserForFrontendRequest(ctx context.Context, path string, payload SetAuthorizedKeysUserForFrontendPayload, contentType string) (*http.Request, error) {
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

// SetDefaultShellUserForFrontendPath computes a request path to the setDefaultShell action of userForFrontend.
func SetDefaultShellUserForFrontendPath() string {

	return fmt.Sprintf("/frontend/v2/user/config/defaultShell")
}

// SetDefaultShellUserForFrontend makes a request to the setDefaultShell action endpoint of the userForFrontend resource
func (c *Client) SetDefaultShellUserForFrontend(ctx context.Context, path string, defaultShell string) (*http.Response, error) {
	req, err := c.NewSetDefaultShellUserForFrontendRequest(ctx, path, defaultShell)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewSetDefaultShellUserForFrontendRequest create the request corresponding to the setDefaultShell action endpoint of the userForFrontend resource.
func (c *Client) NewSetDefaultShellUserForFrontendRequest(ctx context.Context, path string, defaultShell string) (*http.Request, error) {
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