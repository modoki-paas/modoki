// Code generated by goagen v1.3.1, DO NOT EDIT.
//
// API "Modoki API": container Resource Client
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

// CreateContainerPath computes a request path to the create action of container.
func CreateContainerPath() string {

	return fmt.Sprintf("/api/v1/container/create")
}

// create a new container
func (c *Client) CreateContainer(ctx context.Context, path string, image string, name string, command []string, entrypoint []string, env []string, sslRedirect *bool, volumes []string, workingDir *string) (*http.Response, error) {
	req, err := c.NewCreateContainerRequest(ctx, path, image, name, command, entrypoint, env, sslRedirect, volumes, workingDir)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewCreateContainerRequest create the request corresponding to the create action endpoint of the container resource.
func (c *Client) NewCreateContainerRequest(ctx context.Context, path string, image string, name string, command []string, entrypoint []string, env []string, sslRedirect *bool, volumes []string, workingDir *string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "https"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	values := u.Query()
	values.Set("image", image)
	values.Set("name", name)
	for _, p := range command {
		tmp21 := p
		values.Add("command", tmp21)
	}
	for _, p := range entrypoint {
		tmp22 := p
		values.Add("entrypoint", tmp22)
	}
	for _, p := range env {
		tmp23 := p
		values.Add("env", tmp23)
	}
	if sslRedirect != nil {
		tmp24 := strconv.FormatBool(*sslRedirect)
		values.Set("sslRedirect", tmp24)
	}
	for _, p := range volumes {
		tmp25 := p
		values.Add("volumes", tmp25)
	}
	if workingDir != nil {
		values.Set("workingDir", *workingDir)
	}
	u.RawQuery = values.Encode()
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

// DownloadContainerPath computes a request path to the download action of container.
func DownloadContainerPath() string {

	return fmt.Sprintf("/api/v1/container/download")
}

// Copy files from the container
func (c *Client) DownloadContainer(ctx context.Context, path string, id string, internalPath string) (*http.Response, error) {
	req, err := c.NewDownloadContainerRequest(ctx, path, id, internalPath)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewDownloadContainerRequest create the request corresponding to the download action endpoint of the container resource.
func (c *Client) NewDownloadContainerRequest(ctx context.Context, path string, id string, internalPath string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "https"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	values := u.Query()
	values.Set("id", id)
	values.Set("internalPath", internalPath)
	u.RawQuery = values.Encode()
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

// InspectContainerPath computes a request path to the inspect action of container.
func InspectContainerPath() string {

	return fmt.Sprintf("/api/v1/container/inspect")
}

// Return details of a container
func (c *Client) InspectContainer(ctx context.Context, path string, id string) (*http.Response, error) {
	req, err := c.NewInspectContainerRequest(ctx, path, id)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewInspectContainerRequest create the request corresponding to the inspect action endpoint of the container resource.
func (c *Client) NewInspectContainerRequest(ctx context.Context, path string, id string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "https"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	values := u.Query()
	values.Set("id", id)
	u.RawQuery = values.Encode()
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

// ListContainerPath computes a request path to the list action of container.
func ListContainerPath() string {

	return fmt.Sprintf("/api/v1/container/list")
}

// Return a list of containers
func (c *Client) ListContainer(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewListContainerRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewListContainerRequest create the request corresponding to the list action endpoint of the container resource.
func (c *Client) NewListContainerRequest(ctx context.Context, path string) (*http.Request, error) {
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

// LogsContainerPath computes a request path to the logs action of container.
func LogsContainerPath() string {

	return fmt.Sprintf("/api/v1/container/logs")
}

// Get stdout and stderr logs from a container.
func (c *Client) LogsContainer(ctx context.Context, path string, id string, follow *bool, since *time.Time, stderr *bool, stdout *bool, tail *string, timestamps *bool, until *time.Time) (*websocket.Conn, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "ws"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	values := u.Query()
	values.Set("id", id)
	if follow != nil {
		tmp26 := strconv.FormatBool(*follow)
		values.Set("follow", tmp26)
	}
	if since != nil {
		tmp27 := since.Format(time.RFC3339)
		values.Set("since", tmp27)
	}
	if stderr != nil {
		tmp28 := strconv.FormatBool(*stderr)
		values.Set("stderr", tmp28)
	}
	if stdout != nil {
		tmp29 := strconv.FormatBool(*stdout)
		values.Set("stdout", tmp29)
	}
	if tail != nil {
		values.Set("tail", *tail)
	}
	if timestamps != nil {
		tmp30 := strconv.FormatBool(*timestamps)
		values.Set("timestamps", tmp30)
	}
	if until != nil {
		tmp31 := until.Format(time.RFC3339)
		values.Set("until", tmp31)
	}
	u.RawQuery = values.Encode()
	url_ := u.String()
	cfg, err := websocket.NewConfig(url_, url_)
	if err != nil {
		return nil, err
	}
	return websocket.DialConfig(cfg)
}

// RemoveContainerPath computes a request path to the remove action of container.
func RemoveContainerPath() string {

	return fmt.Sprintf("/api/v1/container/remove")
}

// remove a container
func (c *Client) RemoveContainer(ctx context.Context, path string, force bool, id string) (*http.Response, error) {
	req, err := c.NewRemoveContainerRequest(ctx, path, force, id)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewRemoveContainerRequest create the request corresponding to the remove action endpoint of the container resource.
func (c *Client) NewRemoveContainerRequest(ctx context.Context, path string, force bool, id string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "https"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	values := u.Query()
	tmp32 := strconv.FormatBool(force)
	values.Set("force", tmp32)
	values.Set("id", id)
	u.RawQuery = values.Encode()
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

// StartContainerPath computes a request path to the start action of container.
func StartContainerPath() string {

	return fmt.Sprintf("/api/v1/container/start")
}

// start a container
func (c *Client) StartContainer(ctx context.Context, path string, id string) (*http.Response, error) {
	req, err := c.NewStartContainerRequest(ctx, path, id)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewStartContainerRequest create the request corresponding to the start action endpoint of the container resource.
func (c *Client) NewStartContainerRequest(ctx context.Context, path string, id string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "https"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	values := u.Query()
	values.Set("id", id)
	u.RawQuery = values.Encode()
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

// StopContainerPath computes a request path to the stop action of container.
func StopContainerPath() string {

	return fmt.Sprintf("/api/v1/container/stop")
}

// stop a container
func (c *Client) StopContainer(ctx context.Context, path string, id string) (*http.Response, error) {
	req, err := c.NewStopContainerRequest(ctx, path, id)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewStopContainerRequest create the request corresponding to the stop action endpoint of the container resource.
func (c *Client) NewStopContainerRequest(ctx context.Context, path string, id string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "https"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	values := u.Query()
	values.Set("id", id)
	u.RawQuery = values.Encode()
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

// UploadContainerPath computes a request path to the upload action of container.
func UploadContainerPath() string {

	return fmt.Sprintf("/api/v1/container/upload")
}

// Copy files to the container
func (c *Client) UploadContainer(ctx context.Context, path string, payload *UploadPayload, contentType string) (*http.Response, error) {
	req, err := c.NewUploadContainerRequest(ctx, path, payload, contentType)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewUploadContainerRequest create the request corresponding to the upload action endpoint of the container resource.
func (c *Client) NewUploadContainerRequest(ctx context.Context, path string, payload *UploadPayload, contentType string) (*http.Request, error) {
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	{
		fw, err := w.CreateFormField("allowOverwrite")
		if err != nil {
			return nil, err
		}
		s := strconv.FormatBool(payload.AllowOverwrite)
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
		fw, err := w.CreateFormField("id")
		if err != nil {
			return nil, err
		}
		s := payload.ID
		if _, err := fw.Write([]byte(s)); err != nil {
			return nil, err
		}
	}
	{
		fw, err := w.CreateFormField("path")
		if err != nil {
			return nil, err
		}
		s := payload.Path
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
	if c.JWTSigner != nil {
		if err := c.JWTSigner.Sign(req); err != nil {
			return nil, err
		}
	}
	return req, nil
}
