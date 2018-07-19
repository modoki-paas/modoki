package api

import (
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = API("Modoki API", func() {
	Title("Modoki API")
	Scheme("http", "https")
	BasePath("/api/v2")
	Security(JWT)
	Version("1.0.0")
})

var JWT = JWTSecurity("jwt", func() {
	Header("Authorization")
	Scope("api:access", "API access")
})
