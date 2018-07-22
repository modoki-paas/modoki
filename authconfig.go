package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/cs3238-tsuzu/modoki/extensions/auth"
	"github.com/pkg/errors"

	"github.com/goadesign/goa"
)

type contextKeyAuthType int

const contextKeyAuth contextKeyAuthType = iota

type authConfig map[string]interface{}

func loadAuthConfig(path string) (authConfig, error) {
	fp, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer fp.Close()

	var config authConfig
	if err := json.NewDecoder(fp).Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func initAuthMiddleware(path string, security *goa.JWTSecurity) (goa.Middleware, error) {
	config, err := loadAuthConfig(path)

	if err != nil {
		return nil, err
	}

	exts := authbase.LoadAll()

	var mws []goa.Middleware
	var names []string

	handler := func(next goa.Handler) goa.Handler {
		return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			var newReq *http.Request
			reqUpdater := func(ctx context.Context, rw http.ResponseWriter, r *http.Request) error {
				newReq = r

				return nil
			}

			for i := range mws {
				handler := mws[i](reqUpdater)

				if err := handler(ctx, rw, req); err == nil {
					newReq = newReq.WithContext(context.WithValue(newReq.Context(), contextKeyAuth, names[i]))

					break
				}
			}

			return next(ctx, rw, newReq)
		}
	}

	for k := range config {
		ext, ok := exts[k]

		if !ok {
			return nil, errors.New("No such auth type: " + k)
		}

		mw, err := ext.GetMiddleware(config[k], security)

		if err != nil {
			return nil, err
		}

		mws = append(mws, mw)
		names = append(names, k)
	}

	return handler, nil
}
