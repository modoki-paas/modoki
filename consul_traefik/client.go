package consulTraefik

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/docker/libkv/store/consul"

	"github.com/docker/libkv/store"
)

type Client struct {
	Prefix string
	Client store.Store
}

func NewClient(prefix, addr string) (*Client, error) {
	store, err := consul.New([]string{addr}, nil)

	if err != nil {
		return nil, err
	}

	return &Client{
		Prefix: prefix,
		Client: store,
	}, nil
}

func (c *Client) NewBackend(backend, serverName, addr string) error {
	err := c.Client.Put(
		c.Prefix+"/backends/"+backend+"/servers/"+serverName+"/url",
		[]byte(addr),
		nil,
	)

	return err
}

func (c *Client) BackupBackend(backend, serverName string) ([]byte, error) {
	keyPrefix := c.Prefix + "/backends/" + backend + "/servers/" + serverName + "/"
	pairs, err := c.Client.List(keyPrefix)

	if err != nil {
		return nil, err
	}

	m := make(map[string][]byte)
	for i := range pairs {
		m[pairs[i].Key[len(keyPrefix):]] = pairs[i].Value
	}
	b, _ := json.Marshal(m)

	return b, nil
}

func (c *Client) RestoreBackup(backend, serverName string, backup []byte) error {
	keyPrefix := c.Prefix + "/backends/" + backend + "/servers/" + serverName + "/"
	var m map[string][]byte
	if err := json.Unmarshal(backup, &m); err != nil {
		return err
	}

	for k, v := range m {
		if err := c.Client.Put(
			keyPrefix+k,
			v,
			nil,
		); err != nil {
			return err
		}
	}

	return nil
}

// Before executing DeleteBackend, you should execute BackupBackend
func (c *Client) DeleteBackend(backend string) error {
	return c.Client.DeleteTree(c.Prefix + "/backends/" + backend)
}

func (c *Client) NewFrontend(frontendName, rule string) error {
	return c.AddValueForFrontend(frontendName, "routes", "host", "rule", rule)
}

func (c *Client) AddValueForFrontend(frontendName string, values ...interface{}) error {
	if len(values) < 2 {
		return errors.New("Insufficient parameters")
	}

	arr := make([]string, len(values)-1)

	for i := range arr {
		arr[i] = values[i].(string)
	}

	return c.Client.Put(
		c.Prefix+"/frontends/"+frontendName+"/"+strings.Join(arr, "/"),
		[]byte(fmt.Sprint(values[len(values)-1])),
		nil,
	)
}

func (c *Client) HasFrontend(frontendName string) (bool, error) {
	pairs, err := c.Client.List(c.Prefix + "/frontends/" + frontendName + "/")

	if err != nil {
		if err == store.ErrKeyNotFound {
			err = nil
		}
		return false, err
	}

	return len(pairs) != 0, nil
}

// Before executing DeleteBackend, you should execute BackupBackend
func (c *Client) DeleteFrontend(frontend string) error {
	return c.Client.DeleteTree(c.Prefix + "/backends/" + frontend)
}
