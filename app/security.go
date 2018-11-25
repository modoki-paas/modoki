// Code generated by goagen v1.4.0, DO NOT EDIT.
//
// API "Modoki API": Application Security
//
// Command:
// $ goagen
// --design=github.com/modoki-paas/modoki/design
// --out=$(GOPATH)/src/github.com/modoki-paas/modoki
// --version=v1.4.0

package app

import (
	"context"
	"github.com/goadesign/goa"
	"net/http"
)

type (
	// Private type used to store auth handler info in request context
	authMiddlewareKey string
)

// UseAPIKeyMiddleware mounts the api_key auth middleware onto the service.
func UseAPIKeyMiddleware(service *goa.Service, middleware goa.Middleware) {
	service.Context = context.WithValue(service.Context, authMiddlewareKey("api_key"), middleware)
}

// NewAPIKeySecurity creates a api_key security definition.
func NewAPIKeySecurity() *goa.APIKeySecurity {
	def := goa.APIKeySecurity{
		In:   goa.LocHeader,
		Name: "X-Shared-Secret",
	}
	return &def
}

// UseJWTMiddleware mounts the jwt auth middleware onto the service.
func UseJWTMiddleware(service *goa.Service, middleware goa.Middleware) {
	service.Context = context.WithValue(service.Context, authMiddlewareKey("jwt"), middleware)
}

// NewJWTSecurity creates a jwt security definition.
func NewJWTSecurity() *goa.JWTSecurity {
	def := goa.JWTSecurity{
		In:       goa.LocHeader,
		Name:     "Authorization",
		TokenURL: "",
		Scopes: map[string]string{
			"api:access": "API access",
		},
	}
	return &def
}

// handleSecurity creates a handler that runs the auth middleware for the security scheme.
func handleSecurity(schemeName string, h goa.Handler, scopes ...string) goa.Handler {
	return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		scheme := ctx.Value(authMiddlewareKey(schemeName))
		am, ok := scheme.(goa.Middleware)
		if !ok {
			return goa.NoAuthMiddleware(schemeName)
		}
		ctx = goa.WithRequiredScopes(ctx, scopes)
		return am(h)(ctx, rw, req)
	}
}
