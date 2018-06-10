package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/websocket"

	"github.com/k0kubun/pp"

	"github.com/docker/docker/api/types/filters"

	"github.com/cs3238-tsuzu/modoki/consul_traefik"

	"github.com/docker/docker/client"

	"github.com/jmoiron/sqlx"

	"github.com/cs3238-tsuzu/modoki/app"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/libkv/store"
	"github.com/goadesign/goa"
	"github.com/pkg/errors"
)

// ContainerController implements the container resource.
type ContainerController struct {
	*goa.Controller

	DB           *sqlx.DB
	DockerClient *client.Client
	Consul       *consulTraefik.Client
}

// NewContainerController creates a container controller.
func NewContainerController(service *goa.Service) *ContainerController {
	return &ContainerController{Controller: service.NewController("ContainerController")}
}

// Create runs the create action.
func (c *ContainerController) Create(ctx *app.CreateContainerContext) error {
	// ContainerController_Create: start_implement

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.BadRequest(err)
	}

	res, err := c.DB.ExecContext(ctx, `INSERT INTO containers (name, uid, status) VALUES (?, ?, "Waiting")`, ctx.Name, uid)

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
		c.must(c.updateStatus(context.Background(), "Creating", "", id))

		type ImagePullProgress struct {
			Status         string `json:"status"`
			ProgressDetail struct {
				Current int `json:"current"`
				Total   int `json:"total"`
			} `json:"progressDetail,omitempty"`
			Progress string `json:"progress,omitempty"`
			ID       string `json:"id,omitempty"`
		}

		if rc, err := c.DockerClient.ImagePull(context.Background(), ctx.Image, types.ImagePullOptions{}); err != nil {
			c.must(c.updateStatus(context.Background(), "Error", fmt.Sprintf("Downloading the image error: %v", err), id))

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
			Cmd:        strslice.StrSlice(ctx.Command),
			Entrypoint: ctx.Entrypoint,
			Env:        ctx.Env,
			Volumes:    volumesMap,
			Labels: map[string]string{
				DockerLabelModokiID:   strconv.Itoa(id),
				DockerLabelModokiUID:  strconv.Itoa(uid),
				DockerLabelModokiName: ctx.Name,
			},
		}

		if ctx.WorkingDir != nil {
			config.WorkingDir = *ctx.WorkingDir
		}

		hostConfig := &container.HostConfig{}

		networkingConfig := &network.NetworkingConfig{}

		if networkName != nil {
			networkingConfig.EndpointsConfig = map[string]*network.EndpointSettings{
				*networkName: &network.EndpointSettings{},
			}
		}

		body, err := c.DockerClient.ContainerCreate(context.Background(), config, hostConfig, networkingConfig, "")

		if err != nil {
			c.must(c.updateStatus(context.Background(), "Error", fmt.Sprintf("Failed to create a container: %v", err), id))

			return
		}

		_, err = c.DB.Exec("UPDATE containers SET cid=? where id=?", body.ID, id)

		if err != nil {
			c.must(c.updateStatus(context.Background(), "Error", fmt.Sprintf("Update containers table error: %v", err), id))

			return
		}

		frontendName := fmt.Sprintf(FrontendFormat, id)
		backendName := fmt.Sprintf(BackendFormat, id)
		if err := c.Consul.NewFrontend(frontendName, "Host: "+ctx.Name+"."+*publicAddr); err != nil {
			c.must(c.updateStatus(context.Background(), "Error", fmt.Sprintf("Update traefik error: %v", err), id))

			return
		}
		if err := c.Consul.AddValueForFrontend(frontendName, "passHostHeader", true); err != nil {
			c.must(c.updateStatus(context.Background(), "Error", fmt.Sprintf("Update traefik error: %v", err), id))

			return
		}

		if *https {
			if err := c.Consul.AddValueForFrontend(frontendName, "headers", "sslredirect", ctx.SslRedirect); err != nil {
				c.must(c.updateStatus(context.Background(), "Error", fmt.Sprintf("Update traefik error: %v", err), id))

				return
			}
		}

		if err := c.Consul.AddValueForFrontend(frontendName, "backend", backendName); err != nil {
			c.must(c.updateStatus(context.Background(), "Error", fmt.Sprintf("Update traefik error: %v", err), id))

			return
		}

		c.must(c.updateStatus(context.Background(), "Created", "", id))
	}()

	var eps []string

	if *https {
		eps = []string{
			"https://" + ctx.Name + "." + *publicAddr,
			"http://" + ctx.Name + "." + *publicAddr,
		}
	} else {
		eps = []string{
			"http://" + ctx.Name + "." + *publicAddr,
		}
	}

	cres := &app.GoaContainerCreateResults{
		ID:        id,
		Endpoints: eps,
	}

	return ctx.OK(cres)

	// ContainerController_Create: end_implement
}

