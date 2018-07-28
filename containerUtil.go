package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/modoki-paas/modoki/consul_traefik"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/jmoiron/sqlx"
	"github.com/k0kubun/pp"
	"github.com/pkg/errors"
)

// ContainerControllerUtil is a utility library for ContainerController
type ContainerControllerUtil struct {
	DB           *sqlx.DB
	DockerClient *client.Client
	Consul       *consulTraefik.Client
}

func (c *ContainerControllerUtil) updateStatus(ctx context.Context, status, msg string, id int) error {
	_, err := c.DB.ExecContext(ctx, "UPDATE containers SET status=?, message=? WHERE id=?", status, msg, id)

	return err
}

func (c *ContainerController) must(err error) {
	if err != nil {
		log.Println("UpdateStatus error:", err)
	}
}

func (c *ContainerControllerUtil) updateContainerStatus(ctx context.Context, cid string) error {
	j, err := c.DockerClient.ContainerInspect(ctx, cid)

	if err != nil {
		return errors.Wrap(err, "Container Inspect Error")
	}

	var id int
	if idStr, ok := j.Config.Labels[dockerLabelModokiID]; !ok {
		return errors.New("This container is not maintained by modoki")
	} else {
		id, err = strconv.Atoi(idStr)

		if err != nil {
			return errors.Wrap(err, "Invalid id format")
		}
	}

	if j.State.Error != "" {
		if err := c.updateStatus(ctx, "Error", j.State.Error, id); err != nil {
			return errors.Wrap(err, "DB Update error")
		}
	}

	n := "bridge"

	if networkName != nil { // command arguments
		n = *networkName
	}

	addr := j.NetworkSettings.Networks[n].IPAddress

	backendName := fmt.Sprintf(backendFormat, id)

	if addr == "" {
		if err := c.Consul.DeleteBackend(backendName); err != nil {
			if !strings.Contains(err.Error(), "Key not found") {
				return errors.Wrap(err, "Traefik Unregisteration Error")
			}
		}
	} else {
		if err := c.Consul.NewBackend(backendName, serverName, "http://"+addr); err != nil {
			return errors.Wrap(err, "Traefik Registeration Error")
		}
	}

	status := ""
	if j.State.Running {
		status = "Running"
	} else {
		status = "Stopped"
	}

	if err := c.updateStatus(ctx, status, "", id); err != nil {
		return errors.Wrap(err, "DB Update error")
	}

	return nil
}

func (c *ContainerControllerUtil) run(ctx context.Context) {
	var fn func()

	fn = func() {
		select {
		case <-ctx.Done():
			return
		default:
		}

		defer fn()
		msg, err := c.DockerClient.Events(ctx, types.EventsOptions{})

		for {
			select {
			case m := <-msg:
				log.Println("event caught: ", pp.Sprint(m))

				switch m.Status {
				case "start":
					c.updateContainerStatus(context.Background(), m.Actor.ID)
				case "die":
					c.updateContainerStatus(context.Background(), m.Actor.ID)
				}
			case e := <-err:
				log.Println("Watching events error: ", e)
				return
			}
		}

	}

	fn()
}
