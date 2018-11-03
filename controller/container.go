package controller

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"code.cloudfoundry.org/bytefmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/libkv/store"
	"github.com/modoki-paas/modoki/app"
	"github.com/modoki-paas/modoki/const"
	"github.com/pkg/errors"
	"golang.org/x/net/websocket"
)

type ContainerControllerImpl struct {
	*ContainerControllerUtil
}

func (c *ContainerControllerImpl) must(err error) {
	if err != nil {
		log.Println("UpdateStatus error:", err)
	}
}

// CreateWithContext runs the create action with context.
func (c *ContainerControllerImpl) CreateWithContext(ctx context.Context, uid string, payload ContainerCreatePayload) (*app.GoaContainerCreateResults, int, error) {
	res, err := c.DB.ExecContext(ctx, `INSERT INTO containers (name, uid, status) VALUES (?, ?, "Waiting")`, payload.Name, uid)

	if err != nil {
		if strings.Contains(err.Error(), "Duplicate") && strings.Contains(err.Error(), "'name'") {
			return nil, http.StatusConflict, errors.New("The name is already used by another container")
		}

		return nil, http.StatusInternalServerError, err
	}

	var id int
	if id64, err := res.LastInsertId(); err != nil {
		return nil, http.StatusInternalServerError, err
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

		if rc, err := c.DockerClient.ImagePull(context.Background(), payload.Image, types.ImagePullOptions{}); err != nil {
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

		for i := range payload.Volumes {
			volumesMap[payload.Volumes[i]] = struct{}{}
		}

		config := &container.Config{
			Image:      payload.Image,
			Cmd:        strslice.StrSlice(payload.Command),
			Entrypoint: payload.Entrypoint,
			Env:        payload.Env,
			Volumes:    volumesMap,
			Labels: map[string]string{
				constants.DockerLabelModokiID:   strconv.Itoa(id),
				constants.DockerLabelModokiUID:  uid,
				constants.DockerLabelModokiName: payload.Name,
			},
		}

		if payload.WorkingDir != nil {
			config.WorkingDir = *payload.WorkingDir
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

		if c.NetworkName != nil {
			networkingConfig.EndpointsConfig = map[string]*network.EndpointSettings{
				*c.NetworkName: &network.EndpointSettings{},
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

		frontendName := fmt.Sprintf(constants.FrontendFormat, id)
		backendName := fmt.Sprintf(constants.BackendFormat, id)
		if err := c.Consul.NewFrontend(frontendName, "Host: "+payload.Name+"."+c.PublicAddr); err != nil {
			c.must(c.updateStatus(context.Background(), "Error", fmt.Sprintf("Update traefik error: %v", err), id))

			return
		}
		if err := c.Consul.AddValueForFrontend(frontendName, "passHostHeader", true); err != nil {
			c.must(c.updateStatus(context.Background(), "Error", fmt.Sprintf("Update traefik error: %v", err), id))

			return
		}

		if c.HTTPS {
			if err := c.Consul.AddValueForFrontend(frontendName, "headers", "sslredirect", payload.SSLRedirect); err != nil {
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

	if c.HTTPS {
		eps = []string{
			"https://" + payload.Name + "." + c.PublicAddr,
			"http://" + payload.Name + "." + c.PublicAddr,
		}
	} else {
		eps = []string{
			"http://" + payload.Name + "." + c.PublicAddr,
		}
	}

	cres := &app.GoaContainerCreateResults{
		ID:        id,
		Endpoints: eps,
	}

	return cres, http.StatusOK, nil
}

// DownloadWithContext runs the download action.
func (c *ContainerControllerImpl) DownloadWithContext(ctx context.Context, uid, idOrName, internalPath string, headerOnly bool) (*DownloadResult, int, error) {
	rows, err := c.DB.Query("SELECT id, cid FROM containers WHERE uid=? AND (id=? OR name=?)", uid, idOrName, idOrName)

	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "Database Error")
	}

	rows.Next()

	var id int
	var cid sql.NullString
	if err := rows.Scan(&id, &cid); err != nil {
		rows.Close()
		return nil, http.StatusNotFound, errors.New("No container found")
	}
	rows.Close()

	if !cid.Valid {
		return nil, http.StatusNotFound, errors.New("No container found")
	}

	if headerOnly {
		stat, err := c.DockerClient.ContainerStatPath(ctx, cid.String, internalPath)

		if err != nil {
			return nil, http.StatusNotFound, err
		}

		j, _ := json.Marshal(stat)
		dr := &DownloadResult{
			PathStatJSON: string(j),
		}

		return dr, http.StatusOK, nil
	}

	rc, stat, err := c.DockerClient.CopyFromContainer(ctx, cid.String, internalPath)

	if err != nil {
		return nil, http.StatusNotFound, err
	}

	j, _ := json.Marshal(stat)
	dr := &DownloadResult{
		Reader:       rc,
		PathStatJSON: string(j),
	}

	return dr, http.StatusOK, nil
}

// Remove runs the remove action.
func (c *ContainerControllerImpl) Remove(uid, idOrName string, force bool) (int, error) {
	rows, err := c.DB.Query("SELECT id, cid, status FROM containers WHERE uid=? AND (id=? OR name=?)", uid, idOrName, idOrName)

	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "Database Error")
	}

	rows.Next()

	var id int
	var cid sql.NullString
	var status string
	if err := rows.Scan(&id, &cid, &status); err != nil {
		rows.Close()
		return http.StatusNotFound, nil
	}
	rows.Close()

	if status != "Error" {
		if !cid.Valid {
			return http.StatusNotFound, nil
		}

		if err := c.DockerClient.ContainerRemove(
			context.Background(),
			cid.String,
			types.ContainerRemoveOptions{
				RemoveVolumes: true,
				Force:         force,
			},
		); err != nil {
			if strings.Contains(err.Error(), "You cannot remove a running container") {
				return StatusRunningContainer, errors.New("container is running")
			} else {
				return http.StatusInternalServerError, errors.Wrap(err, "Docker API Error")
			}
		}
	}

	if _, err := c.DB.Exec("DELETE FROM containers WHERE id=?", id); err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "Deletion From Database Error")
	}

	frontendName := fmt.Sprintf(constants.FrontendFormat, id)
	backendName := fmt.Sprintf(constants.BackendFormat, id)

	if err := c.Consul.DeleteBackend(backendName); err != nil {
		if err != store.ErrKeyNotFound {
			return http.StatusInternalServerError, errors.Wrap(err, "Consul Error")
		}

	}
	if err := c.Consul.DeleteFrontend(frontendName); err != nil {
		if err != store.ErrKeyNotFound {
			return http.StatusInternalServerError, errors.Wrap(err, "Consul Error")
		}
	}

	return http.StatusNoContent, nil
}

// Start runs the start action.
func (c *ContainerControllerImpl) Start(uid, irOrName string) (int, error) {
	rows, err := c.DB.Query("SELECT id, cid FROM containers WHERE uid=? AND (id=? OR name=?)", uid, irOrName)

	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "Database Error")
	}

	rows.Next()

	var id int
	var cid sql.NullString
	if err := rows.Scan(&id, &cid); err != nil {
		rows.Close()
		return http.StatusNotFound, nil
	}
	rows.Close()

	if !cid.Valid {
		return http.StatusNotFound, nil
	}

	resp, err := c.DockerClient.HTTPClient().Post("http://"+c.DockerAPIVersion+"/containers/"+cid.String+"/start", "", nil)

	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "Docker API error")
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
		return http.StatusNoContent, nil
	case http.StatusNotFound:
		return http.StatusNotFound, nil
	case http.StatusInternalServerError:
		type message struct {
			Message string `json:"message"`
		}

		var msg message
		json.NewDecoder(resp.Body).Decode(&msg)

		return http.StatusInternalServerError, fmt.Errorf("Container starting error: %s", msg.Message)
	}

	if err := c.updateContainerStatus(context.Background(), cid.String); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusNoContent, nil
}

