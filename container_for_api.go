package main

import (
	"io"

	"github.com/goadesign/goa"
	"github.com/modoki-paas/modoki/app"
	"github.com/modoki-paas/modoki/controller"
	"golang.org/x/net/websocket"
)

// ContainerForAPIController implements the containerForApi resource.
type ContainerForAPIController struct {
	*goa.Controller
	controllerImpl *controller.ContainerControllerImpl
}

// NewContainerForAPIController creates a containerForApi controller.
func NewContainerForAPIController(service *goa.Service) *ContainerForAPIController {
	return &ContainerForAPIController{Controller: service.NewController("ContainerForAPIController")}
}

// Create runs the create action.
func (c *ContainerForAPIController) Create(ctx *app.CreateContainerForAPIContext) error {
	// ContainerForAPIController_Create: start_implement

	// Put your logic here

	res := &app.GoaContainerCreateResults{}
	return ctx.OK(res)
	// ContainerForAPIController_Create: end_implement
}

// Download runs the download action.
func (c *ContainerForAPIController) Download(ctx *app.DownloadContainerForAPIContext) error {
	// ContainerForAPIController_Download: start_implement

	// Put your logic here

	return nil
	// ContainerForAPIController_Download: end_implement
}

// Exec runs the exec action.
func (c *ContainerForAPIController) Exec(ctx *app.ExecContainerForAPIContext) error {
	c.ExecWSHandler(ctx).ServeHTTP(ctx.ResponseWriter, ctx.Request)
	return nil
}

// ExecWSHandler establishes a websocket connection to run the exec action.
func (c *ContainerForAPIController) ExecWSHandler(ctx *app.ExecContainerForAPIContext) websocket.Handler {
	return func(ws *websocket.Conn) {
		// ContainerForAPIController_Exec: start_implement

		// Put your logic here

		ws.Write([]byte("exec containerForApi"))
		// Dummy echo websocket server
		io.Copy(ws, ws)
		// ContainerForAPIController_Exec: end_implement
	}
} // GetConfig runs the getConfig action.
func (c *ContainerForAPIController) GetConfig(ctx *app.GetConfigContainerForAPIContext) error {
	// ContainerForAPIController_GetConfig: start_implement

	// Put your logic here

	res := &app.GoaContainerConfig{}
	return ctx.OK(res)
	// ContainerForAPIController_GetConfig: end_implement
}

// Inspect runs the inspect action.
func (c *ContainerForAPIController) Inspect(ctx *app.InspectContainerForAPIContext) error {
	// ContainerForAPIController_Inspect: start_implement

	// Put your logic here

	res := &app.GoaContainerInspect{}
	return ctx.OK(res)
	// ContainerForAPIController_Inspect: end_implement
}

// List runs the list action.
func (c *ContainerForAPIController) List(ctx *app.ListContainerForAPIContext) error {
	// ContainerForAPIController_List: start_implement

	// Put your logic here

	res := app.GoaContainerListEachCollection{}
	return ctx.OK(res)
	// ContainerForAPIController_List: end_implement
}

// Logs runs the logs action.
func (c *ContainerForAPIController) Logs(ctx *app.LogsContainerForAPIContext) error {
	c.LogsWSHandler(ctx).ServeHTTP(ctx.ResponseWriter, ctx.Request)
	return nil
}

// LogsWSHandler establishes a websocket connection to run the logs action.
func (c *ContainerForAPIController) LogsWSHandler(ctx *app.LogsContainerForAPIContext) websocket.Handler {
	return func(ws *websocket.Conn) {
		// ContainerForAPIController_Logs: start_implement

		// Put your logic here

		ws.Write([]byte("logs containerForApi"))
		// Dummy echo websocket server
		io.Copy(ws, ws)
		// ContainerForAPIController_Logs: end_implement
	}
} // Remove runs the remove action.
func (c *ContainerForAPIController) Remove(ctx *app.RemoveContainerForAPIContext) error {
	// ContainerForAPIController_Remove: start_implement

	// Put your logic here

	return nil
	// ContainerForAPIController_Remove: end_implement
}

// SetConfig runs the setConfig action.
func (c *ContainerForAPIController) SetConfig(ctx *app.SetConfigContainerForAPIContext) error {
	// ContainerForAPIController_SetConfig: start_implement

	// Put your logic here

	return nil
	// ContainerForAPIController_SetConfig: end_implement
}

// Start runs the start action.
func (c *ContainerForAPIController) Start(ctx *app.StartContainerForAPIContext) error {
	// ContainerForAPIController_Start: start_implement

	// Put your logic here

	return nil
	// ContainerForAPIController_Start: end_implement
}

// Stop runs the stop action.
func (c *ContainerForAPIController) Stop(ctx *app.StopContainerForAPIContext) error {
	// ContainerForAPIController_Stop: start_implement

	// Put your logic here

	return nil
	// ContainerForAPIController_Stop: end_implement
}

// Upload runs the upload action.
func (c *ContainerForAPIController) Upload(ctx *app.UploadContainerForAPIContext) error {
	// ContainerForAPIController_Upload: start_implement

	// Put your logic here

	return nil
	// ContainerForAPIController_Upload: end_implement
}
