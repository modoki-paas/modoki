package main

import (
	"crypto/rsa"
	"fmt"
	"strconv"
	"time"

	"github.com/cs3238-tsuzu/modoki/app"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/goadesign/goa"
	uuid "github.com/satori/go.uuid"
)

// VironController implements the viron resource.
type VironController struct {
	*goa.Controller

	PrivateKey *rsa.PrivateKey
}

// NewVironController creates a viron controller.
func NewVironController(service *goa.Service) *VironController {
	return &VironController{Controller: service.NewController("VironController")}
}

// Authtype runs the authtype action.
func (c *VironController) Authtype(ctx *app.AuthtypeVironContext) error {
	// VironController_Authtype: start_implement

	res := app.VironauthtypeCollection{
		&app.Vironauthtype{
			Type:     "email", // メールアドレスとパスワードによる独自認証を利用する場合のtype
			Provider: "Modoki",
			URL:      "/api/v1/signin", // サインインフォームでsubmitする際のリクエストURL
			Method:   "POST",
		},
		&app.Vironauthtype{
			Type:     "signout",
			Provider: "",
			URL:      "/signout",
			Method:   "POST",
		},
	}
	return ctx.OK(res)

	// VironController_Authtype: end_implement
}

// Get runs the get action.
func (c *VironController) Get(ctx *app.GetVironContext) error {
	// VironController_Get: start_implement

	// Put your logic here

	res := &app.Vironsetting{}
	return ctx.OK(res)

	// VironController_Get: end_implement
}

// Signin runs the signin action.
func (c *VironController) Signin(ctx *app.SigninVironContext) error {
	// VironController_Signin: start_implement

	if !(ctx.Payload.ID == "admin" && ctx.Payload.Password == "password") {
		return ctx.Unauthorized()
	}

	token := jwtgo.New(jwtgo.SigningMethodRS512)
	deadline := time.Now().AddDate(1, 0, 0).Unix()
	token.Claims = jwtgo.MapClaims{
		jwtKeyUID: strconv.Itoa(1),
		"iss":     "modoki",                         // who creates the token and signs it
		"aud":     "modoki",                         // to whom the token is intended to be sent
		"exp":     deadline,                         // time when the token will expire (10 minutes from now)
		"jti":     uuid.Must(uuid.NewV4()).String(), // a unique identifier for the token
		"iat":     time.Now().Unix(),                // when the token was issued/created (now)
		"nbf":     2,                                // time before which the token is not yet valid (2 minutes ago)
		"sub":     "subject",                        // the subject/principal is whom the token is about
		"scopes":  "api:access",                     // token scope - not a standard claim
	}
	signedToken, err := token.SignedString(c.PrivateKey)
	if err != nil {
		return fmt.Errorf("failed to sign token: %s", err) // internal error
	}

	// Set auth header for client retrieval
	ctx.ResponseData.Header().Set("Authorization", "Bearer "+signedToken)

	// Send response
	return ctx.NoContent()

	// VironController_Signin: end_implement
}
