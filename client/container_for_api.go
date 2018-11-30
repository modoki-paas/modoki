// Code generated by goagen v1.3.1, DO NOT EDIT.
//
// API "Modoki API": containerForApi Resource Client
//
// Command:
// $ goagen
// --design=github.com/modoki-paas/modoki/design
// --out=$(GOPATH)/src/github.com/modoki-paas/modoki
// --version=v1.4.0

package client

import (
	"bytes"
	"context"
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// CreateContainerForAPIPath computes a request path to the create action of containerForApi.
func CreateContainerForAPIPath() string {

	return fmt.Sprintf("/api/v2/container/create")
}

// create a new container
func (c *Client) CreateContainerForAPI(ctx context.Context, path string, image string, name string, command []string, entrypoint []string, env []string, sslRedirect *bool, volumes []string, workingDir *string) (*http.Response, error) {
	req, err := c.NewCreateContainerForAPIRequest(ctx, path, image, name, command, entrypoint, env, sslRedirect, volumes, workingDir)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewCreateContainerForAPIRequest create the request corresponding to the create action endpoint of the containerForApi resource.
func (c *Client) NewCreateContainerForAPIRequest(ctx context.Context, path string, image string, name string, command []string, entrypoint []string, env []string, sslRedirect *bool, volumes []string, workingDir *string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "https"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	values := u.Query()
	values.Set("image", image)
	values.Set("name", name)
	for _, p := range command {
		tmp59 := p
		values.Add("command", tmp59)
	}
	for _, p := range entrypoint {
		tmp60 := p
		values.Add("entrypoint", tmp60)
	}
	for _, p := range env {
		tmp61 := p
		values.Add("env", tmp61)
	}
	if sslRedirect != nil {
		tmp62 := strconv.FormatBool(*sslRedirect)
		values.Set("sslRedirect", tmp62)
	}
	for _, p := range volumes {
		tmp63 := p
		values.Add("volumes", tmp63)
	}
	if workingDir != nil {
		values.Set("workingDir", *workingDir)
	}
	u.RawQuery = values.Encode()
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.APIKeySigner != nil {
		if err := c.APIKeySigner.Sign(req); err != nil {
			return nil, err
		}
	}
	return req, nil
}

// DownloadContainerForAPIPath computes a request path to the download action of containerForApi.
func DownloadContainerForAPIPath(id string) string {
	param0 := id

	return fmt.Sprintf("/api/v2/container/%s/download", param0)
}

// DownloadContainerForAPIPath2 computes a request path to the download action of containerForApi.
func DownloadContainerForAPIPath2() string {

	return fmt.Sprintf("/api/v2/container/download")
}

// Copy files from the container
func (c *Client) DownloadContainerForAPI(ctx context.Context, path string, internalPath string) (*http.Response, error) {
	req, err := c.NewDownloadContainerForAPIRequest(ctx, path, internalPath)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewDownloadContainerForAPIRequest create the request corresponding to the download action endpoint of the containerForApi resource.
func (c *Client) NewDownloadContainerForAPIRequest(ctx context.Context, path string, internalPath string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "https"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	values := u.Query()
	values.Set("internalPath", internalPath)
	u.RawQuery = values.Encode()
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.APIKeySigner != nil {
		if err := c.APIKeySigner.Sign(req); err != nil {
			return nil, err
		}
	}
	return req, nil
}

// ExecContainerForAPIPath computes a request path to the exec action of containerForApi.
func ExecContainerForAPIPath(id string) string {
	param0 := id

	return fmt.Sprintf("/api/v2/container/%s/exec", param0)
}

// Exec a command with attaching to a container using WebSocket(Mainly for xterm.js, using a protocol for terminado)
func (c *Client) ExecContainerForAPI(ctx context.Context, path string, command []string, tty *bool) (*websocket.Conn, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "ws"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	values := u.Query()
	if command != nil {
		for _, p := range command {
			tmp64 := p
			values.Add("command", tmp64)
		}
	}
	if tty != nil {
		tmp65 := strconv.FormatBool(*tty)
		values.Set("tty", tmp65)
	}
	u.RawQuery = values.Encode()
	url_ := u.String()
	cfg, err := websocket.NewConfig(url_, url_)
	if err != nil {
		return nil, err
	}
	return websocket.DialConfig(cfg)
}

// GetConfigContainerForAPIPath computes a request path to the getConfig action of containerForApi.
func GetConfigContainerForAPIPath(id string) string {
	param0 := id

	return fmt.Sprintf("/api/v2/container/%s/config", param0)
}

// Get the config of a container
func (c *Client) GetConfigContainerForAPI(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewGetConfigContainerForAPIRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewGetConfigContainerForAPIRequest create the request corresponding to the getConfig action endpoint of the containerForApi resource.
func (c *Client) NewGetConfigContainerForAPIRequest(ctx context.Context, path string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "https"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.APIKeySigner != nil {
		if err := c.APIKeySigner.Sign(req); err != nil {
			return nil, err
		}
	}
	return req, nil
}

// InspectContainerForAPIPath computes a request path to the inspect action of containerForApi.
func InspectContainerForAPIPath(id string) string {
	param0 := id

	return fmt.Sprintf("/api/v2/container/%s/inspect", param0)
}

// Return details of a container
func (c *Client) InspectContainerForAPI(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewInspectContainerForAPIRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewInspectContainerForAPIRequest create the request corresponding to the inspect action endpoint of the containerForApi resource.
func (c *Client) NewInspectContainerForAPIRequest(ctx context.Context, path string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "https"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.APIKeySigner != nil {
		if err := c.APIKeySigner.Sign(req); err != nil {
			return nil, err
		}
	}
	return req, nil
}

// ListContainerForAPIPath computes a request path to the list action of containerForApi.
func ListContainerForAPIPath() string {

	return fmt.Sprintf("/api/v2/container/list")
}

// Return a list of containers
func (c *Client) ListContainerForAPI(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewListContainerForAPIRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewListContainerForAPIRequest create the request corresponding to the list action endpoint of the containerForApi resource.
func (c *Client) NewListContainerForAPIRequest(ctx context.Context, path string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "https"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.APIKeySigner != nil {
		if err := c.APIKeySigner.Sign(req); err != nil {
			return nil, err
		}
	}
	return req, nil
}

// LogsContainerForAPIPath computes a request path to the logs action of containerForApi.
func LogsContainerForAPIPath(id string) string {
	param0 := id

	return fmt.Sprintf("/api/v2/container/%s/logs", param0)
}

// Get stdout and stderr logs from a container.
func (c *Client) LogsContainerForAPI(ctx context.Context, path string, follow *bool, since *time.Time, stderr *bool, stdout *bool, tail *string, timestamps *bool, until *time.Time) (*websocket.Conn, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "ws"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	values := u.Query()
	if follow != nil {
		tmp66 := strconv.FormatBool(*follow)
		values.Set("follow", tmp66)
	}
	if since != nil {
		tmp67 := since.Format(time.RFC3339)
		values.Set("since", tmp67)
	}
	if stderr != nil {
		tmp68 := strconv.FormatBool(*stderr)
		values.Set("stderr", tmp68)
	}
	if stdout != nil {
		tmp69 := strconv.FormatBool(*stdout)
		values.Set("stdout", tmp69)
	}
	if tail != nil {
		values.Set("tail", *tail)
	}
	if timestamps != nil {
		tmp70 := strconv.FormatBool(*timestamps)
		values.Set("timestamps", tmp70)
	}
	if until != nil {
		tmp71 := until.Format(time.RFC3339)
		values.Set("until", tmp71)
	}
	u.RawQuery = values.Encode()
	url_ := u.String()
	cfg, err := websocket.NewConfig(url_, url_)
	if err != nil {
		return nil, err
	}
	return websocket.DialConfig(cfg)
}

// RemoveContainerForAPIPath computes a request path to the remove action of containerForApi.
func RemoveContainerForAPIPath(id string) string {
	param0 := id

	return fmt.Sprintf("/api/v2/container/%s/remove", param0)
}

// remove a container
func (c *Client) RemoveContainerForAPI(ctx context.Context, path string, force bool) (*http.Response, error) {
	req, err := c.NewRemoveContainerForAPIRequest(ctx, path, force)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewRemoveContainerForAPIRequest create the request corresponding to the remove action endpoint of the containerForApi resource.
func (c *Client) NewRemoveContainerForAPIRequest(ctx context.Context, path string, force bool) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "https"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	values := u.Query()
	tmp72 := strconv.FormatBool(force)
	values.Set("force", tmp72)
	u.RawQuery = values.Encode()
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.APIKeySigner != nil {
		if err := c.APIKeySigner.Sign(req); err != nil {
			return nil, err
		}
	}
	return req, nil
}

// SetConfigContainerForAPIPath computes a request path to the setConfig action of containerForApi.
func SetConfigContainerForAPIPath(id string) string {
	param0 := id

	return fmt.Sprintf("/api/v2/container/%s/config", param0)
}

// Change the config of a container
func (c *Client) SetConfigContainerForAPI(ctx context.Context, path string, payload *ContainerConfig, contentType string) (*http.Response, error) {
	req, err := c.NewSetConfigContainerForAPIRequest(ctx, path, payload, contentType)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewSetConfigContainerForAPIRequest create the request corresponding to the setConfig action endpoint of the containerForApi resource.
func (c *Client) NewSetConfigContainerForAPIRequest(ctx context.Context, path string, payload *ContainerConfig, contentType string) (*http.Request, error) {
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
	if c.APIKeySigner != nil {
		if err := c.APIKeySigner.Sign(req); err != nil {
			return nil, err
		}
	}
	return req, nil
}

// StartContainerForAPIPath computes a request path to the start action of containerForApi.
func StartContainerForAPIPath(id string) string {
	param0 := id

	return fmt.Sprintf("/api/v2/container/%s/start", param0)
}

// start a container
func (c *Client) StartContainerForAPI(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewStartContainerForAPIRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewStartContainerForAPIRequest create the request corresponding to the start action endpoint of the containerForApi resource.
func (c *Client) NewStartContainerForAPIRequest(ctx context.Context, path string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "https"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.APIKeySigner != nil {
		if err := c.APIKeySigner.Sign(req); err != nil {
			return nil, err
		}
	}
	return req, nil
}

// StopContainerForAPIPath computes a request path to the stop action of containerForApi.
func StopContainerForAPIPath(id string) string {
	param0 := id

	return fmt.Sprintf("/api/v2/container/%s/stop", param0)
}

// stop a container
func (c *Client) StopContainerForAPI(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewStopContainerForAPIRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewStopContainerForAPIRequest create the request corresponding to the stop action endpoint of the containerForApi resource.
func (c *Client) NewStopContainerForAPIRequest(ctx context.Context, path string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "https"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.APIKeySigner != nil {
		if err := c.APIKeySigner.Sign(req); err != nil {
			return nil, err
		}
	}
	return req, nil
}

// UploadContainerForAPIPath computes a request path to the upload action of containerForApi.
func UploadContainerForAPIPath(id string) string {
	param0 := id

	return fmt.Sprintf("/api/v2/container/%s/upload", param0)
}

// Copy files to the container
func (c *Client) UploadContainerForAPI(ctx context.Context, path string, payload *UploadPayload, contentType string) (*http.Response, error) {
	req, err := c.NewUploadContainerForAPIRequest(ctx, path, payload, contentType)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewUploadContainerForAPIRequest create the request corresponding to the upload action endpoint of the containerForApi resource.
func (c *Client) NewUploadContainerForAPIRequest(ctx context.Context, path string, payload *UploadPayload, contentType string) (*http.Request, error) {
	var body bytes.Buffer
	w := multipart.NewWriter(&body)

	{
		fw, err := w.CreateFormField("allowOverwrite")
		if err != nil {
			return nil, err
		}
		tmp_AllowOverwrite := *payload.AllowOverwrite
		s := strconv.FormatBool(tmp_AllowOverwrite)
		if _, err := fw.Write([]byte(s)); err != nil {
			return nil, err
		}
	}
	{
		fw, err := w.CreateFormField("copyUIDGID")
		if err != nil {
			return nil, err
		}
		tmp_CopyUIDGID := payload.CopyUIDGID
		s := strconv.FormatBool(tmp_CopyUIDGID)
		if _, err := fw.Write([]byte(s)); err != nil {
			return nil, err
		}
	}
	{
		_, file := filepath.Split(payload.Data)
		fw, err := w.CreateFormFile("data", file)
		if err != nil {
			return nil, err
		}
		fh, err := os.Open(payload.Data)
		if err != nil {
			return nil, err
		}
		defer fh.Close()
		if _, err := io.Copy(fw, fh); err != nil {
			return nil, err
		}
	}
	{
		fw, err := w.CreateFormField("path")
		if err != nil {
			return nil, err
		}
		tmp_Path := payload.Path
		s := tmp_Path
		if _, err := fw.Write([]byte(s)); err != nil {
			return nil, err
		}
	}
	if err := w.Close(); err != nil {
		return nil, err
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
	header.Set("Content-Type", w.FormDataContentType())
	if c.APIKeySigner != nil {
		if err := c.APIKeySigner.Sign(req); err != nil {
			return nil, err
		}
	}
	return req, nil
}