// Download runs the download action.
func (c *ContainerController) Download(ctx *app.DownloadContainerContext) error {
	// ContainerController_Download: start_implement

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(err)
	}

	rows, err := c.DB.Query("SELECT id, cid FROM containers WHERE uid=? AND (id=? OR name=?)", uid, ctx.ID, ctx.ID)

	if err != nil {
		return ctx.InternalServerError(errors.Wrap(err, "Database Error"))
	}

	rows.Next()

	var id int
	var cid sql.NullString
	if err := rows.Scan(&id, &cid); err != nil {
		rows.Close()
		return ctx.NotFound(errors.New("No container found"))
	}
	rows.Close()

	if !cid.Valid {
		return ctx.NotFound(errors.New("No container found"))
	}

	rc, stat, err := c.DockerClient.CopyFromContainer(ctx, cid.String, ctx.InternalPath)

	if err != nil {
		return ctx.NotFound(err)
	}
	defer rc.Close()

	ctx.ResponseWriter.Header().Set("Content-Length", strconv.FormatInt(stat.Size, 10))
	ctx.ResponseWriter.Header().Set("Content-Type", "application/octet-stream")

	io.Copy(ctx.ResponseWriter, rc)

	return nil
	// ContainerController_Download: end_implement
}

// Remove runs the remove action.
func (c *ContainerController) Remove(ctx *app.RemoveContainerContext) error {
	// ContainerController_Remove: start_implement

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(err)
	}

	rows, err := c.DB.Query("SELECT id, cid FROM containers WHERE uid=? AND (id=? OR name=?)", uid, ctx.ID, ctx.ID)

	if err != nil {
		return ctx.InternalServerError(errors.Wrap(err, "Database Error"))
	}

	rows.Next()

	var id int
	var cid sql.NullString
	if err := rows.Scan(&id, &cid); err != nil {
		rows.Close()
		return ctx.NotFound()
	}
	rows.Close()

	if !cid.Valid {
		return ctx.NotFound()
	}

	if err := c.DockerClient.ContainerRemove(
		context.Background(),
		cid.String,
		types.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         ctx.Force,
		},
	); err != nil {
		if strings.Contains(err.Error(), "You cannot remove a running container") {
			return ctx.RunningContainer()
		} else {
			return ctx.InternalServerError(errors.Wrap(err, "Docker API Error"))
		}
	}

	if _, err := c.DB.Exec("DELETE FROM containers WHERE id=?", id); err != nil {
		return ctx.InternalServerError(errors.Wrap(err, "Deletion From Database Error"))
	}

	frontendName := fmt.Sprintf(FrontendFormat, id)
	backendName := fmt.Sprintf(BackendFormat, id)

	if err := c.Consul.DeleteBackend(backendName); err != nil {
		if err != store.ErrKeyNotFound {
			return ctx.InternalServerError(errors.Wrap(err, "Consul Error"))
		}

	}
	if err := c.Consul.DeleteFrontend(frontendName); err != nil {
		if err != store.ErrKeyNotFound {
			return ctx.InternalServerError(errors.Wrap(err, "Consul Error"))
		}
	}

	return ctx.NoContent()

	// ContainerController_Remove: end_implement
}

