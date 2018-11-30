//go:generate goagen bootstrap -d github.com/modoki-paas/modoki/design

package main

import (
	"flag"
	"log"

	"github.com/modoki-paas/modoki/controller"

	"github.com/docker/docker/client"
	_ "github.com/go-sql-driver/mysql"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/modoki-paas/modoki/app"
	"github.com/modoki-paas/modoki/const"
	"github.com/modoki-paas/modoki/consul_traefik"
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

	dockerClient, err := client.NewClient(*docker, *dockerAPIVersion, nil, nil)

	if err != nil {
		log.Fatal("Docker client initialization error", err)
	}

	containerUtil := &controller.ContainerControllerUtil{
		DockerClient:     dockerClient,
		DB:               db,
		Consul:           consul,
		PublicAddr:       *publicAddr,
		HTTPS:            *https,
		DockerAPIVersion: *dockerAPIVersion,
		NetworkName:      networkName,
	}
	userUtil := &controller.UserControllerUtil{
		DockerClient: dockerClient,
		DB:           db,
		Consul:       consul,
	}
	// TODO: Stop running in goroutine(export to another service)
	//go containerUtil.run(context.Background())

	// Create service
	service := goa.New("Modoki API")

	app.UseAPIKeyMiddleware(service, newAPIKeyMiddleware(db))
	app.UseJWTMiddleware(service, jwtMiddleware)

	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	// Mount controllers

	userControllerImpl := &controller.UserControllerImpl{
		userUtil,
	}

	containerControllerImpl := &controller.ContainerControllerImpl{
		containerUtil,
	}

	containerForAPIController := NewContainerForAPIController(service)
	containerForAPIController.controllerImpl = containerControllerImpl

	containerForFrontendController := NewContainerForFrontendController(service)

	containerForFrontendController.controllerImpl = containerControllerImpl

	userForAPIController := NewUserForAPIController(service)
	userForAPIController.controllerImpl = userControllerImpl

	userForFrontendController := NewUserForFrontendController(service)
	userForFrontendController.controllerImpl = userControllerImpl

	app.MountContainerForAPIController(service, containerForAPIController)
	app.MountContainerForFrontendController(service, containerForFrontendController)

	app.MountUserForAPIController(service, userForAPIController)
	app.MountUserForFrontendController(service, userForFrontendController)

	swaggerController := NewSwaggerController(service)

	app.MountSwaggerController(service, swaggerController)

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

	if _, err := db.Exec(constants.ContainerSchema); err != nil {
		log.Fatal("error: Failed to create containers table: ", err)
	}

	if _, err := db.Exec(constants.AuthorizedKeysSchema); err != nil {
		log.Fatal("error: Failed to create authorizedKeys schema: ", err)
	}

	if _, err := db.Exec(constants.APIKeysSchema); err != nil {
		log.Fatal("error: Failed to create apiKeys schema: ", err)
	}

	return db
}

func consulInit() *consulTraefik.Client {
	consul, err := consulTraefik.NewClient("traefik", *consulHost)

	if err != nil {
		log.Fatal("error: Connecting to consul server error", err)
	}

	if ok, err := consul.HasFrontend(constants.TraefikFrontendName); err != nil {
		log.Fatal("error: consul.HasFrontend error", err)
	} else if !ok {

		if err := consul.AddValueForFrontend(
			constants.TraefikFrontendName,
			"routes", "host", "rule", "Host:"+*publicAddr+";PathPrefix:/api/,/frontend/,/swagger/",
		); err != nil {
			log.Fatal("error: consul.NewFrontend error", err)
		}

		if err := consul.AddValueForFrontend(constants.TraefikFrontendName, "passHostHeader", true); err != nil {
			log.Fatal("error: consul.AddValueForFrontend error", err)
		}

		if *https {
			if err := consul.AddValueForFrontend(constants.TraefikFrontendName, "headers", "sslredirect", true); err != nil {
				log.Fatal("error: consul.AddValueForFrontend error", err)
			}
		}

		if err := consul.AddValueForFrontend(constants.TraefikFrontendName, "backend", constants.TraefikBackendName); err != nil {
			log.Fatal("error: consul.AddValueForFrontend error", err)
		}
	}

	if err := consul.NewBackend(constants.TraefikBackendName, constants.ServerName, *traefikAddr); err != nil {
		log.Fatal("error: consul.NewBackend error", err)
	}

	return consul
}

func finalize(consul *consulTraefik.Client) {
	consul.DeleteBackend(constants.TraefikBackendName)
}
