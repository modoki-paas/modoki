package api

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = Resource("containerForApi", func() { // Defines the Operands resource
	Security(APIKey)
	BasePath(APIBasePath + "/container")

	containerActions()
})

var _ = Resource("containerForFrontend", func() { // Defines the Operands resource
	Security(JWT)
	BasePath(FrontendAPIBasePath + "/container")

	containerActions()
})

var containerActions = func() {
	Action("create", func() {
		Routing(GET("/create"))
		Description("create a new container")
		Params(func() {
			Param("name", String, func() {
				Description("Name of container and subdomain")
				Pattern("^[a-zA-Z0-9_]+$")
				Example("Hello_World01")
				MaxLength(64)
				MinLength(1)
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
		Routing(GET("/:id/start"))
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
		Routing(GET("/:id/stop"))
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
		Routing(GET("/:id/remove"))
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
		Routing(GET("/:id/inspect"))
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
		Routing(POST("/:id/upload"))
		Description("Copy files to the container")
		MultipartForm()
		Payload(UploadPayload)

		Params(func() {
			Param("id", String, "ID or name")

			Required("id")
		})

		Response(NoContent)
		Response(NotFound, ErrorMedia)
		Response(BadRequest, ErrorMedia)
		Response(RequestEntityTooLarge)
		Response(InternalServerError, ErrorMedia)
	})

	Action("download", func() {
		Routing(GET("/:id/download"), HEAD("/download"))
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
		Routing(GET("/:id/logs"))
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

	Action("getConfig", func() {
		Routing(GET("/:id/config"))
		Description("Get the config of a container")

		Params(func() {
			Param("id", String, "id or name")
		})
		Response(OK, ContainerConfigOK)
		Response(NotFound)
		Response(InternalServerError, ErrorMedia)
	})

	Action("setConfig", func() {
		Routing(POST("/:id/config"))
		Description("Change the config of a container")

		Payload(ContainerConfig)
		Params(func() {
			Param("id", String, "id or name")

			Required("id")
		})

		Response(NoContent)
		Response(NotFound)
		Response(InternalServerError, ErrorMedia)
	})

	Action("exec", func() { // WebSocket API
		Routing(GET("/:id/exec"))
		Scheme("ws")
		Description("Exec a command with attaching to a container using WebSocket(Mainly for xterm.js, using a protocol for terminado)")

		Params(func() {
			Param("id", String, "id or name")
			Param("command", ArrayOf(String), "The path to the executable file")
			Param("tty", Boolean, "Tty")

			Required("id")
		})

		Response(SwitchingProtocols)
		Response(NotFound, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})
}

var ContainerInspectRawStateMedia = MediaType("vnd.application/goa.container.inspect.raw_state+json", func() {
	Attributes(func() {
		Attribute("exitCode", Integer)
		Attribute("finishedAt", DateTime)
		Attribute("oomKilled", Boolean)
		Attribute("dead", Boolean)
		Attribute("paused", Boolean)
		Attribute("pid", Integer)
		Attribute("restarting", Boolean)
		Attribute("running", Boolean)
		Attribute("startedAt", DateTime)
		Attribute("status", String, func() {
			Enum("created", "running", "paused", "restarting", "removing", "exited", "dead")
		})

		Required("exitCode", "finishedAt", "oomKilled", "dead", "paused", "pid", "restarting", "running", "startedAt", "status")
	})

	View("default", func() {
		Attribute("exitCode")
		Attribute("finishedAt")
		Attribute("oomKilled")
		Attribute("dead")
		Attribute("paused")
		Attribute("pid")
		Attribute("restarting")
		Attribute("running")
		Attribute("startedAt")
		Attribute("status")
	})
})

var ContainerListEachMedia = MediaType("vpn.application/goa.container.list.each+json", func() {
	Attributes(func() {
		Attribute("name", String, "Assign the specified name to the container. Must match /?[a-zA-Z0-9_-]+.")
		Attribute("id", Integer, "ID")
		Attribute("image", String, "The name of the image to use when creating the container")
		Attribute("imageID", String, "The container's image ID")
		Attribute("command", String, "Command to run when starting the container")
		Attribute("created", DateTime, "The time the container was created")
		Attribute("volumes", ArrayOf(String), "Paths to mount volumes in")
		Attribute("status", String, func() {
			Enum("Creating", "Created", "Running", "Stopped", "Error")
		})

		Required("name", "id", "image", "imageID", "command", "created", "status", "volumes")
	})

	View("default", func() {
		Attribute("name")
		Attribute("id")
		Attribute("image")
		Attribute("imageID")
		Attribute("command")
		Attribute("created")
		Attribute("volumes")
		Attribute("status")
	})
})

var ContainerInspectMedia = MediaType("vpn.application/goa.container.inspect+json", func() {
	Attributes(func() {
		Attribute("name", String, "Assign the specified name to the container. Must match /?[a-zA-Z0-9_-]+.")
		Attribute("id", Integer, "ID")
		Attribute("image", String, "The name of the image to use when creating the container")
		Attribute("imageID", String, "The container's image ID")
		Attribute("path", String, "The path to the command being run")
		Attribute("args", ArrayOf(String), "The arguments to the command being run")
		Attribute("created", DateTime, "The time the container was created")
		Attribute("volumes", ArrayOf(String), "Paths to mount volumes in")

		Attribute("status", String, func() {
			Enum("Image Downloading", "Created", "Running", "Stopped", "Error")
		})
		Attribute("raw_state", ContainerInspectRawStateMedia)

		Required("name", "id", "image", "imageID", "path", "args", "created", "status", "raw_state", "volumes")
	})

	View("default", func() {
		Attribute("name")
		Attribute("id")
		Attribute("image")
		Attribute("imageID")
		Attribute("path")
		Attribute("args")
		Attribute("created")
		Attribute("volumes")
		Attribute("status")
		Attribute("raw_state")
	})
})

var ContainerCreateOK = MediaType("vnd.application/goa.container.create.results+json", func() {
	Description("The results of container creation")
	Attributes(func() { // Defines the media type attributes
		Attribute("id", Integer, "container id")
		Attribute("endpoints", ArrayOf(String), "endpoint URL")

		Required("id", "endpoints")
	})

	View("default", func() {
		Attribute("id")
		Attribute("endpoints")
	})
})

var ContainerDownloadOK = MediaType("vpn.application/goa.container.download.result+json", func() {
	Attributes(func() {
		Attribute("file", File)
	})

	View("default", func() {
		Attribute("file")
	})
})

var ContainerConfig = Type("ContainerConfig", func() {
	Attribute("defaultShell", String)
})

var ContainerConfigOK = MediaType("vpn.application/goa.container.config+json", func() {
	Reference(ContainerConfig)
	Attributes(func() {
		Attribute("defaultShell")
	})

	View("default", func() {
		Attribute("defaultShell")
	})
})

var UploadPayload = Type("UploadPayload", func() {
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

	Required("path", "data", "allowOverwrite", "copyUIDGID")
})