// Stop runs the stop action.
func (c *ContainerControllerImpl) Stop(uid, irOrName string) (int, error) {
	rows, err := c.DB.Query("SELECT id, cid FROM containers WHERE uid=? AND (id=? OR name=?)", uid, irOrName, irOrName)

	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "Database Error")
	}

	rows.Next()

	var id int
	var cid sql.NullString
	if err := rows.Scan(&id, &cid); err != nil {
		rows.Close()
		return http.StatusNotFound, nil
	}
	rows.Close()

	if !cid.Valid {
		return http.StatusNotFound, nil
	}

	d := 15 * time.Second

	if err := c.DockerClient.ContainerStop(context.Background(), cid.String, &d); err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "Docker API error")
	}

	if err := c.updateContainerStatus(context.Background(), cid.String); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusNoContent, nil
}

// InspectWithContext runs the inspect action.
func (c *ContainerControllerImpl) InspectWithContext(ctx context.Context, uid, idOrName string) (*app.GoaContainerInspect, int, error) {
	rows, err := c.DB.Query("SELECT id, cid, name, status FROM containers WHERE uid=? AND (id=? OR name=?)", uid, idOrName, idOrName)

	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "Database Error")
	}

	rows.Next()

	var id int
	var cid sql.NullString
	var name, status string
	if err := rows.Scan(&id, &cid, &name, &status); err != nil {
		rows.Close()
		return nil, http.StatusNotFound, nil
	}
	rows.Close()

	if status == "Error" || status == "Creating" {
		insp := &app.GoaContainerInspect{
			ID:     id,
			Name:   name,
			Status: status,
		}

		return insp, http.StatusOK, nil
	}

	if !cid.Valid {
		return nil, http.StatusNotFound, nil
	}

	j, err := c.DockerClient.ContainerInspect(ctx, cid.String)

	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "Container Inspect Error")
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

	return insp, http.StatusOK, nil
}

