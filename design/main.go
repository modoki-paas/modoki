package api

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = API("Modoki API", func() {
	Title("Modoki API Documentation")
	Scheme("http", "https")
	Host("localhost:4434")
	BasePath("/api/v1")
	Security(JWT)
})

var JWT = JWTSecurity("jwt", func() {
	Header("Authorization")
	Scope("api:access", "API access")
})

var _ = Resource("container", func() { // Defines the Operands resource
	Security(JWT)
	BasePath("/container")

	Action("create", func() {
		Routing(GET("/create"))
		Description("create a new container")
		Params(func() {
			Param("name", String, func() {
				Description("Name of container and subdomain")
				Pattern("^[a-zA-Z0-9_]+$")
				Example("Hello_World01")
				MaxLength(64)
			})
			Param("image", String, "Name of image")
			Param("command", ArrayOf(String), "Command to run specified as a string or an array of strings.")
			Param("entrypoint", ArrayOf(String), "The entry point for the container as a string or an array of strings")
			Param("env", ArrayOf(String), "Environment variables")
			Param("volumes", ArrayOf(String), "Path to volumes in a container")
			Param("workingDir", String, "Current directory (PWD) in the command will be launched")
			Param("sslRedirect", Boolean, func() {
				Description("Whether HTTP is redirected to HTTPS")

				Default(true)
			})

			Required("name", "image")
		})
		Response(OK, func() {
			Status(200)
			Media(ContainerCreateOK)
		})
		Response(BadRequest, ErrorMedia)
		Response("Conflict", func() {
			Status(409)
			Media(ErrorMedia)
		})
		Response(InternalServerError, ErrorMedia)
	})

	Action("start", func() {
		Routing(GET("/start"))
		Description("start a container")
		Params(func() {
			Param("id", String, "id or name")

			Required("id")
		})
		Response(NoContent)
		Response(NotFound)
		Response(InternalServerError, ErrorMedia)
	})

	Action("stop", func() {
		Routing(GET("/stop"))
		Description("stop a container")
		Params(func() {
			Param("id", String, "id or name")

			Required("id")
		})
		Response(NoContent)
		Response(NotFound)
		Response(InternalServerError, ErrorMedia)
	})

	Action("remove", func() {
		Routing(GET("/remove"))
		Description("remove a container")
		Params(func() {
			Param("id", String, "id or name")
			Param("force", Boolean, func() {
				Default(false)
				Description("If the container is running, kill it before removing it.")
			})

			Required("id", "force")
		})
		Response(NoContent)
		Response(NotFound)
		Response("RunningContainer", func() {
			Status(409)
			Description("You cannot remove a running container")
		})
		Response(InternalServerError, ErrorMedia)
	})

	Action("inspect", func() {
		Routing(GET("/inspect"))
		Description("Return details of a container")

		Params(func() {
			Param("id", String, "ID or name")

			Required("id")
		})

		Response(OK, ContainerInspectMedia)
		Response(NotFound)
		Response(InternalServerError, ErrorMedia)
	})
	Action("list", func() {
		Routing(GET("/list"))
		Description("Return a list of containers")

		Response(OK, CollectionOf(ContainerListEachMedia))
		Response(InternalServerError, ErrorMedia)
	})
	Action("upload", func() {
		Routing(POST("/upload"))
		Description("Copy files to the container")
		MultipartForm()
		Payload(UploadPayload)

		Response(NoContent)
		Response(NotFound, ErrorMedia)
		Response(BadRequest, ErrorMedia)
		Response(RequestEntityTooLarge)
		Response(InternalServerError, ErrorMedia)
	})

	Action("download", func() {
		Routing(GET("/download"), HEAD("/download"))
		Description("Copy files from the container")
		Params(func() {
			Param("id", String, "ID or name")
			Param("internalPath", String, "Path in the container to save files")

			Required("id", "internalPath")
		})

		Response(OK, "application/octet-stream")
		Response(NotFound, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("logs", func() { // WebSocket API
		Routing(GET("/logs"))
		Scheme("ws")
		Description("Get stdout and stderr logs from a container.")

		Params(func() {
			Param("id", String, "id or name")

			Param("follow", Boolean, func() {
				Default(false)
			})
			Param("stdout", Boolean, func() {
				Default(false)
			})
			Param("stderr", Boolean, func() {
				Default(false)
			})
			Param("since", DateTime)
			Param("until", DateTime)
			Param("timestamps", Boolean, func() {
				Default(false)
			})
			Param("tail", String, func() {
				Default("all")
			})

			Required("id")
		})

		Response(SwitchingProtocols)
		Response(NotFound, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})
})

var UploadPayload = Type("UploadPayload", func() {
	Attribute("id", String, "ID or name")
	Attribute("path", String, "Path in the container to save files")
	Attribute("data", File, "File tar archive")
	Attribute("allowOverwrite", Boolean, func() {
		Description("Allow for a existing directory to be replaced by a file")
		Default(false)
	})
	Attribute("copyUIDGID", Boolean, func() {
		Description("Copy all uid/gid information")
		Default(false)
	})

	Required("id", "path", "data", "copyUIDGID")
})
