package api

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = Resource("viron", func() {
	Action("authtype", func() {
		NoSecurity()
		Description("Get viron authtype")
		Routing(GET("/viron_authtype"))
		Response(OK, func() {
			Media(CollectionOf(VironAuthType))
		})
	})
	Action("get", func() {
		Description("Get viron menu")
		Routing(GET("/viron"))
		Response(OK, func() {
			Media(VironSetting)
		})
	})
	Action("signin", func() {
		Description("Creates a valid JWT")
		NoSecurity()
		Payload(SigninPayload)
		Routing(POST("/signin"))
		Response(NoContent, func() {
			Headers(func() {
				Header("Authorization", String, "Generated JWT")
			})
		})
		Response(Unauthorized)
		Response(InternalServerError)
	})
})

var SigninPayload = Type("SigninPayload", func() {
	Member("id", String, "ID or Email", func() {
		Example("identify key")
	})
	Member("password", String, "Password", func() {
		MaxLength(256)
	})
	Required("id", "password")
})

var PostPayload = Type("PostPayload", func() {
	Member("url_name", String, "url name", func() {
		Example("hello")
	})
	Member("title", String, "title", func() {
		Example("hello viron-goa example")
		MaxLength(120)
	})
	Member("contents", String, "contents", func() {
		Example("Hi gopher")
		MaxLength(120)
	})
	Member("status", String, "status", func() {
		Example("draft")
		Enum("draft", "published")
	})
	Member("published_at", DateTime, "published_at", func() {
	})
	Required("url_name", "title", "contents", "status")
})

// VironAuthType
var VironAuthType = MediaType("application/vnd.vironauthtype+json", func() {
	Description("viron authtype media")
	Attributes(func() {
		Attribute("type", String, "auth type", func() {
			Example("signin")
		})
		Attribute("provider", String, "auth provider", func() {
			Example("")
		})
		Attribute("url", String, "url", func() {
			Example("/signin")
		})
		Attribute("method", String, "method", func() {
			Example("POST")
		})
		Required("type", "provider", "url", "method")
	})
	View("default", func() {
		Attribute("type")
		Attribute("provider")
		Attribute("url")
		Attribute("method")
	})
})

var VironSetting = MediaType("application/vnd.vironsetting+json", func() {
	Description("viron setting")
	Attributes(func() {
		Attribute("name", String, "name")
		Attribute("color", String, "color")
		Attribute("theme", String, "theme")
		Attribute("thumbnail", String, "thumbnail")
		Attribute("tags", ArrayOf(String), "tags")
		Attribute("pages", ArrayOf(VironPage), "pages")
		Required("name", "color", "theme", "pages", "tags", "thumbnail")
	})
	View("default", func() {
		Attribute("name")
		Attribute("color")
		Attribute("theme")
		Attribute("thumbnail")
		Attribute("tags")
		Attribute("pages")
	})
})

var VironPage = MediaType("application/vnd.vironpage+json", func() {
	Description("viron page")
	Attributes(func() {
		Attribute("id", String, "id")
		Attribute("name", String, "name")
		Attribute("section", String, "section")
		Attribute("group", String, "group")
		Attribute("components", ArrayOf(VironComponent), "pages")
		Required("id", "name", "section", "group", "components")
	})
	View("default", func() {
		Attribute("id")
		Attribute("name")
		Attribute("section")
		Attribute("group")
		Attribute("components")
	})
})

var VironComponent = MediaType("application/vnd.vironcomponent+json", func() {
	Description("viron component")
	Attributes(func() {
		Attribute("name", String, "name")
		Attribute("style", String, "style")
		Attribute("api", VironAPI, "api")
		Attribute("pagination", Boolean, "pagination")
		Attribute("primary", String, "primary key")
		Attribute("query", ArrayOf(VironQuery), "query")
		Attribute("table_labels", ArrayOf(String), "table label")
		Required("name", "style", "api")
	})
	View("default", func() {
		Attribute("name")
		Attribute("style")
		Attribute("api")
		Attribute("pagination")
		Attribute("primary")
		Attribute("query")
		Attribute("table_labels")
	})
})

var VironQuery = MediaType("application/vnd.query+json", func() {
	Description("viron query")
	Attributes(func() {
		Attribute("key", String, "key")
		Attribute("type", String, "type")
		Required("key", "type")
	})
	View("default", func() {
		Attribute("key")
		Attribute("type")
	})
})

var VironAPI = MediaType("application/vnd.vironapi+json", func() {
	Description("viron api")
	Attributes(func() {
		Attribute("method", String, "name")
		Attribute("path", String, "path")
		Required("method", "path")
	})
	View("default", func() {
		Attribute("method")
		Attribute("path")
	})
})