// ListWithContext runs the list action.
func (c *ContainerControllerImpl) ListWithContext(ctx context.Context, uid string) (app.GoaContainerListEachCollection, int, error) {
	filter := filters.NewArgs()

	filter.Add("label", constants.DockerLabelModokiUID+"="+uid)

	list, err := c.DockerClient.ContainerList(ctx, types.ContainerListOptions{
		All:     true,
		Filters: filter,
	})

	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "Docker API Error")
	}
	res := make(app.GoaContainerListEachCollection, 0, len(list)+10)

	rows, err := c.DB.Query(`SELECT id, name, message, status FROM containers WHERE status="Error" OR status="Creating"`)

	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "Database Error")
	}

	defer rows.Close()

	for rows.Next() {
		var id int
		var name, msg, status string

		if err := rows.Scan(&id, &name, &msg, &status); err != nil {
			return nil, http.StatusInternalServerError, errors.Wrap(err, "Database Error")
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
		id, _ := strconv.Atoi(j.Labels[constants.DockerLabelModokiID])
		name := j.Labels[constants.DockerLabelModokiName]

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

	return res, http.StatusOK, nil
}

// UploadWithContext runs the upload action.
func (c *ContainerControllerImpl) UploadWithContext(ctx context.Context, uid, idOrName string, payload *app.UploadPayload) (int, error) {
	rows, err := c.DB.Query("SELECT id, cid FROM containers WHERE uid=? AND (id=? OR name=?)", uid, idOrName, idOrName)

	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "Database Error")
	}

	rows.Next()

	var id int
	var cid sql.NullString
	if err := rows.Scan(&id, &cid); err != nil {
		rows.Close()
		return http.StatusNotFound, errors.New("No container found")
	}
	rows.Close()

	if !cid.Valid {
		return http.StatusNotFound, errors.New("No container found")
	}

	reader, err := payload.Data.Open()

	if err != nil {
		return http.StatusBadRequest, errors.Wrap(err, "Opening the form error")
	}
	defer reader.Close()

	if err := c.DockerClient.CopyToContainer(ctx, cid.String, payload.Path, reader, types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: payload.AllowOverwrite,
		CopyUIDGID:                payload.CopyUIDGID,
	}); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return http.StatusNotFound, errors.New("The path does not exist")
		}

		return http.StatusInternalServerError, errors.Wrap(err, "Failed to copy a file via Docker API")
	}

	return http.StatusNoContent, nil
}

