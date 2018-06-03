package api

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var ContainerType = Type("Container", func() {
	Attribute("name", String, "Assign the specified name to the container. Must match /?[a-zA-Z0-9_-]+.")
	Attribute("id", String, "ID")
	Attribute("image", String, "The name of the image to use when creating the container")
	Attribute("imageID", String, "The container's image ID")
	Attribute("path", String, "The path to the command being run")
	Attribute("args", ArrayOf(String), "The arguments to the command being run")

	Attribute("created", DateTime, "The time the container was created")
	Attribute("state", func() {
		Attribute("error")
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

		Required("error", "exitCode", "finishedAt", "oomKilled", "dead", "paused", "pid", "restarting", "running", "startedAt", "status")
	})
	Attribute("volumes", ArrayOf(String))

	Required("name", "id", "image", "imageID", "path", "args", "created", "state", "volumes")
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
