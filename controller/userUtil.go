package controller

import (
	"github.com/docker/docker/client"
	"github.com/jmoiron/sqlx"
	"github.com/modoki-paas/modoki/consul_traefik"
)

type UserControllerUtil struct {
	DB           *sqlx.DB
	DockerClient *client.Client
	Consul       *consulTraefik.Client
}

type SSHKey struct {
	Label string
	Key   string
}

type UserAPIKey struct {
	ID     int
	UID    string
	ApiKey string
}

func generateAPIKey() string {
	// TODO:
	return ""
}
