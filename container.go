package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"code.cloudfoundry.org/bytefmt"
	"github.com/cs3238-tsuzu/modoki/app"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/libkv/store"
	"github.com/goadesign/goa"
	"github.com/pkg/errors"
	"golang.org/x/net/websocket"
)

// ContainerController implements the container resource.
type ContainerController struct {
	*goa.Controller
	*ContainerControllerUtil
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
		return ctx.BadRequest(goa.ErrBadRequest(err))
	}

	res, err := c.DB.ExecContext(ctx, `INSERT INTO containers (name, uid, status) VALUES (?, ?, "Waiting")`, ctx.Name, uid)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	var id int
	if id64, err := res.LastInsertId(); err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
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
				dockerLabelModokiID:   strconv.Itoa(id),
				dockerLabelModokiUID:  strconv.Itoa(uid),
				dockerLabelModokiName: ctx.Name,
			},
		}

		if ctx.WorkingDir != nil {
			config.WorkingDir = *ctx.WorkingDir
		}

		var cpuMaxUsage int64 = 100
		pair, err := c.Consul.Client.Get("modoki/cpu/max_usage")

		if err == nil {
			if c, err := strconv.Atoi(string(pair.Value)); err == nil && c > 0 && c <= 100 {
				cpuMaxUsage = int64(c)
			}
		}

		var memMaxUsage int64
		pair, err = c.Consul.Client.Get("modoki/memory/max_usage")

		if err == nil {
			if u, err := bytefmt.ToBytes(string(pair.Value)); err == nil && u > 0 {
				memMaxUsage = int64(u)
			}
		}

		var storageMaxSize string
		pair, err = c.Consul.Client.Get("modoki/storage/max_usage")

		if err == nil {
			if v, err := bytefmt.ToBytes(string(pair.Value)); err == nil && v > 0 {
				storageMaxSize = string(pair.Value)
			}
		}

		hostConfig := &container.HostConfig{}

		if cpuMaxUsage != 100 {
			hostConfig.Resources.CPUPeriod = 100000
			hostConfig.Resources.CPUQuota = cpuMaxUsage * 1000
		}
		if memMaxUsage != 0 {
			hostConfig.Resources.Memory = memMaxUsage
		}

		if storageMaxSize != "" {
			hostConfig.StorageOpt = map[string]string{
				"size": storageMaxSize,
			}
		}

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

		frontendName := fmt.Sprintf(frontendFormat, id)
		backendName := fmt.Sprintf(backendFormat, id)
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
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	rows, err := c.DB.Query("SELECT id, cid FROM containers WHERE uid=? AND (id=? OR name=?)", uid, ctx.ID, ctx.ID)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(errors.Wrap(err, "Database Error")))
	}

	rows.Next()

	var id int
	var cid sql.NullString
	if err := rows.Scan(&id, &cid); err != nil {
		rows.Close()
		return ctx.NotFound(goa.ErrNotFound(errors.New("No container found")))
	}
	rows.Close()

	if !cid.Valid {
		return ctx.NotFound(goa.ErrNotFound(errors.New("No container found")))
	}

	if ctx.Method == "HEAD" {
		stat, err := c.DockerClient.ContainerStatPath(ctx, cid.String, ctx.InternalPath)

		if err != nil {
			return ctx.NotFound(goa.ErrNotFound(err))
		}

		j, _ := json.Marshal(stat)
		ctx.ResponseWriter.Header().Set("X-Docker-Container-Path-Stat", string(j))
		ctx.ResponseWriter.Header().Set("Content-Type", "application/octet-stream")
		ctx.ResponseWriter.WriteHeader(http.StatusOK)

		return nil
	}

	rc, stat, err := c.DockerClient.CopyFromContainer(ctx, cid.String, ctx.InternalPath)

	if err != nil {
		return ctx.NotFound(goa.ErrNotFound(err))
	}
	defer rc.Close()

	j, _ := json.Marshal(stat)
	ctx.ResponseWriter.Header().Set("X-Docker-Container-Path-Stat", string(j))
	ctx.ResponseWriter.Header().Set("Content-Type", "application/octet-stream")

	ctx.ResponseWriter.WriteHeader(http.StatusOK)
	if _, err := io.Copy(ctx.ResponseWriter, rc); err != nil {
		return err
	}

	return nil
	// ContainerController_Download: end_implement
}

