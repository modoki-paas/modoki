package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

const version = "0.1-beta"
const apiVersion = "1"

const versionFormat = `modoki client version: %s
API version: %s`

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf(
			versionFormat,
			version,
			apiVersion,
		)
	}

	app := cli.NewApp()
	app.Usage = "Use modoki with CLI like Docker"
	app.Version = version
	app.UsageText = "modoki [global options] command [command options] [arguments...]"
	app.Name = "modoki"

	app.Commands = []cli.Command{
		cli.Command{
			Name:        "create",
			ArgsUsage:   "[options] [iamge name] [commands...]",
			Description: "Create a new container",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name",
					Usage: "Service name(sub domain)",
				},
				cli.StringFlag{
					Name:  "entrypoint",
					Usage: "Entrypoint",
				},
				cli.StringSliceFlag{
					Name:  "env",
					Usage: "Environment variables",
				},
				cli.StringSliceFlag{
					Name:  "volumes, v",
					Usage: "Path to volumes in a container",
				},
				cli.StringFlag{
					Name:  "workdir",
					Usage: "Working directory",
				},
				cli.BoolTFlag{
					Name:  "ssl-redirect",
					Usage: "Force clients to redirec to https",
				},
			},
		},
		cli.Command{
			Name:        "start",
			ArgsUsage:   "[id or name]",
			Description: "Start a container",
		},
		cli.Command{
			Name:        "stop",
			ArgsUsage:   "[id or name]",
			Description: "Stop a container",
		},
		cli.Command{
			Name:        "remove",
			ArgsUsage:   "[options] [id or name]",
			Description: "Remove a container",

			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "force",
					Usage: "Remove if a container is running",
				},
			},
		},
		cli.Command{
			Name:        "inspect",
			ArgsUsage:   "[id or name]",
			Description: "Show inspection of a container",
		},
		cli.Command{
			Name:        "ps",
			ArgsUsage:   "[options]",
			Description: "Show a list of containers",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "json",
					Usage: "in json format",
				},
			},
		},
		cli.Command{
			Name:        "logs",
			ArgsUsage:   "[options...] [id or name]",
			Description: "Show logs of a container",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "follow, f",
				},
				cli.BoolFlag{
					Name: "stdout",
				},
				cli.BoolFlag{
					Name: "stderr",
				},
				cli.IntFlag{
					Name:  "since",
					Usage: "UNIX Time",
				},
				cli.IntFlag{
					Name:  "until",
					Usage: "UNIX Time",
				},
				cli.BoolFlag{
					Name: "timestamps",
				},
				cli.StringFlag{
					Name: "tail",
				},
			},
		},
		cli.Command{
			Name:      "cp",
			ArgsUsage: `[container_name:]source_path [container_name:]dest_path (container_name can't be used for the both)`,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("error: %s", err.Error())
	}
}
