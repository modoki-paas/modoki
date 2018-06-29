package api

import (
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = API("Modoki API", func() {
	Title("Modoki API Documentation")
	Scheme("http", "https")
	Host("localhost:4434")
	BasePath("/api/v2")
	Security(JWT)
})

var JWT = JWTSecurity("jwt", func() {
	Header("Authorization")
	Scope("api:access", "API access")
})
