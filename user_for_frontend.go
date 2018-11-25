package main

import (
	"github.com/goadesign/goa"
	"github.com/modoki-paas/modoki/app"
	"github.com/modoki-paas/modoki/controller"
)

// UserForFrontendController implements the userForFrontend resource.
type UserForFrontendController struct {
	*goa.Controller
	controllerImpl *controller.UserControllerImpl
}

// NewUserForFrontendController creates a userForFrontend controller.
func NewUserForFrontendController(service *goa.Service) *UserForFrontendController {
	return &UserForFrontendController{Controller: service.NewController("UserForFrontendController")}
}

// AddAuthorizedKeys runs the addAuthorizedKeys action.
func (c *UserForFrontendController) AddAuthorizedKeys(ctx *app.AddAuthorizedKeysUserForFrontendContext) error {
	// UserForFrontendController_AddAuthorizedKeys: start_implement
	h := newErrorHandler(ctx).handleNoContent().handleInternalServerError().handleBadRequest()

	uid, err := GetUIDFromContext(ctx)

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
	// UserForFrontendController_AddAuthorizedKeys: end_implement
}

// GetAPIKey runs the getAPIKey action.
func (c *UserForFrontendController) GetAPIKey(ctx *app.GetAPIKeyUserForFrontendContext) error {
	// UserForFrontendController_GetAPIKey: start_implement
	h := newErrorHandler(ctx).handleInternalServerError()

	uid, err := GetUIDFromContext(ctx)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	apiKey, status, err := c.controllerImpl.GetAPIKey(uid)

	if isError(status) {
		if err := h.Call(status, err); err != nil {
			return err
		}

		return ctx.InternalServerError(err)
	}

	return ctx.OK(apiKey)
	// UserForFrontendController_GetAPIKey: end_implement
}

// GetConfig runs the getConfig action.
func (c *UserForFrontendController) GetConfig(ctx *app.GetConfigUserForFrontendContext) error {
	// UserForFrontendController_GetConfig: start_implement
	h := newErrorHandler(ctx).handleInternalServerError()

	uid, err := GetUIDFromContext(ctx)

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
	// UserForFrontendController_GetConfig: end_implement
}

// GetDefaultShell runs the getDefaultShell action.
func (c *UserForFrontendController) GetDefaultShell(ctx *app.GetDefaultShellUserForFrontendContext) error {
	// UserForFrontendController_GetDefaultShell: start_implement
	h := newErrorHandler(ctx).handleInternalServerError()

	uid, err := GetUIDFromContext(ctx)

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
	// UserForFrontendController_GetDefaultShell: end_implement
}

// ListAuthorizedKeys runs the listAuthorizedKeys action.
func (c *UserForFrontendController) ListAuthorizedKeys(ctx *app.ListAuthorizedKeysUserForFrontendContext) error {
	// UserForFrontendController_ListAuthorizedKeys: start_implement
	h := newErrorHandler(ctx).handleInternalServerError().handleNotFound()

	uid, err := GetUIDFromContext(ctx)

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
	// UserForFrontendController_ListAuthorizedKeys: end_implement
}

// ReissueAPIKey runs the reissueAPIKey action.
func (c *UserForFrontendController) ReissueAPIKey(ctx *app.ReissueAPIKeyUserForFrontendContext) error {
	// UserForFrontendController_ReissueAPIKey: start_implement
	h := newErrorHandler(ctx).handleInternalServerError()

	uid, err := GetUIDFromContext(ctx)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	apiKey, status, err := c.controllerImpl.ReissueAPIKey(uid)

	if isError(status) {
		if err := h.Call(status, err); err != nil {
			return err
		}

		return ctx.InternalServerError(err)
	}

	return ctx.OK(apiKey)
	// UserForFrontendController_ReissueAPIKey: end_implement
}

// RemoveAuthorizedKeys runs the removeAuthorizedKeys action.
func (c *UserForFrontendController) RemoveAuthorizedKeys(ctx *app.RemoveAuthorizedKeysUserForFrontendContext) error {
	// UserForFrontendController_RemoveAuthorizedKeys: start_implement
	h := newErrorHandler(ctx).handleInternalServerError().handleNotFound()

	uid, err := GetUIDFromContext(ctx)

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
	// UserForFrontendController_RemoveAuthorizedKeys: end_implement
}

// SetAuthorizedKeys runs the setAuthorizedKeys action.
func (c *UserForFrontendController) SetAuthorizedKeys(ctx *app.SetAuthorizedKeysUserForFrontendContext) error {
	// UserForFrontendController_SetAuthorizedKeys: start_implement
	h := newErrorHandler(ctx).handleInternalServerError().handleNotFound()

	uid, err := GetUIDFromContext(ctx)

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
	// UserForFrontendController_SetAuthorizedKeys: end_implement
}

// SetDefaultShell runs the setDefaultShell action.
func (c *UserForFrontendController) SetDefaultShell(ctx *app.SetDefaultShellUserForFrontendContext) error {
	// UserForFrontendController_SetDefaultShell: start_implement
	h := newErrorHandler(ctx).handleInternalServerError()

	uid, err := GetUIDFromContext(ctx)

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
	// UserForFrontendController_SetDefaultShell: end_implement
}
