package authbase

import (
	"github.com/cs3238-tsuzu/modoki/extensions/auth/auth0"
	"github.com/cs3238-tsuzu/modoki/extensions/auth/file"
	"github.com/cs3238-tsuzu/modoki/extensions/auth/firebase"
	"github.com/goadesign/goa"
)

// AuthExtension is a base interface for authentication extension
type AuthExtension interface {
	GetName() string
	GetMiddleware(config interface{}, security *goa.JWTSecurity) (goa.Middleware, error)
}

func LoadAll() map[string]AuthExtension {
	fileAuth := &fileauth.FileAuth{}
	firebaseAuth := &firebaseauth.FirebaseAuth{}
	auth0auth := &auth0auth.Auth0auth{}

	return map[string]AuthExtension{
		fileAuth.GetName():     fileAuth,
		firebaseAuth.GetName(): firebaseAuth,
		auth0auth.GetName():    auth0auth,
	}
}
