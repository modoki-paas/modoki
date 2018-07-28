package main

import (
	"github.com/modoki-paas/modoki/consul_traefik"
	"github.com/docker/docker/client"
	"github.com/jmoiron/sqlx"
)

type UserControllerUtil struct {
	DB           *sqlx.DB
	DockerClient *client.Client
	Consul       *consulTraefik.Client
}
