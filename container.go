package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/cs3238-tsuzu/modoki/consul_traefik"

	"github.com/pkg/errors"

	"github.com/docker/docker/api/types"

	"github.com/cs3238-tsuzu/modoki/app"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/goadesign/goa"
	"github.com/jmoiron/sqlx"
)

const containerSchema = `
CREATE TABLE IF NOT EXISTS containers (
	id INT NOT NULL AUTO_INCREMENT,
	cid VARCHAR(128) UNIQUE,
	name VARCHAR(64) NOT NULL UNIQUE,
	uid INT NOT NULL,
	status VARCHAR(32),
	message TEXT,
	PRIMARY KEY (id),
	INDEX(cid, name, uid)
);
`

func InitSchemaForContainer(db *sqlx.DB) error {
	_, err := db.Exec(containerSchema)

	return err
}

type Container struct {
	ID      int
	CID     sql.NullString
	Name    string
	UID     int
	Status  string
	Message string
}

// ContainerController implements the container resource.
type ContainerController struct {
	*goa.Controller
	dockerClient *client.Client
	db           *sqlx.DB
	consul       *consulTraefik.Client
}

// NewContainerController creates a container controller.
func NewContainerController(service *goa.Service, dockerClient *client.Client, db *sqlx.DB, consul *consulTraefik.Client) *ContainerController {
	return &ContainerController{
		Controller:   service.NewController("ContainerController"),
		dockerClient: dockerClient,
		db:           db,
		consul:       consul,
	}
}

func (c *ContainerController) updateStatus(ctx context.Context, status, msg string, id int) error {
	_, err := c.db.ExecContext(ctx, "UPDATE containers SET status=?, message=? WHERE id=?", status, msg, id)

	return err
}

func (c *ContainerController) must(err error) {
	if err != nil {
		log.Println("UpdateStatus error:", err)
	}
}

// Create runs the create action.
func (c *ContainerController) Create(ctx *app.CreateContainerContext) error {
	// ContainerController_Create: start_implement

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.BadRequest(err)
	}

	res, err := c.db.ExecContext(ctx, `INSERT INTO containers (name, uid, status) VALUES (?, ?, "Waiting")`, ctx.Name, uid)

	if err != nil {
		return ctx.InternalServerError(err)
	}

	var id int
	if id64, err := res.LastInsertId(); err != nil {
		return ctx.InternalServerError(err)
	} else {
		id = int(id64)
	}

	go func() {
		c.must(c.updateStatus(context.Background(), "Image Downloading", "", id))

		type ImagePullProgress struct {
			Status         string `json:"status"`
			ProgressDetail struct {
				Current int `json:"current"`
				Total   int `json:"total"`
			} `json:"progressDetail,omitempty"`
			Progress string `json:"progress,omitempty"`
			ID       string `json:"id,omitempty"`
		}

		if rc, err := c.dockerClient.ImagePull(context.Background(), ctx.Image, types.ImagePullOptions{}); err != nil {
			c.must(c.updateStatus(context.Background(), "Error", fmt.Sprintf("Image downloading error: %v", err), id))

			return
		} else {
			defer rc.Close()

			decoder := json.NewDecoder(rc)

			var status string
			for {
				var progress ImagePullProgress

				if err := decoder.Decode(&progress); err != nil {
					break
				}

				status = progress.Status
			}

			if !(strings.Contains(status, "Downloaded") || strings.Contains(status, "up to date")) {
				c.must(c.updateStatus(context.Background(), "Error", fmt.Sprintf("Image downloading error: %v", status), id))

				return
			}
		}

		volumesMap := make(map[string]struct{})

		for i := range ctx.Volumes {
			volumesMap[ctx.Volumes[i]] = struct{}{}
		}

		config := &container.Config{
			Image:      ctx.Image,
			Cmd:        strslice.StrSlice(ctx.Cmd),
			Entrypoint: ctx.Entrypoint,
			Env:        ctx.Env,
			Volumes:    volumesMap,
		}

		if ctx.WorkingDir != nil {
			config.WorkingDir = *ctx.WorkingDir
		}

		hostConfig := &container.HostConfig{}

		networkingConfig := &network.NetworkingConfig{}

		if networkName != nil {
			networkingConfig.EndpointsConfig[*networkName] = &network.EndpointSettings{}
		}

		body, err := c.dockerClient.ContainerCreate(context.Background(), config, hostConfig, networkingConfig, "")

		if err != nil {
			c.must(c.updateStatus(context.Background(), "Error", fmt.Sprintf("Failed to create a container: %v", err), id))

			return
		}

		_, err = c.db.Exec("UPDATE containers SET cid=? where id=?", body.ID, id)

		if err != nil {
			c.must(c.updateStatus(context.Background(), "Error", fmt.Sprintf("Update containers table error: %v", err), id))

			return
		}

		frontendName := fmt.Sprintf(FrontendFormat, id)
		backendName := fmt.Sprintf(BackendFormat, id)
		if err := c.consul.NewFrontend(frontendName, "Host: "+ctx.Name+"."+*publicAddr); err != nil {
			c.must(c.updateStatus(context.Background(), "Error", fmt.Sprintf("Update traefik error: %v", err), id))

			return
		}
		if err := c.consul.AddValueForFrontend(frontendName, "passHostHeader", true); err != nil {
			c.must(c.updateStatus(context.Background(), "Error", fmt.Sprintf("UUpdate traefik error: %v", err), id))

			return
		}

		if err := c.consul.AddValueForFrontend(frontendName, "backend", backendName); err != nil {
			c.must(c.updateStatus(context.Background(), "Error", fmt.Sprintf("Update traefik error: %v", err), id))

			return
		}

		c.must(c.updateStatus(context.Background(), "Created", "", id))
	}()
	cres := &app.GoaContainerCreateResults{
		ID: id,
	}

	return ctx.OK(cres)
	// ContainerController_Create: end_implement
}

