package main

import (
	"github.com/goadesign/goa"
	"github.com/modoki-paas/modoki/app"
	"github.com/modoki-paas/modoki/controller"
)

// UserForAPIController implements the userForApi resource.
type UserForAPIController struct {
	*goa.Controller
	controllerImpl *controller.UserControllerImpl
}

// NewUserForAPIController creates a userForApi controller.
func NewUserForAPIController(service *goa.Service) *UserForAPIController {
	return &UserForAPIController{Controller: service.NewController("UserForAPIController")}
}

// AddAuthorizedKeys runs the addAuthorizedKeys action.
func (c *UserForAPIController) AddAuthorizedKeys(ctx *app.AddAuthorizedKeysUserForAPIContext) error {
	// UserForAPIController_AddAuthorizedKeys: start_implement
	h := newErrorHandler(ctx).handleNoContent().handleInternalServerError().handleBadRequest()

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	status, err := c.controllerImpl.AddAuthorizedKeys(uid, ctx.Payload.Label, ctx.Payload.Key)

	if isError(status) {

		if err := h.Call(status, err); err != nil {
			return err
		}

		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return nil
	// UserForAPIController_AddAuthorizedKeys: end_implement
}

// GetConfig runs the getConfig action.
func (c *UserForAPIController) GetConfig(ctx *app.GetConfigUserForAPIContext) error {
	// UserForAPIController_GetConfig: start_implement
	h := newErrorHandler(ctx).handleInternalServerError()

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	config, status, err := c.controllerImpl.GetConfig(uid)

	if isError(status) {
		if err := h.Call(status, err); err != nil {
			return err
		}

		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.OK(config)
	// UserForAPIController_GetConfig: end_implement
}

// GetDefaultShell runs the getDefaultShell action.
func (c *UserForAPIController) GetDefaultShell(ctx *app.GetDefaultShellUserForAPIContext) error {
	// UserForAPIController_GetDefaultShell: start_implement

	h := newErrorHandler(ctx).handleInternalServerError()

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	defaultShell, status, err := c.controllerImpl.GetDefaultShell(uid)

	if isError(status) {
		if err := h.Call(status, err); err != nil {
			return err
		}

		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.OK(defaultShell)
	// UserForAPIController_GetDefaultShell: end_implement
}

// ListAuthorizedKeys runs the listAuthorizedKeys action.
func (c *UserForAPIController) ListAuthorizedKeys(ctx *app.ListAuthorizedKeysUserForAPIContext) error {
	// UserForAPIController_ListAuthorizedKeys: start_implement
	h := newErrorHandler(ctx).handleInternalServerError().handleNotFound()

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	config, status, err := c.controllerImpl.ListAuthorizedKeys(uid)

	if isError(status) {
		if err := h.Call(status, err); err != nil {
			return err
		}

		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.OK(config)
	// UserForAPIController_ListAuthorizedKeys: end_implement
}

// RemoveAuthorizedKeys runs the removeAuthorizedKeys action.
func (c *UserForAPIController) RemoveAuthorizedKeys(ctx *app.RemoveAuthorizedKeysUserForAPIContext) error {
	// UserForAPIController_RemoveAuthorizedKeys: start_implement
	h := newErrorHandler(ctx).handleInternalServerError().handleNotFound()

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	status, err := c.controllerImpl.RemoveAuthorizedKeys(uid, ctx.Label)

	if isError(status) {
		if err := h.Call(status, err); err != nil {
			return err
		}

		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.NoContent()
	// UserForAPIController_RemoveAuthorizedKeys: end_implement
}

// SetAuthorizedKeys runs the setAuthorizedKeys action.
func (c *UserForAPIController) SetAuthorizedKeys(ctx *app.SetAuthorizedKeysUserForAPIContext) error {
	// UserForAPIController_SetAuthorizedKeys: start_implement
	h := newErrorHandler(ctx).handleInternalServerError().handleNotFound()

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	keys := make([]controller.SSHKey, 0, len(ctx.Payload))

	for i := range ctx.Payload {
		keys = append(keys, controller.SSHKey{
			Key:   ctx.Payload[i].Key,
			Label: ctx.Payload[i].Label,
		})
	}

	status, err := c.controllerImpl.SetAuthorizedKeys(uid, keys)

	if isError(status) {
		if err := h.Call(status, err); err != nil {
			return err
		}

		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.NoContent()
	// UserForAPIController_SetAuthorizedKeys: end_implement
}

// SetDefaultShell runs the setDefaultShell action.
func (c *UserForAPIController) SetDefaultShell(ctx *app.SetDefaultShellUserForAPIContext) error {
	// UserForAPIController_SetDefaultShell: start_implement
	h := newErrorHandler(ctx).handleInternalServerError()

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	status, err := c.controllerImpl.SetDefaultShell(uid, ctx.DefaultShell)

	if isError(status) {
		if err := h.Call(status, err); err != nil {
			return err
		}

		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.NoContent()
	// UserForAPIController_SetDefaultShell: end_implement
}
