package main

import (
	"context"
	"fmt"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/goadesign/goa/middleware/security/jwt"
)

func GetUIDFromJWT(ctx context.Context) (int, error) {
	token := jwt.ContextJWT(ctx)
	if token == nil {
		return 0, fmt.Errorf("JWT token is missing from context") // internal error
	}
	claims := token.Claims.(jwtgo.MapClaims)

	uid := claims[jwtKeyUID].(int)

	return uid, nil
}
