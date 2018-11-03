package api

import (
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = API("Modoki API", func() {
	Title("Modoki API")
	Scheme("http", "https")
	Version("1.0.0")
})

const APIVersion = "v2"

const APIBasePath = "/api/" + APIVersion
const FrontendAPIBasePath = "/frontend/" + APIVersion

var JWT = JWTSecurity("jwt", func() {
	Header("Authorization")
	Scope("api:access", "API access")
})

var APIKey = APIKeySecurity("api_key", func() {
	Header("X-Shared-Secret")
})
