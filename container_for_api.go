package main

import (
	"io"
	"net/http"

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
	h := newErrorHandler(ctx).handleBadRequestWithError().handleConflict().handleInternalServerError()

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	payload := controller.ContainerCreatePayload{
		Name:        ctx.Name,
		Image:       ctx.Image,
		Command:     ctx.Command,
		Entrypoint:  ctx.Entrypoint,
		Env:         ctx.Env,
		Volumes:     ctx.Volumes,
		WorkingDir:  ctx.WorkingDir,
		SSLRedirect: ctx.SslRedirect,
	}

	res, status, err := c.controllerImpl.CreateWithContext(ctx, uid, payload)

	if isError(status) {
		if err := h.Call(status, err); err != nil {
			return err
		}

		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.OK(res)
	// ContainerForAPIController_Create: end_implement
}

// Download runs the download action.
func (c *ContainerForAPIController) Download(ctx *app.DownloadContainerForAPIContext) error {
	// ContainerForAPIController_Download: start_implement
	h := newErrorHandler(ctx).handleNotFoundWithError().handleInternalServerError()

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	headerOnly := ctx.Method == "HEAD"

	res, status, err := c.controllerImpl.DownloadWithContext(ctx, uid, ctx.ID, ctx.InternalPath, headerOnly)

	if isError(status) {
		if err := h.Call(status, err); err != nil {
			return err
		}

		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	defer res.Reader.Close()

	ctx.ResponseWriter.Header().Set("X-Docker-Container-Path-Stat", res.PathStatJSON)
	ctx.ResponseWriter.Header().Set("Content-Type", "application/octet-stream")
	ctx.ResponseWriter.WriteHeader(http.StatusOK)

	if _, err := io.Copy(ctx.ResponseWriter, res.Reader); err != nil {
		return ctx.InternalServerError(err)
	}

	return nil
	// ContainerForAPIController_Download: end_implement
}

// Exec runs the exec action.
func (c *ContainerForAPIController) Exec(ctx *app.ExecContainerForAPIContext) error {
	h := newErrorHandler(ctx).handleNotFoundWithError().handleInternalServerError()

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	handler, status, err := c.controllerImpl.ExecWithContext(ctx, uid, ctx.ID, &controller.ExecParameters{
		Command: ctx.Command,
		Tty:     ctx.Tty,
	})

	if isError(status) {
		if err := h.Call(status, err); err != nil {
			return err
		}

		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	handler.ServeHTTP(ctx.ResponseWriter, ctx.Request)
	return nil
}

// GetConfig runs the getConfig action.
func (c *ContainerForAPIController) GetConfig(ctx *app.GetConfigContainerForAPIContext) error {
	// ContainerForAPIController_GetConfig: start_implement
	h := newErrorHandler(ctx).handleNotFound().handleInternalServerError()

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	res, status, err := c.controllerImpl.GetConfig(uid, ctx.ID)

	if isError(status) {
		if err := h.Call(status, err); err != nil {
			return err
		}

		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.OK(res)
	// ContainerForAPIController_GetConfig: end_implement
}

// Inspect runs the inspect action.
func (c *ContainerForAPIController) Inspect(ctx *app.InspectContainerForAPIContext) error {
	// ContainerForAPIController_Inspect: start_implement
	h := newErrorHandler(ctx).handleNotFound().handleInternalServerError()

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	res, status, err := c.controllerImpl.InspectWithContext(ctx, uid, ctx.ID)

	if isError(status) {
		if err := h.Call(status, err); err != nil {
			return err
		}

		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.OK(res)
	// ContainerForAPIController_Inspect: end_implement
}

// List runs the list action.
func (c *ContainerForAPIController) List(ctx *app.ListContainerForAPIContext) error {
	// ContainerForAPIController_List: start_implement
	h := newErrorHandler(ctx).handleInternalServerError()

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	res, status, err := c.controllerImpl.ListWithContext(ctx, uid)

	if isError(status) {
		if err := h.Call(status, err); err != nil {
			return err
		}

		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.OK(res)
	// ContainerForAPIController_List: end_implement
}

// Logs runs the logs action.
func (c *ContainerForAPIController) Logs(ctx *app.LogsContainerForAPIContext) error {
	h := newErrorHandler(ctx).handleInternalServerError()

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	payload := controller.LogsPayload{
		Stdout:     ctx.Stdout,
		Stderr:     ctx.Stderr,
		Timestamps: ctx.Timestamps,
		Follow:     ctx.Follow,
		Tail:       ctx.Tail,
		Since:      ctx.Since,
		Until:      ctx.Until,
	}

	rc, status, err := c.controllerImpl.LogsWithContext(ctx, uid, ctx.ID, payload)

	if isError(status) {
		if err := h.Call(status, err); err != nil {
			return err
		}

		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	c.LogsWSHandler(ctx, rc).ServeHTTP(ctx.ResponseWriter, ctx.Request)

	return nil
}

// LogsWSHandler establishes a websocket connection to run the logs action.
func (c *ContainerForAPIController) LogsWSHandler(ctx *app.LogsContainerForAPIContext, rc io.ReadCloser) websocket.Handler {
	return func(ws *websocket.Conn) {
		// ContainerForAPIController_Logs: start_implement

		defer rc.Close()
		io.Copy(ws, rc)

		// ContainerForAPIController_Logs: end_implement
	}
}

// Remove runs the remove action.
func (c *ContainerForAPIController) Remove(ctx *app.RemoveContainerForAPIContext) error {
	// ContainerForAPIController_Remove: start_implement
	h := newErrorHandler(ctx).handleNotFound().handleRunningContainer().handleInternalServerError()

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	status, err := c.controllerImpl.Remove(uid, ctx.ID, ctx.Force)

	if isError(status) {
		if err := h.Call(status, err); err != nil {
			return err
		}

		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.NoContent()
	// ContainerForAPIController_Remove: end_implement
}

// SetConfig runs the setConfig action.
func (c *ContainerForAPIController) SetConfig(ctx *app.SetConfigContainerForAPIContext) error {
	// ContainerForAPIController_SetConfig: start_implement
	h := newErrorHandler(ctx).handleNotFound().handleInternalServerError()

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	status, err := c.controllerImpl.SetConfig(uid, ctx.ID, ctx.Payload)

	if isError(status) {
		if err := h.Call(status, err); err != nil {
			return err
		}

		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.NoContent()
	// ContainerForAPIController_SetConfig: end_implement
}

// Start runs the start action.
func (c *ContainerForAPIController) Start(ctx *app.StartContainerForAPIContext) error {
	// ContainerForAPIController_Start: start_implement
	h := newErrorHandler(ctx).handleNotFound().handleInternalServerError()

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	status, err := c.controllerImpl.Start(uid, ctx.ID)

	if err != nil {
		if err := h.Call(status, err); err != nil {
			return err
		}

		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.NoContent()
	// ContainerForAPIController_Start: end_implement
}

// Stop runs the stop action.
func (c *ContainerForAPIController) Stop(ctx *app.StopContainerForAPIContext) error {
	// ContainerForAPIController_Stop: start_implement
	h := newErrorHandler(ctx).handleNotFound().handleInternalServerError()

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	status, err := c.controllerImpl.Stop(uid, ctx.ID)

	if err != nil {
		if err := h.Call(status, err); err != nil {
			return err
		}

		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.NoContent()
	// ContainerForAPIController_Stop: end_implement
}

// Upload runs the upload action.
func (c *ContainerForAPIController) Upload(ctx *app.UploadContainerForAPIContext) error {
	// ContainerForAPIController_Upload: start_implement
	h := newErrorHandler(ctx).handleNotFoundWithError().handleBadRequestWithError().handleRequestEntityTooLarge().handleInternalServerError()

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	status, err := c.controllerImpl.UploadWithContext(ctx, uid, ctx.ID, ctx.Payload)

	if err != nil {
		if err := h.Call(status, err); err != nil {
			return err
		}

		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.NoContent()
	// ContainerForAPIController_Upload: end_implement
}
