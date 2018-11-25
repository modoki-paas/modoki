package main

import (
	"context"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/goadesign/goa/middleware/security/jwt"
	"github.com/modoki-paas/modoki/const"
)

func GetUIDFromContext(ctx context.Context) (string, error) {
	if token := jwt.ContextJWT(ctx); token != nil {
		claims := token.Claims.(jwtgo.MapClaims)

		uid := claims[constants.JWTKeyUID].(string)

		return uid, nil
	}

	if uid, ok := ctx.Value(constants.UIDContextKeyForAPIKey).(string); ok {
		return uid, nil
	}

	return "", errUnauthorized
}
