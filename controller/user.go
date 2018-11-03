package controller

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/docker/libkv/store"
	"github.com/modoki-paas/modoki/app"
	"github.com/modoki-paas/modoki/const"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

type UserControllerImpl struct {
	*UserControllerUtil
}

func NewUserControllerImpl(util *UserControllerUtil) *UserControllerImpl {
	return &UserControllerImpl{
		UserControllerUtil: util,
	}
}

// GetConfig runs the getConfig action.
func (c *UserControllerImpl) GetConfig(uid string) (*app.GoaUserConfig, int, error) {
	var config app.GoaUserConfig

	if p, err := c.Consul.Client.Get(fmt.Sprint(constants.DefaultShellKVFormat, uid)); err != nil {
		if err != store.ErrKeyNotFound {
			return nil, http.StatusInternalServerError, errors.Wrap(err, "consul error")
		}
	} else {
		config.DefaultShell = string(p.Value)
	}

	var keys []*app.GoaUserAuthorizedkey

	if err := c.DB.Select(&keys, "SELECT `key`, label FROM authorizedKeys WHERE uid=?", uid); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "DB error")
	}

	config.AuthorizedKeys = app.GoaUserAuthorizedkeyCollection(keys)

	return &config, http.StatusOK, nil
}

// AddAuthorizedKeys runs the addAuthorizedKeys action.
func (c *UserControllerImpl) AddAuthorizedKeys(uid, label, key string) (int, error) {
	_, _, _, _, err := ssh.ParseAuthorizedKey([]byte(key))
	if err != nil {
		return http.StatusBadRequest, errors.New("invalid key format")
	}

	_, err = c.DB.Exec("INSERT INTO authorizedKeys (label, `key`, uid) VALUES (?, ?, ?)", label, key, uid)

	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "DB error")
	}

	return http.StatusNoContent, nil
}

// ListAuthorizedKeys runs the listAuthorizedKeys action.
func (c *UserControllerImpl) ListAuthorizedKeys(uid string) (app.GoaUserAuthorizedkeyCollection, int, error) {
	var keys []*app.GoaUserAuthorizedkey

	if err := c.DB.Select(&keys, "SELECT `key`, label FROM authorizedKeys WHERE uid=?", uid); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "DB error")
	}

	return app.GoaUserAuthorizedkeyCollection(keys), http.StatusOK, nil
}

// RemoveAuthorizedKeys runs the removeAuthorizedKeys action.
func (c *UserControllerImpl) RemoveAuthorizedKeys(uid, label string) (int, error) {
	res, err := c.DB.Exec("DELETE FROM authorizedKeys WHERE uid=? AND label=?", uid, label)

	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "DB error")
	}

	if r, _ := res.RowsAffected(); r == 0 {
		return http.StatusNotFound, errors.New("key not found")
	}

	return http.StatusNoContent, nil
}

// SetAuthorizedKeys runs the setAuthorizedKeys action.
func (c *UserControllerImpl) SetAuthorizedKeys(uid string, keys []SSHKey) (int, error) {
	_, err := c.DB.Exec("DELETE FROM authorizedKeys WHERE uid=?", uid)

	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "DB error")
	}

	for i := range keys {
		_, err := c.DB.Exec("INSERT INTO authorizedKeys (label, `key`, uid) VALUES (?, ?, ?)", keys[i].Label, keys[i].Key, uid)

		if err != nil {
			return http.StatusInternalServerError, errors.Wrap(err, "DB error")
		}
	}

	return http.StatusNoContent, nil
}

// GetDefaultShell runs the getDefaultShell action.
func (c *UserControllerImpl) GetDefaultShell(uid string) (*app.GoaUserDefaultshell, int, error) {
	res := &app.GoaUserDefaultshell{}
	if p, err := c.Consul.Client.Get(fmt.Sprint(constants.DefaultShellKVFormat, uid)); err != nil {
		if err != store.ErrKeyNotFound {
			return nil, http.StatusInternalServerError, errors.Wrap(err, "consul error")
		}
	} else {
		res.DefaultShell = string(p.Value)
	}

	return res, http.StatusOK, nil
}

// SetDefaultShell runs the setDefaultShell action.
func (c *UserControllerImpl) SetDefaultShell(uid, defaultShell string) (int, error) {
	if err := c.Consul.Client.Put(fmt.Sprint(constants.DefaultShellKVFormat, uid), []byte(defaultShell), nil); err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "consul error")
	}

	return http.StatusNoContent, nil
}

// GetAPIKey runs the getAPIKey action.
func (c *UserControllerImpl) GetAPIKey(uid string) (*app.GoaUserApikey, int, error) {
	tx, err := c.DB.Beginx()

	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "db error")
	}

	var apiKey string
	rollback, err := func() (bool, error) {
		var userApiKey UserAPIKey
		err := tx.Get(&userApiKey, "SELECT * FROM apiKeys WHERE uid=?", uid)

		if err != nil && err != sql.ErrNoRows {
			return false, err
		}

		if err != sql.ErrNoRows {
			apiKey = userApiKey.ApiKey

			return false, nil
		}

		apiKey = generateAPIKey()
		_, err = tx.Exec("INSERT INTO apiKeys (uid, apiKey) VALUES (?, ?)", uid, apiKey)

		if err != nil {
			return false, err
		}

		return false, nil
	}()

	if rollback {
		if err := tx.Rollback(); err != nil {
			return nil, http.StatusInternalServerError, err
		}
	} else {
		if err := tx.Commit(); err != nil {
			return nil, http.StatusInternalServerError, err
		}
	}

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	res := &app.GoaUserApikey{
		Key: apiKey,
	}
	return res, http.StatusOK, nil
}

// ReissueAPIKey runs the reissueAPIKey action.
func (c *UserControllerImpl) ReissueAPIKey(uid string) (*app.GoaUserApikey, int, error) {
	tx, err := c.DB.Beginx()

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	apiKey := generateAPIKey()

	rollback, err := func() (bool, error) {
		res, err := tx.Exec("UPDATE apiKeys SET apiKey WHERE uid=?", apiKey, uid)

		if err != nil {
			return false, err
		}

		if affected, err := res.RowsAffected(); err != nil {
			return false, err
		} else if affected != 0 {
			return false, nil
		}

		_, err = tx.Exec("INSERT INTO apiKeys (uid, apiKey) VALUES (?, ?)", uid, apiKey)

		if err != nil {
			return false, err
		}

		return false, nil
	}()

	if rollback {
		if err := tx.Rollback(); err != nil {
			return nil, http.StatusInternalServerError, err
		}
	} else {
		if err := tx.Commit(); err != nil {
			return nil, http.StatusInternalServerError, err
		}
	}

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	res := &app.GoaUserApikey{
		Key: apiKey,
	}
	return res, http.StatusOK, nil
}
