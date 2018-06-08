package api

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var ContainerInspectRawStateMedia = MediaType("vnd.application/goa.container.inspect.raw_state", func() {
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

var ContainerListEachMedia = MediaType("vpn.application/goa.container.list.each", func() {
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

var ContainerInspectMedia = MediaType("vpn.application/goa.container.inspect", func() {
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

		Required("id")
	})

	View("default", func() {
		Attribute("id")
	})
})
