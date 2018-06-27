package main

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"time"

	modoki "github.com/cs3238-tsuzu/modoki/client"
	"golang.org/x/net/websocket"
)

// Copy files from the container
func downloadContainerHEAD(ctx context.Context, c *modoki.Client, path string, internalPath string) (*http.Response, error) {
	req, err := newDownloadContainerRequestHEAD(ctx, c, path, internalPath)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// newDownloadContainerRequestHEAD create the request corresponding to the download action endpoint of the container resource.
func newDownloadContainerRequestHEAD(ctx context.Context, c *modoki.Client, path string, internalPath string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "https"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	values := u.Query()
	values.Set("internalPath", internalPath)
	u.RawQuery = values.Encode()
	req, err := http.NewRequest("HEAD", u.String(), nil)
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

// Get stdout and stderr logs from a container.
func modokiLogsContainer(c *modoki.Client, ctx context.Context, path string, follow *bool, since *time.Time, stderr *bool, stdout *bool, tail *string, timestamps *bool, until *time.Time) (*websocket.Conn, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "ws"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	values := u.Query()
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

	var header http.Header
	if c.JWTSigner != nil {
		req, _ := http.NewRequest("ws", "", nil)
		req.URL = &u
		if err := c.JWTSigner.Sign(req); err != nil {
			return nil, err
		}

		u = *req.URL
		header = req.Header
	}
	url_ := u.String()
	cfg, err := websocket.NewConfig(url_, url_)
	if err != nil {
		return nil, err
	}

	cfg.Header = header

	return websocket.DialConfig(cfg)
}

// Copy files to the container
func uploadContainer(ctx context.Context, c *modoki.Client, path string, payload *modoki.UploadPayload, contentType string, reader io.Reader, filename string) (*http.Response, error) {
	req, err := newUploadContainerRequest(ctx, c, path, payload, contentType, reader, filename)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewUploadContainerRequest create the request corresponding to the upload action endpoint of the container resource.
func newUploadContainerRequest(ctx context.Context, c *modoki.Client, path string, payload *modoki.UploadPayload, contentType string, reader io.Reader, filename string) (*http.Request, error) {
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
		fw, err := w.CreateFormField("copyUIDGID")
		if err != nil {
			return nil, err
		}
		s := strconv.FormatBool(payload.CopyUIDGID)
		if _, err := fw.Write([]byte(s)); err != nil {
			return nil, err
		}
	}
	{
		fw, err := w.CreateFormFile("data", filename)
		if err != nil {
			return nil, err
		}
		if _, err := io.Copy(fw, reader); err != nil {
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
