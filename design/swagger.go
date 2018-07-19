package api

import (
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = Resource("swagger", func() { // Defines the Operands resource
	Origin("*", func() { // CORS policy that applies to all actions and file servers
		Methods("GET") // of "public" resource
	})
	NoSecurity()

	Files("/api/v2/swagger/swagger.json", "./swagger/swagger.json")
	Files("/api/v2/swagger/swagger.yaml", "./swagger/swagger.yaml")

})
