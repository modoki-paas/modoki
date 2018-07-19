//go:generate goagen bootstrap -d github.com/cs3238-tsuzu/modoki/design

package main

import (
	"context"

	"github.com/cs3238-tsuzu/modoki/consul_traefik"
	_ "github.com/go-sql-driver/mysql"

	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/cs3238-tsuzu/modoki/app"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/docker/docker/client"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/goadesign/goa/middleware/security/jwt"
	"github.com/jmoiron/sqlx"
)

var (
	sqlDriver        = flag.String("driver", "mysql", "SQL Driver")
	jwtPub           = flag.String("jwtpub", "/usr/local/modoki/cred/*.pub", "Glob of JWT public key")
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

	// Create service
	service := goa.New("Modoki API")

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

	consul, err := consulTraefik.NewClient("traefik", *consulHost)

	if err != nil {
		log.Fatal("error: Connecting to Zookeeper server error", err)
	}

	if ok, err := consul.HasFrontend(traefikFrontendName); err != nil {
		log.Fatal("error: zookeeper.HasFrontend error", err)
	} else if !ok {
		if err := consul.NewFrontend(traefikFrontendName, "Host:"+*publicAddr); err != nil {
			log.Fatal("error: zookeeper.NewFrontend error", err)
		}

		if err := consul.AddValueForFrontend(traefikFrontendName, "passHostHeader", true); err != nil {
			log.Fatal("error: zookeeper.AddValueForFrontend error", err)
		}

		if *https {
			if err := consul.AddValueForFrontend(traefikFrontendName, "headers", "sslredirect", true); err != nil {
				log.Fatal("error: zookeeper.AddValueForFrontend error", err)
			}
		}

		if err := consul.AddValueForFrontend(traefikFrontendName, "backend", traefikBackendName); err != nil {
			log.Fatal("error: zookeeper.AddValueForFrontend error", err)
		}
	}

	if err := consul.NewBackend(traefikBackendName, serverName, *traefikAddr); err != nil {
		log.Fatal("error: zookeeper.NewBackend error", err)
	}

	defer consul.DeleteBackend(traefikBackendName)

	keys, err := LoadJWTPublicKeys(*jwtPub)

	if err != nil {
		log.Fatal("error: Failed to load public keys", err)
	}

	// Mount middleware
	jwtMiddlware := jwt.New(jwt.NewSimpleResolver(keys), nil, app.NewJWTSecurity())

	app.UseJWTMiddleware(service, jwtMiddlware)

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

	// Start service

	if err := service.ListenAndServe(":80"); err != nil {
		service.LogError("startup", "err", err)
	}
}

// LoadJWTPublicKeys loads PEM encoded RSA public keys used to validata and decrypt the JWT.
func LoadJWTPublicKeys(path string) ([]jwt.Key, error) {
	keyFiles, err := filepath.Glob(path)
	if err != nil {
		return nil, err
	}
	keys := make([]jwt.Key, len(keyFiles))
	for i, keyFile := range keyFiles {
		pem, err := ioutil.ReadFile(keyFile)
		if err != nil {
			return nil, err
		}
		key, err := jwtgo.ParseRSAPublicKeyFromPEM([]byte(pem))
		if err != nil {
			return nil, fmt.Errorf("failed to load key %s: %s", keyFile, err)
		}
		keys[i] = key
	}
	if len(keys) == 0 {
		return nil, fmt.Errorf("couldn't load public keys for JWT security")
	}

	return keys, nil
}
