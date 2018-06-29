package api

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var UserAuthorizedKeyType = Type("UserAuthorizedKey", func() {
	Attribute("key", String, func() {
		MaxLength(2048)
	})
	Attribute("label", String, func() {
		Pattern("^[a-zA-Z0-9_]+$")
		MaxLength(32)
		MinLength(1)
	})

	Required("key", "label")
})
var UserConfig = Type("UserConfig", func() {
	Attribute("defaultShell", String)
	Attribute("authorizedKeys", ArrayOf(UserAuthorizedKeyType))
})

var UserAuthorizedKeyOK = MediaType("vpn.application/goa.user.authorizedKey", func() {
	Reference(UserAuthorizedKeyType)

	Attributes(func() {
		Attribute("key")
		Attribute("label")

		Required("key", "label")
	})

	View("default", func() {
		Attribute("key")
		Attribute("label")
	})
})
var UserConfigOK = MediaType("vpn.application/goa.user.config", func() {
	Attributes(func() {
		Attribute("defaultShell", String)
		Attribute("authorizedKeys", CollectionOf(UserAuthorizedKeyOK))

		Required("defaultShell", "authorizedKeys")
	})

	View("default", func() {
		Attribute("defaultShell")
		Attribute("authorizedKeys")
	})
})

var UserDefaultShellOK = MediaType("vpn.application/goa.user.defaultShell", func() {
	Attributes(func() {
		Attribute("defaultShell", String)

		Required("defaultShell")
	})

	View("default", func() {
		Attribute("defaultShell")
	})
})

var _ = Resource("user", func() {
	Security(JWT)
	BasePath("/user")

	Action("getConfig", func() {
		Routing(GET("/config"))

		Response(OK, UserConfigOK)
		Response(InternalServerError, ErrorMedia)
	})

	Action("getDefaultShell", func() {
		Routing(GET("/config/defaultShell"))

		Response(OK, UserDefaultShellOK)
		Response(InternalServerError, ErrorMedia)
	})

	Action("setDefaultShell", func() {
		Routing(POST("/config/defaultShell"))

		Params(func() {
			Param("defaultShell", String)

			Required("defaultShell")
		})

		Response(NoContent)
		Response(InternalServerError, ErrorMedia)
	})

	Action("setAuthorizedKeys", func() {
		Routing(POST("/config/authorizedKeys"))

		Payload(ArrayOf(UserAuthorizedKeyType))

		Response(NoContent)
		Response(InternalServerError, ErrorMedia)
	})
	Action("addAuthorizedKeys", func() {
		Routing(PUT("/config/authorizedKeys"))

		Payload(UserAuthorizedKeyType)

		Response(NoContent)
		Response(InternalServerError, ErrorMedia)
	})
	Action("removeAuthorizedKeys", func() {
		Routing(DELETE("/config/authorizedKeys"))

		Params(func() {
			Param("label", String)

			Required("label")
		})

		Response(NoContent)
		Response(NotFound)
		Response(InternalServerError, ErrorMedia)
	})
	Action("listAuthorizedKeys", func() {
		Routing(GET("/config/authorizedKeys"))

		Response(OK, CollectionOf(UserAuthorizedKeyOK))
		Response(NotFound)
		Response(InternalServerError, ErrorMedia)
	})
})
