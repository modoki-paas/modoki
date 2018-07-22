//go:generate goagen bootstrap -d github.com/cs3238-tsuzu/modoki/design

package main

import (
	"context"
	"flag"
	"log"

	"github.com/cs3238-tsuzu/modoki/app"
	"github.com/cs3238-tsuzu/modoki/consul_traefik"
	"github.com/docker/docker/client"
	_ "github.com/go-sql-driver/mysql"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/jmoiron/sqlx"
)

var (
	sqlDriver        = flag.String("driver", "mysql", "SQL Driver")
	authConfigPath   = flag.String("auth", "/usr/local/modoki/authconfig.json", "Path to auth config")
	docker           = flag.String("docker", "unix:///var/run/docker.sock", "Docker path")
	dockerAPIVersion = flag.String("docker-api", "v1.37", "Docker API version")
	sqlHost          = flag.String("db", "root:password@tcp(localhost:3306)/modoki?charset=utf8mb4&parseTime=True", "SQL")
	consulHost       = flag.String("consul", "localhost:8500", "Consul(KV)")
	traefikAddr      = flag.String("traefikAddr", "http://modoki", "Address to register on traefik")
	publicAddr       = flag.String("addr", "modoki.example.com", "API ep: modoki.example.com Service ep: *.modoki.example.com")
	networkName      = flag.String("net", "", "network for containers to join")
	https            = flag.Bool("https", true, "Enable HTTPS")
	help             = flag.Bool("help", false, "Show this")
)

// TODO: SQL retry実装

func main() {
	flag.Parse()

	if *help {
		flag.Usage()

		return
	}

	jwtMiddleware, err := initAuthMiddleware(*authConfigPath, app.NewJWTSecurity())

	if err != nil {
		log.Fatal("error: Failed to load the auth config file: ", err)
	}

	db := dbInit()
	consul := consulInit()

	defer finalize(consul)

	// Create service
	service := goa.New("Modoki API")

	app.UseJWTMiddleware(service, jwtMiddleware)

	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	dockerClient, err := client.NewClient(*docker, *dockerAPIVersion, nil, nil)

	if err != nil {
		log.Fatal("Docker client initialization error", err)
	}

	containerUtil := &ContainerControllerUtil{
		DockerClient: dockerClient,
		DB:           db,
		Consul:       consul,
	}
	userUtil := &UserControllerUtil{
		DockerClient: dockerClient,
		DB:           db,
		Consul:       consul,
	}
	go containerUtil.run(context.Background())

	// Mount "container" controller
	c := NewContainerController(service)

	c.ContainerControllerUtil = containerUtil

	app.MountContainerController(service, c)

	// Mount "user" controller
	c2 := NewUserController(service)

	c2.UserControllerUtil = userUtil

	app.MountUserController(service, c2)

	// Mount "swagger" controller
	c3 := NewSwaggerController(service)

	app.MountSwaggerController(service, c3)

	// Start service

	if err := service.ListenAndServe(":80"); err != nil {
		service.LogError("startup", "err", err)
	}
}

func dbInit() *sqlx.DB {
	db, err := sqlx.Connect(*sqlDriver, *sqlHost)

	if err != nil {
		log.Fatal("error: Connecting to SQL server error: ", err)
	}

	if _, err := db.Exec(containerSchema); err != nil {
		log.Fatal("error: Failed to create containers table: ", err)
	}

	if _, err := db.Exec(authorizedKeysSchema); err != nil {
		log.Fatal("error: Failed to create authorizedKeys schema: ", err)
	}

	return db
}

func consulInit() *consulTraefik.Client {
	consul, err := consulTraefik.NewClient("traefik", *consulHost)

	if err != nil {
		log.Fatal("error: Connecting to consul server error", err)
	}

	if ok, err := consul.HasFrontend(traefikFrontendName); err != nil {
		log.Fatal("error: consul.HasFrontend error", err)
	} else if !ok {
		if err := consul.NewFrontend(traefikFrontendName, "Host:"+*publicAddr); err != nil {
			log.Fatal("error: consul.NewFrontend error", err)
		}

		if err := consul.AddValueForFrontend(traefikFrontendName, "passHostHeader", true); err != nil {
			log.Fatal("error: consul.AddValueForFrontend error", err)
		}

		if *https {
			if err := consul.AddValueForFrontend(traefikFrontendName, "headers", "sslredirect", true); err != nil {
				log.Fatal("error: consul.AddValueForFrontend error", err)
			}
		}

		if err := consul.AddValueForFrontend(traefikFrontendName, "backend", traefikBackendName); err != nil {
			log.Fatal("error: consul.AddValueForFrontend error", err)
		}
	}

	if err := consul.NewBackend(traefikBackendName, serverName, *traefikAddr); err != nil {
		log.Fatal("error: consul.NewBackend error", err)
	}

	return consul
}

func finalize(consul *consulTraefik.Client) {
	consul.DeleteBackend(traefikBackendName)
}
