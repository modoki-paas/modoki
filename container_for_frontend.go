package main

import (
	"io"
	"net/http"

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
	h := newErrorHandler(ctx).handleBadRequestWithError().handleConflict().handleInternalServerError()

	uid, err := GetUIDFromContext(ctx)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	payload := controller.ContainerCreateParameters{
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

	resp := &app.GoaContainerCreateResults{
		Endpoints: res.Endpoints,
		ID:        res.ID,
	}

	return ctx.OK(resp)
	// ContainerForFrontendController_Create: end_implement
}

// Download runs the download action.
func (c *ContainerForFrontendController) Download(ctx *app.DownloadContainerForFrontendContext) error {
	// ContainerForFrontendController_Download: start_implement
	h := newErrorHandler(ctx).handleNotFoundWithError().handleInternalServerError()

	uid, err := GetUIDFromContext(ctx)

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
	// ContainerForFrontendController_Download: end_implement
}

// Exec runs the exec action.
func (c *ContainerForFrontendController) Exec(ctx *app.ExecContainerForFrontendContext) error {
	h := newErrorHandler(ctx).handleNotFoundWithError().handleInternalServerError()

	uid, err := GetUIDFromContext(ctx)

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
func (c *ContainerForFrontendController) GetConfig(ctx *app.GetConfigContainerForFrontendContext) error {
	// ContainerForFrontendController_GetConfig: start_implement
	h := newErrorHandler(ctx).handleNotFound().handleInternalServerError()

	uid, err := GetUIDFromContext(ctx)

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

	resp := &app.GoaContainerConfig{
		DefaultShell: res.DefaultShell,
	}

	return ctx.OK(resp)
	// ContainerForFrontendController_GetConfig: end_implement
}

// Inspect runs the inspect action.
func (c *ContainerForFrontendController) Inspect(ctx *app.InspectContainerForFrontendContext) error {
	// ContainerForFrontendController_Inspect: start_implement
	h := newErrorHandler(ctx).handleNotFound().handleInternalServerError()

	uid, err := GetUIDFromContext(ctx)

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

	resp := &app.GoaContainerInspect{
		Args:    res.Args,
		Created: res.Created,
		ID:      res.ID,
		Image:   res.Image,
		ImageID: res.ImageID,
		Name:    res.Name,
		Path:    res.Path,
		RawState: &app.GoaContainerInspectRawState{
			Dead:       res.RawState.Dead,
			ExitCode:   res.RawState.ExitCode,
			FinishedAt: res.RawState.FinishedAt,
			OomKilled:  res.RawState.OomKilled,
			Paused:     res.RawState.Paused,
			Pid:        res.RawState.Pid,
			Restarting: res.RawState.Restarting,
			Running:    res.RawState.Running,
			StartedAt:  res.RawState.StartedAt,
			Status:     res.RawState.Status,
		},
		Status:  res.Status,
		Volumes: res.Volumes,
	}

	return ctx.OK(resp)
	// ContainerForFrontendController_Inspect: end_implement
}

// List runs the list action.
func (c *ContainerForFrontendController) List(ctx *app.ListContainerForFrontendContext) error {
	// ContainerForFrontendController_List: start_implement
	h := newErrorHandler(ctx).handleInternalServerError()

	uid, err := GetUIDFromContext(ctx)

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

	resp := make(app.GoaContainerListEachCollection, 0, len(res))

	for i := range res {
		resp = append(resp, &app.GoaContainerListEach{
			Command: res[i].Command,
			Created: res[i].Created,
			ID:      res[i].ID,
			Image:   res[i].Image,
			ImageID: res[i].ImageID,
			Name:    res[i].Name,
			Status:  res[i].Status,
			Volumes: res[i].Volumes,
		})
	}

	return ctx.OK(resp)
	// ContainerForFrontendController_List: end_implement
}

// Logs runs the logs action.
func (c *ContainerForFrontendController) Logs(ctx *app.LogsContainerForFrontendContext) error {
	h := newErrorHandler(ctx).handleInternalServerError()

	uid, err := GetUIDFromContext(ctx)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	payload := controller.LogsParameters{
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
func (c *ContainerForFrontendController) LogsWSHandler(ctx *app.LogsContainerForFrontendContext, rc io.ReadCloser) websocket.Handler {
	return func(ws *websocket.Conn) {
		// ContainerForFrontendController_Logs: start_implement

		defer rc.Close()
		io.Copy(ws, rc)

		// ContainerForFrontendController_Logs: end_implement
	}
} // Remove runs the remove action.
func (c *ContainerForFrontendController) Remove(ctx *app.RemoveContainerForFrontendContext) error {
	// ContainerForFrontendController_Remove: start_implement
	h := newErrorHandler(ctx).handleNotFound().handleRunningContainer().handleInternalServerError()

	uid, err := GetUIDFromContext(ctx)

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
	// ContainerForFrontendController_Remove: end_implement
}

// SetConfig runs the setConfig action.
func (c *ContainerForFrontendController) SetConfig(ctx *app.SetConfigContainerForFrontendContext) error {
	// ContainerForFrontendController_SetConfig: start_implement
	h := newErrorHandler(ctx).handleNotFound().handleInternalServerError()

	uid, err := GetUIDFromContext(ctx)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	config := &controller.ContainerConfig{
		DefaultShell: ctx.Payload.DefaultShell,
	}

	status, err := c.controllerImpl.SetConfig(uid, ctx.ID, config)

	if isError(status) {
		if err := h.Call(status, err); err != nil {
			return err
		}

		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.NoContent()
	// ContainerForFrontendController_SetConfig: end_implement
}

// Start runs the start action.
func (c *ContainerForFrontendController) Start(ctx *app.StartContainerForFrontendContext) error {
	// ContainerForFrontendController_Start: start_implement
	h := newErrorHandler(ctx).handleNotFound().handleInternalServerError()

	uid, err := GetUIDFromContext(ctx)

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
	// ContainerForFrontendController_Start: end_implement
}

// Stop runs the stop action.
func (c *ContainerForFrontendController) Stop(ctx *app.StopContainerForFrontendContext) error {
	// ContainerForFrontendController_Stop: start_implement
	h := newErrorHandler(ctx).handleNotFound().handleInternalServerError()

	uid, err := GetUIDFromContext(ctx)

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
	// ContainerForFrontendController_Stop: end_implement
}

// Upload runs the upload action.
func (c *ContainerForFrontendController) Upload(ctx *app.UploadContainerForFrontendContext) error {
	// ContainerForFrontendController_Upload: start_implement
	h := newErrorHandler(ctx).handleNotFoundWithError().handleBadRequestWithError().handleRequestEntityTooLarge().handleInternalServerError()

	uid, err := GetUIDFromContext(ctx)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	file, err := ctx.Payload.Data.Open()

	if err != nil {
		return ctx.InternalServerError(err)
	}
	defer file.Close()

	param := &controller.UploadParameters{
		AllowOverwrite: ctx.Payload.AllowOverwrite,
		CopyUIDGID:     ctx.Payload.CopyUIDGID,
		Data:           file,
		Path:           ctx.Payload.Path,
	}

	status, err := c.controllerImpl.UploadWithContext(ctx, uid, ctx.ID, param)

	if err != nil {
		if err := h.Call(status, err); err != nil {
			return err
		}

		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.NoContent()
	// ContainerForFrontendController_Upload: end_implement
}
