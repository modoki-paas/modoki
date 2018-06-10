package main

import (
	"context"
	"net/http"
	"net/url"
	"os/user"
	"strconv"
	"time"

	modoki "github.com/cs3238-tsuzu/modoki/client"
	"golang.org/x/net/websocket"

	goaclient "github.com/goadesign/goa/client"
)

func getHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}

// newJWTSigner returns the request signer used for authenticating
// against the jwt security scheme.
func newJWTSigner(key string) goaclient.Signer {
	return &goaclient.APIKeySigner{
		SignQuery: false,
		KeyName:   "Authorization",
		KeyValue:  key,
		Format:    "Bearer %s",
	}

}

// Get stdout and stderr logs from a container.
func modokiLogsContainer(c *modoki.Client, ctx context.Context, path string, id string, follow *bool, since *time.Time, stderr *bool, stdout *bool, tail *string, timestamps *bool, until *time.Time) (*websocket.Conn, error) {
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