// Remove runs the remove action.
func (c *ContainerController) Remove(ctx *app.RemoveContainerContext) error {
	// ContainerController_Remove: start_implement

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(err)
	}

	rows, err := c.db.Query("SELECT id, cid FROM containers WHERE uid=? AND (id=? OR name=?)", uid, ctx.ID, ctx.ID)

	if err != nil {
		return ctx.InternalServerError(errors.Wrap(err, "Database Error"))
	}

	rows.Next()

	var id int
	var cid sql.NullString
	if err := rows.Scan(&id, &cid); err != nil {
		return ctx.NotFound()
	}

	if !cid.Valid {
		return ctx.NotFound()
	}

	if err := c.dockerClient.ContainerRemove(
		context.Background(),
		cid.String,
		types.ContainerRemoveOptions{
			RemoveLinks:   true,
			RemoveVolumes: true,
			Force:         ctx.Force,
		},
	); err != nil {
		if strings.Contains(err.Error(), "You cannot remove a running container") {
			return ctx.RunningContainer()
		}
	}

	if _, err := c.db.Exec("DELETE FROM containers WHERE id=?", id); err != nil {
		return ctx.InternalServerError(errors.Wrap(err, "Deletion From Database Error"))
	}

	frontendName := fmt.Sprintf(FrontendFormat, id)
	backendName := fmt.Sprintf(BackendFormat, id)

	if err := c.consul.DeleteBackend(backendName); err != nil {
		return ctx.InternalServerError(errors.Wrap(err, "consul Error"))
	}
	if err := c.consul.DeleteFrontend(frontendName); err != nil {
		return ctx.InternalServerError(errors.Wrap(err, "consul Error"))
	}

	return ctx.OK(nil)
	// ContainerController_Remove: end_implement
}

// Start runs the start action.
func (c *ContainerController) Start(ctx *app.StartContainerContext) error {
	// ContainerController_Start: start_implement

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(err)
	}

	rows, err := c.db.Query("SELECT id, cid FROM containers WHERE uid=? AND (id=? OR name=?)", uid, ctx.ID, ctx.ID)

	if err != nil {
		return ctx.InternalServerError(errors.Wrap(err, "Database Error"))
	}

	rows.Next()

	var id int
	var cid sql.NullString
	if err := rows.Scan(&id, &cid); err != nil {
		return ctx.NotFound()
	}

	if !cid.Valid {
		return ctx.NotFound()
	}

	resp, err := c.dockerClient.HTTPClient().Post("http://"+*dockerAPIVersion+"/containers/"+cid.String+"/start", "", nil)

	if err != nil {
		return ctx.InternalServerError(errors.Wrap(err, "Docker API error"))
	}
	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()
	switch resp.StatusCode {
	case http.StatusOK:
		// Do nothing
	case 304: // Already started
		return ctx.OK(nil)
	case http.StatusNotFound:
		return ctx.NotFound()
	case http.StatusInternalServerError:
		type message struct {
			Message string `json:"message"`
		}

		var msg message
		json.NewDecoder(resp.Body).Decode(&msg)

		return ctx.InternalServerError(fmt.Errorf("Container starting error: %s", msg.Message))
	}

	if err := c.updateContainerStatus(context.Background(), id, cid.String); err != nil {
		return ctx.InternalServerError(err)
	}

	return ctx.OK(nil)
	// ContainerController_Start: end_implement
}

// Stop runs the stop action.
func (c *ContainerController) Stop(ctx *app.StopContainerContext) error {
	// ContainerController_Stop: start_implement

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(err)
	}

	rows, err := c.db.Query("SELECT id, cid FROM containers WHERE uid=? AND (id=? OR name=?)", uid, ctx.ID, ctx.ID)

	if err != nil {
		return ctx.InternalServerError(errors.Wrap(err, "Database Error"))
	}

	rows.Next()

	var id int
	var cid sql.NullString
	if err := rows.Scan(&id, &cid); err != nil {
		return ctx.NotFound()
	}

	if !cid.Valid {
		return ctx.NotFound()
	}

	d := 15 * time.Second

	if err := c.dockerClient.ContainerStop(context.Background(), cid.String, &d); err != nil {
		return ctx.InternalServerError(errors.Wrap(err, "Docker API error"))
	}

	if err := c.updateContainerStatus(context.Background(), id, cid.String); err != nil {
		return ctx.InternalServerError(err)
	}

	return ctx.OK(nil)
	// ContainerController_Stop: end_implement
}

func (c *ContainerController) updateContainerStatus(ctx context.Context, id int, cid string) error {
	j, err := c.dockerClient.ContainerInspect(ctx, cid)

	if err != nil {
		return errors.Wrap(err, "Container Inspect Error")
	}

	n := "bridge"

	if networkName != nil {
		n = *networkName
	}

	addr := j.NetworkSettings.Networks[n].IPAddress

	backendName := fmt.Sprintf(BackendFormat, id)

	if addr == "" {
		if err := c.consul.DeleteBackend(backendName); err != nil {
			return errors.Wrap(err, "Traefik Unregisteration Error")
		}
	} else {
		if err := c.consul.NewBackend(backendName, ServerName, addr); err != nil {
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