// Start runs the start action.
func (c *ContainerController) Start(ctx *app.StartContainerContext) error {
	// ContainerController_Start: start_implement

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(err)
	}

	rows, err := c.DB.Query("SELECT id, cid FROM containers WHERE uid=? AND (id=? OR name=?)", uid, ctx.ID, ctx.ID)

	if err != nil {
		return ctx.InternalServerError(errors.Wrap(err, "Database Error"))
	}

	rows.Next()

	var id int
	var cid sql.NullString
	if err := rows.Scan(&id, &cid); err != nil {
		rows.Close()
		return ctx.NotFound()
	}
	rows.Close()

	if !cid.Valid {
		return ctx.NotFound()
	}

	resp, err := c.DockerClient.HTTPClient().Post("http://"+*dockerAPIVersion+"/containers/"+cid.String+"/start", "", nil)

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
		return ctx.NoContent()
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

	if err := c.updateContainerStatus(context.Background(), cid.String); err != nil {
		return ctx.InternalServerError(err)
	}

	return ctx.NoContent()

	// ContainerController_Start: end_implement
}

// Stop runs the stop action.
func (c *ContainerController) Stop(ctx *app.StopContainerContext) error {
	// ContainerController_Stop: start_implement

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(err)
	}

	rows, err := c.DB.Query("SELECT id, cid FROM containers WHERE uid=? AND (id=? OR name=?)", uid, ctx.ID, ctx.ID)

	if err != nil {
		return ctx.InternalServerError(errors.Wrap(err, "Database Error"))
	}

	rows.Next()

	var id int
	var cid sql.NullString
	if err := rows.Scan(&id, &cid); err != nil {
		rows.Close()
		return ctx.NotFound()
	}
	rows.Close()

	if !cid.Valid {
		return ctx.NotFound()
	}

	d := 15 * time.Second

	if err := c.DockerClient.ContainerStop(context.Background(), cid.String, &d); err != nil {
		return ctx.InternalServerError(errors.Wrap(err, "Docker API error"))
	}

	if err := c.updateContainerStatus(context.Background(), cid.String); err != nil {
		return ctx.InternalServerError(err)
	}

	return ctx.NoContent()

	// ContainerController_Stop: end_implement
}

// Inspect runs the inspect action.
func (c *ContainerController) Inspect(ctx *app.InspectContainerContext) error {
	// ContainerController_Inspect: start_implement

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(err)
	}

	rows, err := c.DB.Query("SELECT id, cid, name, status FROM containers WHERE uid=? AND (id=? OR name=?)", uid, ctx.ID, ctx.ID)

	if err != nil {
		return ctx.InternalServerError(errors.Wrap(err, "Database Error"))
	}

	rows.Next()

	var id int
	var cid sql.NullString
	var name, status string
	if err := rows.Scan(&id, &cid, &name, &status); err != nil {
		rows.Close()
		return ctx.NotFound()
	}
	rows.Close()

	if status == "Error" || status == "Created" {
		insp := &app.GoaContainerInspect{
			ID:     id,
			Name:   name,
			Status: status,
		}

		return ctx.OK(insp)
	}

	if !cid.Valid {
		return ctx.NotFound()
	}

	j, err := c.DockerClient.ContainerInspect(ctx, cid.String)

	if err != nil {
		return errors.Wrap(err, "Container Inspect Error")
	}

	t, _ := time.Parse(time.RFC3339Nano, j.Created)

	vols := make([]string, 0, len(j.Config.Volumes))
	for k := range j.Config.Volumes {
		vols = append(vols, k)
	}

	insp := &app.GoaContainerInspect{
		Args:    j.Args,
		Created: t,
		ID:      id,
		Image:   j.Config.Image,
		ImageID: j.Image,
		Name:    name,
		Path:    j.Path,
		Volumes: vols,
	}

	rawState := j.State

	switch rawState.Status {
	case "created":
		insp.Status = "Created"
	case "Running":
		insp.Status = "Running"
	default:
		insp.Status = "Stopped"
	}

	s, _ := time.Parse(time.RFC3339Nano, rawState.StartedAt)
	f, _ := time.Parse(time.RFC3339Nano, rawState.FinishedAt)

	insp.RawState = &app.GoaContainerInspectRawState{
		Dead:       rawState.Dead,
		ExitCode:   rawState.ExitCode,
		FinishedAt: f,
		OomKilled:  rawState.OOMKilled,
		Paused:     rawState.Paused,
		Pid:        rawState.Pid,
		Restarting: rawState.Restarting,
		Running:    rawState.Running,
		StartedAt:  s,
		Status:     rawState.Status,
	}

	return ctx.OK(insp)
	// ContainerController_Inspect: end_implement
}

