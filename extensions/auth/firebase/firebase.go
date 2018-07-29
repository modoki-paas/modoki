package firebaseauth

import (
	"errors"
	"net/http"

	"github.com/modoki-paas/firebase-goa-util"

	"github.com/mitchellh/mapstructure"

	"github.com/goadesign/goa"
)

type firebaseConfig struct {
	Aud string `mapstructure:"aud"`
	Iss string `mapstructure:"iss"`
}

// FirebaseAuth is an AuthExtension for Firebase
type FirebaseAuth struct {
}

// GetName implements AuthExtension.GetName
func (a *FirebaseAuth) GetName() string {
	return "firebase"
}

// GetMiddleware implements AuthExtension.GetMiddleware
func (a *FirebaseAuth) GetMiddleware(cfg interface{}, _ *goa.JWTSecurity) (goa.Middleware, error) {
	if cfg == nil {
		return nil, errors.New("Config must be not nil")
	}

	var config firebaseConfig
	mapstructure.Decode(cfg, &config)

	getPemCert := firebasegoa.NewGetPemCert()

	middleware := firebasegoa.NewJWTMiddleware(config.Aud, config.Iss, getPemCert)
	middleware.Options.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err string) {}

	handler := firebasegoa.BridgeMiddlewareHandler{
		Middleware: middleware,
	}

	fn := handler.Handle

	return fn, nil
}
