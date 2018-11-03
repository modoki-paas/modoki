package controller

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
)

type ExecResult struct {
	Reader io.ReadCloser
	Writer io.WriteCloser
}

type ttySize struct {
	w, h uint
}

func (c *ContainerControllerUtil) createExec(ctx context.Context, cid string, execConfig types.ExecConfig) (string, error) {
	resp, err := c.DockerClient.ContainerExecCreate(ctx, cid, execConfig)

	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func (c *ContainerControllerUtil) initExec(ctx context.Context, cid string, execConfig types.ExecConfig) (string, *types.HijackedResponse, error) {
	execID, err := c.createExec(ctx, cid, execConfig)

	if err != nil {
		return "", nil, err
	}

	resp, err := c.DockerClient.ContainerExecAttach(context.Background(), execID, types.ExecStartCheck{
		Tty:    execConfig.Tty,
		Detach: false,
	})

	if err != nil {
		return "", nil, err
	}

	return execID, &resp, nil
}

func (c *ContainerControllerUtil) resizeTty(ctx context.Context, execID string, size ttySize) error {
	if len(execID) == 0 {
		return nil
	}

	if size.w == 0 && size.h == 0 {
		return nil
	}

	return c.DockerClient.ContainerExecResize(ctx, execID, types.ResizeOptions{Width: size.w, Height: size.h})
}