// Remove runs the remove action.
func (c *ContainerController) Remove(ctx *app.RemoveContainerContext) error {
	// ContainerController_Remove: start_implement

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	rows, err := c.DB.Query("SELECT id, cid, status FROM containers WHERE uid=? AND (id=? OR name=?)", uid, ctx.ID, ctx.ID)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(errors.Wrap(err, "Database Error")))
	}

	rows.Next()

	var id int
	var cid sql.NullString
	var status string
	if err := rows.Scan(&id, &cid, &status); err != nil {
		rows.Close()
		return ctx.NotFound()
	}
	rows.Close()

	if status != "Error" {
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
				return ctx.InternalServerError(goa.ErrInternal(errors.Wrap(err, "Docker API Error")))
			}
		}
	}

	if _, err := c.DB.Exec("DELETE FROM containers WHERE id=?", id); err != nil {
		return ctx.InternalServerError(goa.ErrInternal(errors.Wrap(err, "Deletion From Database Error")))
	}

	frontendName := fmt.Sprintf(frontendFormat, id)
	backendName := fmt.Sprintf(backendFormat, id)

	if err := c.Consul.DeleteBackend(backendName); err != nil {
		if err != store.ErrKeyNotFound {
			return ctx.InternalServerError(goa.ErrInternal(errors.Wrap(err, "Consul Error")))
		}

	}
	if err := c.Consul.DeleteFrontend(frontendName); err != nil {
		if err != store.ErrKeyNotFound {
			return ctx.InternalServerError(goa.ErrInternal(errors.Wrap(err, "Consul Error")))
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
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	rows, err := c.DB.Query("SELECT id, cid FROM containers WHERE uid=? AND (id=? OR name=?)", uid, ctx.ID, ctx.ID)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(errors.Wrap(err, "Database Error")))
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
		return ctx.InternalServerError(goa.ErrInternal(errors.Wrap(err, "Docker API error")))
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
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.NoContent()

	// ContainerController_Start: end_implement
}

// Stop runs the stop action.
func (c *ContainerController) Stop(ctx *app.StopContainerContext) error {
	// ContainerController_Stop: start_implement

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	rows, err := c.DB.Query("SELECT id, cid FROM containers WHERE uid=? AND (id=? OR name=?)", uid, ctx.ID, ctx.ID)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(errors.Wrap(err, "Database Error")))
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
		return ctx.InternalServerError(goa.ErrInternal(errors.Wrap(err, "Docker API error")))
	}

	if err := c.updateContainerStatus(context.Background(), cid.String); err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.NoContent()

	// ContainerController_Stop: end_implement
}

