package main

import (
	"os/user"

	goaclient "github.com/goadesign/goa/client"
)

func getHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}

// newJWTSigner returns the request signer used for authenticating
// against the jwt security scheme.
func newJWTSigner(key string) goaclient.Signer {
	return &goaclient.APIKeySigner{
		SignQuery: false,
		KeyName:   "Authorization",
		KeyValue:  key,
		Format:    "Bearer %s",
	}

}