// Inspect runs the inspect action.
func (c *ContainerController) List(ctx *app.ListContainerContext) error {
	// ContainerController_List: start_implement

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(err)
	}

	filter := filters.NewArgs()

	filter.Add("label", DockerLabelModokiUID+"="+strconv.Itoa(uid))

	list, err := c.DockerClient.ContainerList(ctx, types.ContainerListOptions{
		All:     true,
		Filters: filter,
	})

	if err != nil {
		return ctx.InternalServerError(errors.Wrap(err, "Docker API Error"))
	}
	res := make(app.GoaContainerListEachCollection, 0, len(list)+10)

	rows, err := c.DB.Query(`SELECT id, name, message, status FROM containers WHERE status="Error" OR status="Creating"`)

	if err != nil {
		return ctx.InternalServerError(errors.Wrap(err, "Database Error"))
	}

	defer rows.Close()

	for rows.Next() {
		var id int
		var name, msg, status string

		if err := rows.Scan(&id, &name, &msg, &status); err != nil {
			return ctx.InternalServerError(errors.Wrap(err, "Database Error"))
		}

		res = append(res, &app.GoaContainerListEach{
			ID:      id,
			Name:    name,
			Command: msg,
			Status:  status,
		})
	}
	rows.Close()

	for i := range list {
		j := list[i]

		vols := make([]string, 0, len(j.Mounts))
		for k := range j.Mounts {
			vols = append(vols, j.Mounts[k].Destination)
		}

		t := time.Unix(j.Created, 0)
		id, _ := strconv.Atoi(j.Labels[DockerLabelModokiID])
		name := j.Labels[DockerLabelModokiName]

		var state string

		switch strings.ToLower(j.State) {
		case "running":
			state = "Running"
		case "created":
			state = "Created"
		default:
			state = "Stopped"
		}

		each := &app.GoaContainerListEach{
			Command: j.Command,
			Created: t,
			ID:      id,
			Image:   j.Image,
			ImageID: j.ImageID,
			Name:    name,
			Status:  state,
			Volumes: vols,
		}

		res = append(res, each)
	}

	return ctx.OK(res)
	// ContainerController_List: end_implement
}