// LogsWithContext runs the logs action.
func (c *ContainerControllerImpl) LogsWithContext(ctx context.Context, uid, idOrName string, payload LogsPayload) (io.ReadCloser, int, error) {
	rows, err := c.DB.Query("SELECT id, cid FROM containers WHERE uid=? AND (id=? OR name=?)", uid, idOrName, idOrName)

	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "Database Error")
	}

	rows.Next()

	var id int
	var cid sql.NullString
	if err := rows.Scan(&id, &cid); err != nil {
		rows.Close()
		return nil, http.StatusNotFound, errors.New("No container found")
	}
	rows.Close()

	if !cid.Valid {
		return nil, http.StatusNotFound, errors.New("No container found")
	}

	opts := types.ContainerLogsOptions{
		ShowStderr: payload.Stderr,
		ShowStdout: payload.Stdout,
		Timestamps: payload.Timestamps,
		Follow:     payload.Follow,
		Tail:       payload.Tail,
	}

	if payload.Since != nil {
		opts.Since = payload.Since.Format(time.RFC3339)
	}
	if payload.Until != nil {
		opts.Until = payload.Until.Format(time.RFC3339)
	}

	rc, err := c.DockerClient.ContainerLogs(ctx, cid.String, opts)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return rc, http.StatusOK, nil
}

// SetConfig runs the setConfig action.
func (c *ContainerControllerImpl) SetConfig(uid, idOrName string, payload *app.ContainerConfig) (int, error) {
	var setQuery []string
	var placeholders []interface{}

	if payload.DefaultShell != nil {
		setQuery = append(setQuery, "defaultShell=?")
		placeholders = append(placeholders, *payload.DefaultShell)
	}

	placeholders = append(placeholders, uid, idOrName, idOrName)

	tx, err := c.DB.Begin()

	if err != nil {
		return http.StatusInternalServerError, err
	}

	var cnt int

	err = tx.QueryRow(
		"SELECT COUNT(id) FROM containers WHERE uid=? AND (id=? OR name=?)",
		uid, idOrName, idOrName,
	).Scan(&cnt)

	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}

	if cnt == 0 {
		tx.Rollback()

		return http.StatusNotFound, nil
	}

	_, err = tx.Exec("UPDATE containers SET "+strings.Join(setQuery, " ")+" WHERE uid=? AND (id=? OR name=?)", placeholders...)

	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}

	tx.Commit()

	return http.StatusNoContent, nil
}

// GetConfig runs the getConfig action.
func (c *ContainerControllerImpl) GetConfig(uid, idOrName string) (*app.GoaContainerConfig, int, error) {
	var configs []sql.NullString
	err := c.DB.Select(&configs, "SELECT defaultShell FROM containers WHERE uid=? AND (id=? OR name=?)", uid, idOrName, idOrName)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if len(configs) == 0 {
		return nil, http.StatusNotFound, nil
	}

	if !configs[0].Valid {
		return &app.GoaContainerConfig{nil}, http.StatusOK, nil
	}

	return &app.GoaContainerConfig{&configs[0].String}, http.StatusOK, nil
}

// ExecWithContext runs the exec action.
func (c *ContainerControllerImpl) ExecWithContext(ctx context.Context, uid, idOrName string, payload *ExecParameters) (websocket.Handler, int, error) {
	tty := false
	if payload.Tty != nil && *payload.Tty {
		tty = true
	}

	rows, err := c.DB.Query("SELECT cid, defaultShell FROM containers WHERE uid=? AND (id=? OR name=?)", uid, idOrName, idOrName)

	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, "DB error")
	}

	var cid string
	var defaultShell string
	rows.Next()
	if err := rows.Scan(&cid, &defaultShell); err != nil {
		rows.Close()
		return nil, http.StatusInternalServerError, errors.Wrap(err, "DB error")
	}

	rows.Close()

	if len(payload.Command) == 0 {
		payload.Command = []string{defaultShell}
	}

	if len(payload.Command) == 0 {
		if p, err := c.Consul.Client.Get(fmt.Sprint(constants.DefaultShellKVFormat, uid)); err == nil {
			payload.Command = []string{string(p.Value)}
		}
	}
	if len(payload.Command) == 0 {
		payload.Command = []string{os.Getenv("MODOKI_DEFAULT_SHELL")}
	}
	if len(payload.Command) == 0 {
		payload.Command = []string{"sh"}
	}

	return c.execWSHandler(cid, payload.Command, tty), http.StatusOK, nil
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

// execWSHandler establishes a websocket connection to run the exec action.
func (c *ContainerControllerImpl) execWSHandler(cid string, command []string, tty bool) websocket.Handler {
	return func(ws *websocket.Conn) {
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

		defer pr.Close()
		defer pw.Close()

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
	}
}