// Inspect runs the inspect action.
func (c *ContainerController) Inspect(ctx *app.InspectContainerContext) error {
	// ContainerController_Inspect: start_implement

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	rows, err := c.DB.Query("SELECT id, cid, name, status FROM containers WHERE uid=? AND (id=? OR name=?)", uid, ctx.ID, ctx.ID)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(errors.Wrap(err, "Database Error")))
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
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	filter := filters.NewArgs()

	filter.Add("label", dockerLabelModokiUID+"="+strconv.Itoa(uid))

	list, err := c.DockerClient.ContainerList(ctx, types.ContainerListOptions{
		All:     true,
		Filters: filter,
	})

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(errors.Wrap(err, "Docker API Error")))
	}
	res := make(app.GoaContainerListEachCollection, 0, len(list)+10)

	rows, err := c.DB.Query(`SELECT id, name, message, status FROM containers WHERE status="Error" OR status="Creating"`)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(errors.Wrap(err, "Database Error")))
	}

	defer rows.Close()

	for rows.Next() {
		var id int
		var name, msg, status string

		if err := rows.Scan(&id, &name, &msg, &status); err != nil {
			return ctx.InternalServerError(goa.ErrInternal(errors.Wrap(err, "Database Error")))
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
		id, _ := strconv.Atoi(j.Labels[dockerLabelModokiID])
		name := j.Labels[dockerLabelModokiName]

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
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	rows, err := c.DB.Query("SELECT id, cid FROM containers WHERE uid=? AND (id=? OR name=?)", uid, ctx.ID, ctx.ID)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(errors.Wrap(err, "Database Error")))
	}

	rows.Next()

	var id int
	var cid sql.NullString
	if err := rows.Scan(&id, &cid); err != nil {
		rows.Close()
		return ctx.NotFound(goa.ErrInternal(errors.New("No container found")))
	}
	rows.Close()

	if !cid.Valid {
		return ctx.NotFound(goa.ErrInternal(errors.New("No container found")))
	}

	reader, err := ctx.Payload.Data.Open()

	if err != nil {
		return ctx.BadRequest(goa.ErrBadRequest(errors.Wrap(err, "Opening the form error")))
	}
	defer reader.Close()

	if err := c.DockerClient.CopyToContainer(ctx, cid.String, ctx.Payload.Path, reader, types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: ctx.Payload.AllowOverwrite,
		CopyUIDGID:                ctx.Payload.CopyUIDGID,
	}); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return ctx.NotFound(goa.ErrNotFound(errors.New("The path does not exist")))
		}

		return ctx.InternalServerError(goa.ErrInternal(errors.Wrap(err, "Failed to copy a file via Docker API")))
	}

	return ctx.NoContent()
	// ContainerController_Upload: end_implement
}

// Logs runs the logs action.
func (c *ContainerController) Logs(ctx *app.LogsContainerContext) error {
	// ContainerController_Logs: start_implement

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	rows, err := c.DB.Query("SELECT id, cid FROM containers WHERE uid=? AND (id=? OR name=?)", uid, ctx.ID, ctx.ID)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(errors.Wrap(err, "Database Error")))
	}

	rows.Next()

	var id int
	var cid sql.NullString
	if err := rows.Scan(&id, &cid); err != nil {
		rows.Close()
		return ctx.NotFound(goa.ErrNotFound(errors.New("No container found")))
	}
	rows.Close()

	if !cid.Valid {
		return ctx.NotFound(goa.ErrNotFound(errors.New("No container found")))
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
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	defer rc.Close()

	handler := websocket.Handler(func(conn *websocket.Conn) {
		io.Copy(conn, rc)
	})

	handler.ServeHTTP(ctx.ResponseWriter, ctx.Request)
	return nil
	// ContainerController_Logs: end_implement
}

// SetConfig runs the setConfig action.
func (c *ContainerController) SetConfig(ctx *app.SetConfigContainerContext) error {
	// ContainerController_SetConfig: start_implement

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	var setQuery []string
	var placeholders []interface{}

	if ctx.Payload.DefaultShell != nil {
		setQuery = append(setQuery, "defaultShell=?")
		placeholders = append(placeholders, *ctx.Payload.DefaultShell)
	}

	placeholders = append(placeholders, uid, ctx.ID, ctx.ID)

	res, err := c.DB.Exec("UPDATE containers SET "+strings.Join(setQuery, " ")+" WHERE uid=? AND (id=? OR name=?)", placeholders...)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	if r, _ := res.RowsAffected(); r == 0 {
		return ctx.NotFound()
	}

	return ctx.NoContent()

	// ContainerController_SetConfig: end_implement
}

// GetConfig runs the getConfig action.
func (c *ContainerController) GetConfig(ctx *app.GetConfigContainerContext) error {
	// ContainerController_GetConfig: start_implement

	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	var configs []sql.NullString
	err = c.DB.Select(&configs, "SELECT defaultShell FROM containers WHERE uid=? AND (id=? OR name=?)", uid, ctx.ID, ctx.ID)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	if len(configs) == 0 {
		return ctx.NotFound()
	}

	if !configs[0].Valid {
		return ctx.OK(&app.GoaContainerConfig{nil})
	}

	return ctx.OK(&app.GoaContainerConfig{&configs[0].String})
	// ContainerController_GetConfig: end_implement
}

