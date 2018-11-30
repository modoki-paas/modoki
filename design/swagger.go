package api

import (
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = Resource("swagger", func() { // Defines the Operands resource
	Origin("*", func() { // CORS policy that applies to all actions and file servers
		Methods("GET") // of "public" resource
		Credentials()
	})
	NoSecurity()

	Files("/swagger/swagger.json", "./swagger/swagger.json")
	Files("/swagger/swagger.yaml", "./swagger/swagger.yaml")

})