// Upload runs the upload action.
func (c *ContainerController) Upload(ctx *app.UploadContainerContext) error {
	// ContainerController_Upload: start_implement

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(err)
	}

	rows, err := c.DB.Query("SELECT id, cid FROM containers WHERE uid=? AND (id=? OR name=?)", uid, ctx.Payload.ID, ctx.Payload.ID)

	if err != nil {
		return ctx.InternalServerError(errors.Wrap(err, "Database Error"))
	}

	rows.Next()

	var id int
	var cid sql.NullString
	if err := rows.Scan(&id, &cid); err != nil {
		rows.Close()
		return ctx.NotFound(errors.New("No container found"))
	}
	rows.Close()

	if !cid.Valid {
		return ctx.NotFound(errors.New("No container found"))
	}

	reader, err := ctx.Payload.Data.Open()

	if err != nil {
		return ctx.BadRequest(errors.Wrap(err, "Opening the form error"))
	}
	defer reader.Close()

	if err := c.DockerClient.CopyToContainer(ctx, cid.String, ctx.Payload.Path, reader, types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: ctx.Payload.AllowOverwrite,
	}); err != nil {
		return ctx.InternalServerError(errors.Wrap(err, "Failed to copy a file via Docker API"))
	}

	return ctx.NoContent()
	// ContainerController_Upload: end_implement
}

// Logs runs the logs action.
func (c *ContainerController) Logs(ctx *app.LogsContainerContext) error {
	// ContainerController_Logs: start_implement

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(err)
	}

	rows, err := c.DB.Query("SELECT id, cid FROM containers WHERE uid=? AND (id=? OR name=?)", uid, ctx.ID, ctx.ID)

	if err != nil {
		return ctx.InternalServerError(errors.Wrap(err, "Database Error"))
	}

	rows.Next()

	var id int
	var cid sql.NullString
	if err := rows.Scan(&id, &cid); err != nil {
		rows.Close()
		return ctx.NotFound(errors.New("No container found"))
	}
	rows.Close()

	if !cid.Valid {
		return ctx.NotFound(errors.New("No container found"))
	}

	opts := types.ContainerLogsOptions{
		ShowStderr: ctx.Stderr,
		ShowStdout: ctx.Stdout,
		Timestamps: ctx.Timestamps,
		Follow:     ctx.Follow,
		Tail:       ctx.Tail,
	}

	if ctx.Since != nil {
		opts.Since = ctx.Since.Format(time.RFC3339)
	}
	if ctx.Until != nil {
		opts.Until = ctx.Until.Format(time.RFC3339)
	}

	rc, err := c.DockerClient.ContainerLogs(ctx, cid.String, opts)

	if err != nil {
		return ctx.InternalServerError(err)
	}
	defer rc.Close()

	handler := websocket.Handler(func(conn *websocket.Conn) {
		io.Copy(conn, rc)
	})

	handler.ServeHTTP(ctx.ResponseWriter, ctx.Request)
	return nil
	// ContainerController_Logs: end_implement
}

func (c *ContainerController) updateStatus(ctx context.Context, status, msg string, id int) error {
	_, err := c.DB.ExecContext(ctx, "UPDATE containers SET status=?, message=? WHERE id=?", status, msg, id)

	return err
}

func (c *ContainerController) must(err error) {
	if err != nil {
		log.Println("UpdateStatus error:", err)
	}
}

func (c *ContainerController) updateContainerStatus(ctx context.Context, cid string) error {
	j, err := c.DockerClient.ContainerInspect(ctx, cid)

	if err != nil {
		return errors.Wrap(err, "Container Inspect Error")
	}

	n := "bridge"

	if networkName != nil {
		n = *networkName
	}

	var id int
	if idStr, ok := j.Config.Labels[DockerLabelModokiID]; !ok {
		return errors.New("This container is not maintained by modoki")
	} else {
		id, err = strconv.Atoi(idStr)

		if err != nil {
			return errors.Wrap(err, "Invalid id format")
		}
	}

	addr := j.NetworkSettings.Networks[n].IPAddress

	backendName := fmt.Sprintf(BackendFormat, id)

	if addr == "" {
		if err := c.Consul.DeleteBackend(backendName); err != nil {
			return errors.Wrap(err, "Traefik Unregisteration Error")
		}
	} else {
		if err := c.Consul.NewBackend(backendName, ServerName, "http://"+addr); err != nil {
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

func (c *ContainerController) run(ctx context.Context) {
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