// Exec runs the exec action.
func (c *ContainerController) Exec(ctx *app.ExecContainerContext) error {
	uid, err := GetUIDFromJWT(ctx)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	tty := false
	if ctx.Tty != nil && *ctx.Tty {
		tty = true
	}

	rows, err := c.DB.Query("SELECT cid, defaultShell FROM containers WHERE uid=? AND (id=? OR name=?)", uid, ctx.ID, ctx.ID)

	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(errors.Wrap(err, "DB error")))
	}

	var cid string
	var defaultShell string
	rows.Next()
	if err := rows.Scan(&cid, &defaultShell); err != nil {
		rows.Close()
		return ctx.InternalServerError(goa.ErrInternal(errors.Wrap(err, "DB error")))
	}

	rows.Close()

	if len(ctx.Command) == 0 {
		ctx.Command = []string{defaultShell}
	}

	if len(ctx.Command) == 0 {
		if p, err := c.Consul.Client.Get(fmt.Sprint(defaultShellKVFormat, uid)); err == nil {
			ctx.Command = []string{string(p.Value)}
		}
	}
	if len(ctx.Command) == 0 {
		ctx.Command = []string{os.Getenv("MODOKI_DEFAULT_SHELL")}
	}
	if len(ctx.Command) == 0 {
		ctx.Command = []string{"sh"}
	}

	c.ExecWSHandler(ctx, cid, ctx.Command, tty).ServeHTTP(ctx.ResponseWriter, ctx.Request)
	return nil
}

func createExecOutgointData(encoder *json.Encoder, kind string, data ...string) error {
	arr := []string{kind}
	arr = append(arr, data...)

	if err := encoder.Encode(arr); err != nil {
		return err
	}

	return nil
}

func parseExecIncoming(decoder *json.Decoder) (string, []string, error) {
	var arr []string
	if err := decoder.Decode(&arr); err != nil {
		return "", nil, err
	}

	if len(arr) < 2 {
		return "", nil, errors.New("invalid format")
	}

	return arr[0], arr[1:], nil
}

// ExecWSHandler establishes a websocket connection to run the exec action.
func (c *ContainerController) ExecWSHandler(ctx *app.ExecContainerContext, cid string, command []string, tty bool) websocket.Handler {
	return func(ws *websocket.Conn) {
		// ContainerController_Exec: start_implement

		execConfig := types.ExecConfig{
			AttachStdin:  true,
			AttachStdout: true,
			AttachStderr: true,
			Detach:       false,
		}

		decoder := json.NewDecoder(ws)
		encoder := json.NewEncoder(ws)

		execID, resp, err := c.initExec(context.Background(), cid, execConfig)

		finalize := func() {
			ws.Close()
			resp.Close()
		}

		if err != nil {
			createExecOutgointData(encoder, "error", err.Error())
			finalize()

			return
		}

		go func() {
			for {
				kind, data, err := parseExecIncoming(decoder)

				if err != nil {
					finalize()

					return
				}

				switch kind {
				case "stdin":
					if len(data) != 0 {
						if _, err := resp.Conn.Write([]byte(data[0])); err != nil {
							finalize()
							return
						}
					}

				case "set_size":
					if len(data) < 2 {
						break
					}

					rows, _ := strconv.Atoi(data[0])
					cols, _ := strconv.Atoi(data[1])

					c.resizeTty(context.Background(), execID, ttySize{
						h: uint(rows),
						w: uint(cols),
					})
				}
			}
		}()

		pr, pw := io.Pipe()

		if execConfig.Tty {
			io.Copy(pw, resp.Conn)
		} else {
			stdcopy.StdCopy(pw, pw, resp.Conn)
		}

		buf := make([]byte, 32*1024)
		for {
			if l, err := pr.Read(buf); err != nil {
				finalize()
				return
			} else {
				if err := createExecOutgointData(encoder, "stdout", string(buf[:l])); err != nil {
					finalize()
					return
				}
			}
		}

		// ContainerController_Exec: end_implement
	}
}
