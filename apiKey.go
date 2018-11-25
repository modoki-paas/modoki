package main

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/goadesign/goa"
	"github.com/jmoiron/sqlx"
	"github.com/modoki-paas/modoki/app"
	"github.com/modoki-paas/modoki/const"
	"github.com/pkg/errors"
)

var errUnauthorized = errors.New("unauthorized")

func newAPIKeyMiddleware(db *sqlx.DB) goa.Middleware {
	scheme := app.NewAPIKeySecurity()

	apiKeyChecker := func(apiKey string) (string, error) {
		var uid string
		if err := db.Get(&uid, "SELECT uid FROM apiKeys WHERE apiKey=?", apiKey); err != nil {
			if err == sql.ErrNoRows {
				return "", errUnauthorized
			}

			return "", err
		}

		return uid, nil

	}

	return func(h goa.Handler) goa.Handler {
		return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			key := req.Header.Get(scheme.Name)

			id, err := apiKeyChecker(key)

			if err != nil {
				if err == errUnauthorized {
					return goa.NewErrorClass("unauthozied", 401)("api key is missing or invalid")
				}

				return goa.ErrInternal(err)
			}

			return h(context.WithValue(ctx, constants.UIDContextKeyForAPIKey, id), rw, req)
		}
	}
}
