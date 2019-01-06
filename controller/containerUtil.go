package controller

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	constants "github.com/modoki-paas/modoki/const"
	consulTraefik "github.com/modoki-paas/modoki/consul_traefik"

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

	PublicAddr       string
	HTTPS            bool
	DockerAPIVersion string
	NetworkName      *string
}

func (c *ContainerControllerUtil) updateStatus(ctx context.Context, status, msg string, id int) error {
	_, err := c.DB.ExecContext(ctx, "UPDATE containers SET status=?, message=? WHERE id=?", status, msg, id)

	return err
}

func (c *ContainerControllerUtil) updateContainerStatus(ctx context.Context, cid string) error {
	j, err := c.DockerClient.ContainerInspect(ctx, cid)

	if err != nil {
		return errors.Wrap(err, "Container Inspect Error")
	}

	var id int
	if idStr, ok := j.Config.Labels[constants.DockerLabelModokiID]; !ok {
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
	// TODO: change network

	if c.NetworkName != nil { // command arguments
		n = *c.NetworkName
	}

	addr := j.NetworkSettings.Networks[n].IPAddress

	backendName := fmt.Sprintf(constants.BackendFormat, id)

	if addr == "" {
		if err := c.Consul.DeleteBackend(backendName); err != nil {
			if !strings.Contains(err.Error(), "Key not found") {
				return errors.Wrap(err, "Traefik Unregisteration Error")
			}
		}
	} else {
		if err := c.Consul.NewBackend(backendName, constants.ServerName, "http://"+addr); err != nil {
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

// Run runs event watcher in background TODO: exported as microservice
func (c *ContainerControllerUtil) Run(ctx context.Context) {
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

// TODO: change from payload to parameters

type ContainerCreateParameters struct {
	Name        string
	Image       string
	Command     []string
	Entrypoint  []string
	Env         []string
	Volumes     []string
	WorkingDir  *string
	SSLRedirect bool
}

type ContainerCreateResult struct {
	Endpoints []string
	ID        int
}

type LogsParameters struct {
	Stderr     bool
	Stdout     bool
	Timestamps bool
	Follow     bool
	Tail       string
	Since      *time.Time
	Until      *time.Time
}

type ExecParameters struct {
	Command []string
	Tty     *bool
}

type ContainerInspectRawState struct {
	Dead       bool
	ExitCode   int
	FinishedAt time.Time
	OomKilled  bool
	Paused     bool
	Pid        int
	Restarting bool
	Running    bool
	StartedAt  time.Time
	Status     string
}

type ContainerInspectResult struct {
	Args     []string
	Created  time.Time
	ID       int
	Image    string
	ImageID  string
	Name     string
	Path     string
	RawState *ContainerInspectRawState
	Status   string
	Volumes  []string
}

type ContainerListResult []ContainerListResultElement

type ContainerListResultElement struct {
	Command string
	Created time.Time
	ID      int
	Image   string
	ImageID string
	Name    string
	Status  string
	Volumes []string
}

type DownloadResult struct {
	PathStatJSON string
	Reader       io.ReadCloser
}

type ContainerConfig struct {
	DefaultShell *string
}

type UploadParameters struct {
	AllowOverwrite bool
	CopyUIDGID     bool
	Data           io.Reader
	Path           string
}

const (
	StatusRunningContainer = 409
)
