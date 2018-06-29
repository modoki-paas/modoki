package main

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/cs3238-tsuzu/modoki/app"
	"github.com/cs3238-tsuzu/modoki/consul_traefik"
	"github.com/docker/docker/client"
	"github.com/goadesign/goa"
	"github.com/jmoiron/sqlx"
)

// UserController implements the user resource.
type UserController struct {
	*goa.Controller

	DB           *sqlx.DB
	DockerClient *client.Client
	Consul       *consulTraefik.Client
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
		return ctx.InternalServerError(goa.ErrInternal(errors.Wrap(err, "consul error")))
	} else {
		config.DefaultShell = string(p.Value)
	}

	var keys []*app.GoaUserAuthorizedkey

	if err := c.DB.Select(&keys, "SELECT * FROM authorizedKeys WHERE uid=?", uid); err != nil {
		return ctx.InternalServerError(goa.ErrInternal(errors.Wrap(err, "DB error")))
	}

	config.AuthorizedKeys = app.GoaUserAuthorizedkeyCollection(keys)

	return ctx.OK(&config)

	// UserController_GetConfig: end_implement
}

// SetConfig runs the setConfig action.
func (c *UserController) SetConfig(ctx *app.SetConfigUserContext) error {
	// UserController_SetConfig: start_implement

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	if ctx.Payload.DefaultShell != nil {
		if err := c.Consul.Client.Put(fmt.Sprint(defaultShellKVFormat, uid), []byte(*ctx.Payload.DefaultShell), nil); err != nil {
			return ctx.InternalServerError(goa.ErrInternal(errors.Wrap(err, "consul error")))
		}
	}

	if ctx.Payload.AuthorizedKeys != nil {
		for i := range ctx.Payload.AuthorizedKeys {
			if ctx.Payload.AuthorizedKeys[i] != nil {
				_, err := c.DB.Exec("INSERT INTO authorizedKeys (label, `key`, uid) VALUES (?, ?, ?)", ctx.Payload.AuthorizedKeys[i].Label, ctx.Payload.AuthorizedKeys[i].Key, uid)

				if err != nil {
					return ctx.InternalServerError(goa.ErrInternal(errors.Wrap(err, "DB error")))
				}
			}
		}
	}

	return ctx.NoContent()
	// UserController_SetConfig: end_implement
}

// AddAuthorizedKeys runs the addAuthorizedKeys action.
func (c *UserController) AddAuthorizedKeys(ctx *app.AddAuthorizedKeysUserContext) error {
	// UserController_AddAuthorizedKeys: start_implement

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	_, err = c.DB.Exec("INSERT INTO authorizedKeys (label, `key`, uid) VALUES (?, ?, uid)", ctx.Payload.Label, ctx.Payload.Key, uid)

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

	if err := c.DB.Select(&keys, "SELECT * FROM authorizedKeys WHERE uid=?", uid); err != nil {
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
