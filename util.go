package main

import (
	"context"
	"fmt"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/goadesign/goa/middleware/security/jwt"
)

func GetUIDFromJWT(ctx context.Context) (string, error) {
	token := jwt.ContextJWT(ctx)
	if token == nil {
		return "", fmt.Errorf("JWT token is missing from context") // internal error
	}

	claims := token.Claims.(jwtgo.MapClaims)

	uidStr := claims[jwtKeyUID].(string)

	return uidStr, nil
}
