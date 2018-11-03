package main

import (
	"io"

	"github.com/goadesign/goa"
	"github.com/modoki-paas/modoki/app"
	"github.com/modoki-paas/modoki/controller"
	"golang.org/x/net/websocket"
)

// ContainerForFrontendController implements the containerForFrontend resource.
type ContainerForFrontendController struct {
	*goa.Controller
	controllerImpl *controller.ContainerControllerImpl
}

// NewContainerForFrontendController creates a containerForFrontend controller.
func NewContainerForFrontendController(service *goa.Service) *ContainerForFrontendController {
	return &ContainerForFrontendController{Controller: service.NewController("ContainerForFrontendController")}
}

// Create runs the create action.
func (c *ContainerForFrontendController) Create(ctx *app.CreateContainerForFrontendContext) error {
	// ContainerForFrontendController_Create: start_implement

	// Put your logic here

	res := &app.GoaContainerCreateResults{}
	return ctx.OK(res)
	// ContainerForFrontendController_Create: end_implement
}

// Download runs the download action.
func (c *ContainerForFrontendController) Download(ctx *app.DownloadContainerForFrontendContext) error {
	// ContainerForFrontendController_Download: start_implement

	// Put your logic here

	return nil
	// ContainerForFrontendController_Download: end_implement
}

// Exec runs the exec action.
func (c *ContainerForFrontendController) Exec(ctx *app.ExecContainerForFrontendContext) error {
	c.ExecWSHandler(ctx).ServeHTTP(ctx.ResponseWriter, ctx.Request)
	return nil
}

// ExecWSHandler establishes a websocket connection to run the exec action.
func (c *ContainerForFrontendController) ExecWSHandler(ctx *app.ExecContainerForFrontendContext) websocket.Handler {
	return func(ws *websocket.Conn) {
		// ContainerForFrontendController_Exec: start_implement

		// Put your logic here

		ws.Write([]byte("exec containerForFrontend"))
		// Dummy echo websocket server
		io.Copy(ws, ws)
		// ContainerForFrontendController_Exec: end_implement
	}
} // GetConfig runs the getConfig action.
func (c *ContainerForFrontendController) GetConfig(ctx *app.GetConfigContainerForFrontendContext) error {
	// ContainerForFrontendController_GetConfig: start_implement

	// Put your logic here

	res := &app.GoaContainerConfig{}
	return ctx.OK(res)
	// ContainerForFrontendController_GetConfig: end_implement
}

// Inspect runs the inspect action.
func (c *ContainerForFrontendController) Inspect(ctx *app.InspectContainerForFrontendContext) error {
	// ContainerForFrontendController_Inspect: start_implement

	// Put your logic here

	res := &app.GoaContainerInspect{}
	return ctx.OK(res)
	// ContainerForFrontendController_Inspect: end_implement
}

// List runs the list action.
func (c *ContainerForFrontendController) List(ctx *app.ListContainerForFrontendContext) error {
	// ContainerForFrontendController_List: start_implement

	// Put your logic here

	res := app.GoaContainerListEachCollection{}
	return ctx.OK(res)
	// ContainerForFrontendController_List: end_implement
}

// Logs runs the logs action.
func (c *ContainerForFrontendController) Logs(ctx *app.LogsContainerForFrontendContext) error {
	c.LogsWSHandler(ctx).ServeHTTP(ctx.ResponseWriter, ctx.Request)
	return nil
}

// LogsWSHandler establishes a websocket connection to run the logs action.
func (c *ContainerForFrontendController) LogsWSHandler(ctx *app.LogsContainerForFrontendContext) websocket.Handler {
	return func(ws *websocket.Conn) {
		// ContainerForFrontendController_Logs: start_implement

		// Put your logic here

		ws.Write([]byte("logs containerForFrontend"))
		// Dummy echo websocket server
		io.Copy(ws, ws)
		// ContainerForFrontendController_Logs: end_implement
	}
} // Remove runs the remove action.
func (c *ContainerForFrontendController) Remove(ctx *app.RemoveContainerForFrontendContext) error {
	// ContainerForFrontendController_Remove: start_implement

	// Put your logic here

	return nil
	// ContainerForFrontendController_Remove: end_implement
}

// SetConfig runs the setConfig action.
func (c *ContainerForFrontendController) SetConfig(ctx *app.SetConfigContainerForFrontendContext) error {
	// ContainerForFrontendController_SetConfig: start_implement

	// Put your logic here

	return nil
	// ContainerForFrontendController_SetConfig: end_implement
}

// Start runs the start action.
func (c *ContainerForFrontendController) Start(ctx *app.StartContainerForFrontendContext) error {
	// ContainerForFrontendController_Start: start_implement

	// Put your logic here

	return nil
	// ContainerForFrontendController_Start: end_implement
}

// Stop runs the stop action.
func (c *ContainerForFrontendController) Stop(ctx *app.StopContainerForFrontendContext) error {
	// ContainerForFrontendController_Stop: start_implement

	// Put your logic here

	return nil
	// ContainerForFrontendController_Stop: end_implement
}

// Upload runs the upload action.
func (c *ContainerForFrontendController) Upload(ctx *app.UploadContainerForFrontendContext) error {
	// ContainerForFrontendController_Upload: start_implement

	// Put your logic here

	return nil
	// ContainerForFrontendController_Upload: end_implement
}
