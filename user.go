package main

import (
	"fmt"

	"github.com/docker/libkv/store"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"

	"github.com/cs3238-tsuzu/modoki/app"
	"github.com/goadesign/goa"
)

// UserController implements the user resource.
type UserController struct {
	*goa.Controller
	*UserControllerUtil
}

// NewUserController creates a user controller.
func NewUserController(service *goa.Service) *UserController {
	return &UserController{Controller: service.NewController("UserController")}
}

// GetConfig runs the getConfig action.
func (c *UserController) GetConfig(ctx *app.GetConfigUserContext) error {
	// UserController_GetConfig: start_implement

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	var config app.GoaUserConfig

	if p, err := c.Consul.Client.Get(fmt.Sprint(defaultShellKVFormat, uid)); err != nil {
		if err != store.ErrKeyNotFound {
			return ctx.InternalServerError(goa.ErrInternal(errors.Wrap(err, "consul error")))
		}
	} else {
		config.DefaultShell = string(p.Value)
	}

	var keys []*app.GoaUserAuthorizedkey

	if err := c.DB.Select(&keys, "SELECT `key`, label FROM authorizedKeys WHERE uid=?", uid); err != nil {
		return ctx.InternalServerError(goa.ErrInternal(errors.Wrap(err, "DB error")))
	}

	config.AuthorizedKeys = app.GoaUserAuthorizedkeyCollection(keys)

	return ctx.OK(&config)

	// UserController_GetConfig: end_implement
}

// AddAuthorizedKeys runs the addAuthorizedKeys action.
func (c *UserController) AddAuthorizedKeys(ctx *app.AddAuthorizedKeysUserContext) error {
	// UserController_AddAuthorizedKeys: start_implement

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	_, _, _, _, err = ssh.ParseAuthorizedKey([]byte(ctx.Payload.Key))
	if err != nil {
		return ctx.BadRequest()
	}

	_, err = c.DB.Exec("INSERT INTO authorizedKeys (label, `key`, uid) VALUES (?, ?, ?)", ctx.Payload.Label, ctx.Payload.Key, uid)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(errors.Wrap(err, "DB error")))
	}

	return ctx.NoContent()
	// UserController_AddAuthorizedKeys: end_implement
}

// ListAuthorizedKeys runs the listAuthorizedKeys action.
func (c *UserController) ListAuthorizedKeys(ctx *app.ListAuthorizedKeysUserContext) error {
	// UserController_ListAuthorizedKeys: start_implement

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	var keys []*app.GoaUserAuthorizedkey

	if err := c.DB.Select(&keys, "SELECT `key`, label FROM authorizedKeys WHERE uid=?", uid); err != nil {
		return ctx.InternalServerError(goa.ErrInternal(errors.Wrap(err, "DB error")))
	}

	res := app.GoaUserAuthorizedkeyCollection(keys)
	return ctx.OK(res)
	// UserController_ListAuthorizedKeys: end_implement
}

// RemoveAuthorizedKeys runs the removeAuthorizedKeys action.
func (c *UserController) RemoveAuthorizedKeys(ctx *app.RemoveAuthorizedKeysUserContext) error {
	// UserController_RemoveAuthorizedKeys: start_implement

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	res, err := c.DB.Exec("DELETE FROM authorizedKeys WHERE uid=? AND label=?", uid, ctx.Label)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(errors.Wrap(err, "DB error")))
	}

	if r, _ := res.RowsAffected(); r == 0 {
		return ctx.NotFound()
	}

	return ctx.NoContent()
	// UserController_RemoveAuthorizedKeys: end_implement
}

// SetAuthorizedKeys runs the setAuthorizedKeys action.
func (c *UserController) SetAuthorizedKeys(ctx *app.SetAuthorizedKeysUserContext) error {
	// UserController_SetAuthorizedKeys: start_implement

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	_, err = c.DB.Exec("DELETE FROM authorizedKeys WHERE uid=?", uid)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(errors.Wrap(err, "DB error")))
	}

	for i := range ctx.Payload {
		if ctx.Payload[i] != nil {
			_, err := c.DB.Exec("INSERT INTO authorizedKeys (label, `key`, uid) VALUES (?, ?, ?)", ctx.Payload[i].Label, ctx.Payload[i].Key, uid)

			if err != nil {
				return ctx.InternalServerError(goa.ErrInternal(errors.Wrap(err, "DB error")))
			}
		}
	}

	return ctx.NoContent()
	// UserController_SetAuthorizedKeys: end_implement
}

// GetDefaultShell runs the getDefaultShell action.
func (c *UserController) GetDefaultShell(ctx *app.GetDefaultShellUserContext) error {
	// UserController_GetDefaultShell: start_implement

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	res := &app.GoaUserDefaultshell{}
	if p, err := c.Consul.Client.Get(fmt.Sprint(defaultShellKVFormat, uid)); err != nil {
		if err != store.ErrKeyNotFound {
			return ctx.InternalServerError(goa.ErrInternal(errors.Wrap(err, "consul error")))
		}
	} else {
		res.DefaultShell = string(p.Value)
	}

	return ctx.OK(res)
	// UserController_GetDefaultShell: end_implement
}

// SetDefaultShell runs the setDefaultShell action.
func (c *UserController) SetDefaultShell(ctx *app.SetDefaultShellUserContext) error {
	// UserController_SetDefaultShell: start_implement

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	if err := c.Consul.Client.Put(fmt.Sprint(defaultShellKVFormat, uid), []byte(ctx.DefaultShell), nil); err != nil {
		return ctx.InternalServerError(goa.ErrInternal(errors.Wrap(err, "consul error")))
	}

	return ctx.NoContent()
	// UserController_SetDefaultShell: end_implement
}
