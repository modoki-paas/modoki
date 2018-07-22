package fileauth

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	jwtgo "github.com/dgrijalva/jwt-go"

	"github.com/mitchellh/mapstructure"

	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware/security/jwt"
)

type fileconfig struct {
	File string `mapstructure:"file"`
}

// FileAuth is an AuthExtension for local files
type FileAuth struct {
}

// FileAuth implements AuthExtension.GetName
func (a *FileAuth) GetName() string {
	return "file"
}

// GetMiddleware implements AuthExtension.GetMiddleware
func (a *FileAuth) GetMiddleware(cfg interface{}, security *goa.JWTSecurity) (goa.Middleware, error) {
	if cfg == nil {
		return nil, errors.New("Config must be not nil")
	}

	var config fileconfig
	mapstructure.Decode(cfg, &config)

	keys, err := LoadJWTPublicKeys(config.File)

	if err != nil {
		return nil, err
	}

	return jwt.New(jwt.NewSimpleResolver(keys), nil, security), nil
}

// LoadJWTPublicKeys loads PEM encoded RSA public keys used to validata and decrypt the JWT.
func LoadJWTPublicKeys(path string) ([]jwt.Key, error) {
	keyFiles, err := filepath.Glob(path)
	if err != nil {
		return nil, err
	}
	keys := make([]jwt.Key, len(keyFiles))
	for i, keyFile := range keyFiles {
		pem, err := ioutil.ReadFile(keyFile)
		if err != nil {
			return nil, err
		}
		key, err := jwtgo.ParseRSAPublicKeyFromPEM([]byte(pem))
		if err != nil {
			return nil, fmt.Errorf("failed to load key %s: %s", keyFile, err)
		}
		keys[i] = key
	}
	if len(keys) == 0 {
		return nil, fmt.Errorf("couldn't load public keys for JWT security")
	}

	return keys, nil
}
