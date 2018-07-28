package auth0auth

import (
	"errors"
	"fmt"

	"github.com/auth0/go-jwt-middleware"

	"github.com/modoki-paas/auth0-goa-util"

	"github.com/mitchellh/mapstructure"

	"github.com/goadesign/goa"
)

type auth0config struct {
	Aud  string `mapstructure:"aud"`
	Iss  string `mapstructure:"iss"`
	JWKS string `mapstructure:"jwks"`
}

// Auth0auth is an AuthExtension for Auth0
type Auth0auth struct {
}

// GetName implements AuthExtension.GetName
func (a *Auth0auth) GetName() string {
	return "auth0"
}

// GetMiddleware implements AuthExtension.GetMiddleware
func (a *Auth0auth) GetMiddleware(cfg interface{}, security *goa.JWTSecurity) (goa.Middleware, error) {
	if cfg == nil {
		return nil, errors.New("Config must be not nil")
	}

	var config auth0config
	mapstructure.Decode(cfg, &config)

	getPemCert := auth0goa.NewGetPemCert(config.JWKS)

	middleware := auth0goa.NewJWTMiddleware(config.Aud, config.Iss, getPemCert)

	if security.In == goa.LocHeader {
		middleware.Options.Extractor = jwtmiddleware.FromAuthHeader // Must be Authorization header
	} else if security.In == goa.LocQuery {
		middleware.Options.Extractor = jwtmiddleware.FromParameter(security.Name)
	} else {
		return nil, fmt.Errorf("whoops, security scheme with location (in) %q not supported", security.In)
	}

	handler := auth0goa.BridgeMiddlewareHandler{
		Middleware: middleware,
	}

	fn := handler.Handle

	return fn, nil
}
