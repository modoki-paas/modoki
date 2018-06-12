package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"

	modoki "github.com/cs3238-tsuzu/modoki/client"
	"github.com/goadesign/goa/client"
	"github.com/k0kubun/pp"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

const version = "0.1-beta"
const apiVersion = "1"

const versionFormat = `modoki client version: %s
API version: %s`

type uploadedType struct {
	SrcPath string
	DstPath string
}

type configType struct {
	Token    string
	Scheme   string
	Host     string
	Uploaded map[int]uploadedType
}

func main() {
	httpClient := &http.Client{}

	doer := client.HTTPClientDoer(httpClient)

	modokiClient := modoki.New(doer)

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf(
			versionFormat,
			version,
			apiVersion,
		)
	}

	app := cli.NewApp()
	app.Usage = "Use modoki on CLI like Docker"
	app.Version = version
	app.UsageText = "modoki [global options] command [command options] [arguments...]"
	app.Name = "modoki"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config",
			Usage:  "Path to a config file",
			Value:  "~/.modoki.config",
			EnvVar: "MODOKI_CONFIG",
		},
	}

	var config configType
	var configPath string

	app.Before = func(ctx *cli.Context) error {
		p := ctx.String("config")

		if strings.HasPrefix(p, "~/") {
			home, err := getHomeDir()

			if err != nil {
				log.Fatal("Failed to get home directory: ", err)
			}

			p = home + p[1:]
		}

		configPath = p

		fp, err := os.Open(p)

		if err == nil {
			d := json.NewDecoder(fp)

			if err := d.Decode(&config); err != nil {
				return errors.Wrap(err, "Invalid config format")
			}

			if config.Uploaded == nil {
				config.Uploaded = map[int]uploadedType{}
			}

			if config.Host == "" {
				config.Host = "modoki.tsuzu.xyz"
			}

			if config.Scheme == "" {
				config.Scheme = "https"
			}

			modokiClient.Scheme = config.Scheme
			modokiClient.Host = config.Host
			modokiClient.SetJWTSigner(newJWTSigner(config.Token))

			return nil
		}

		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	app.Commands = []cli.Command{
		cli.Command{
			Name:      "create",
			ArgsUsage: "[options] [iamge name] [commands...]",
			Usage:     "Create a new container",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name",
					Usage: "Service name(sub domain)",
				},
				cli.StringSliceFlag{
					Name:  "entrypoint",
					Usage: "Entrypoint",
				},
				cli.StringSliceFlag{
					Name:  "env, e",
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
				cli.BoolFlag{
					Name:  "dotenv",
					Usage: "Use .env for environment variables",
				},
			},

			Action: func(ctx *cli.Context) error {
				if ctx.NArg() < 1 {
					return errors.New("Image name is not specified")
				}
				if ctx.String("name") == "" {
					return errors.New("--name is not specified")
				}

				image := ctx.Args()[0]
				cmd := ctx.Args()[1:]

				sslRedirect := ctx.BoolT("ssl-redirect")
				var workDir *string
				if s := ctx.String("workdir"); len(s) != 0 {
					workDir = &s
				}

				envs := ctx.StringSlice("env")

				if ctx.Bool("dotenv") {
					dotenvvMap, err := godotenv.Read(".env")
					if err != nil {
						return errors.Wrap(err, ".env error")
					}

					for k, v := range dotenvvMap {
						envs = append(envs, k+"="+v)
					}
				}

				resp, err := modokiClient.CreateContainer(context.Background(), modoki.CreateContainerPath(), image, ctx.String("name"), cmd, ctx.StringSlice("entrypoint"), envs, &sslRedirect, ctx.StringSlice("volumes"), workDir)

				if err != nil {
					return err
				}

				if resp.StatusCode != http.StatusOK {
					res, err := modokiClient.DecodeErrorResponse(resp)

					if err != nil {
						return errors.Wrap(err, resp.Status)
					}

					return errors.Wrap(res, resp.Status)
				}

				res, err := modokiClient.DecodeGoaContainerCreateResults(resp)

				if err != nil {
					return err
				}

				fmt.Println("ID:", res.ID)
				fmt.Println("Endpoints:", strings.Join(res.Endpoints, ", "))

				return nil
			},
		},
		cli.Command{
			Name:      "start",
			ArgsUsage: "[id or name]",
			Usage:     "Start a container",
			Action: func(ctx *cli.Context) error {
				if ctx.NArg() < 1 {
					return errors.New("ID or name is not specified")
				}

				resp, err := modokiClient.StartContainer(context.Background(), modoki.StartContainerPath(), ctx.Args()[0])

				if err != nil {
					return err
				}

				switch resp.StatusCode {
				case http.StatusNoContent:
					return nil
				case http.StatusNotFound:
					return errors.New("No such container")
				default:
					res, err := modokiClient.DecodeErrorResponse(resp)

					if err != nil {
						return errors.Wrap(err, resp.Status)
					}

					return errors.Wrap(res, resp.Status)
				}
			},
		},
		cli.Command{
			Name:      "stop",
			ArgsUsage: "[id or name]",
			Usage:     "Stop a container",
			Action: func(ctx *cli.Context) error {
				if ctx.NArg() < 1 {
					return errors.New("ID or name is not specified")
				}

				resp, err := modokiClient.StopContainer(context.Background(), modoki.StopContainerPath(), ctx.Args()[0])

				if err != nil {
					return err
				}

				switch resp.StatusCode {
				case http.StatusNoContent:
					return nil
				case http.StatusNotFound:
					return errors.New("No such container")
				default:
					res, err := modokiClient.DecodeErrorResponse(resp)

					if err != nil {
						return errors.Wrap(err, resp.Status)
					}

					return errors.Wrap(res, resp.Status)
				}
			},
		},
		cli.Command{
			Name:      "remove",
			Aliases:   []string{"rm"},
			ArgsUsage: "[options] [id or name]",
			Usage:     "Remove a container",

			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "force, f",
					Usage: "Remove if a container is running",
				},
			},
			Action: func(ctx *cli.Context) error {
				if ctx.NArg() < 1 {
					return errors.New("ID or name is not specified")
				}

				resp, err := modokiClient.RemoveContainer(context.Background(), modoki.RemoveContainerPath(), ctx.Bool("force"), ctx.Args()[0])

				if err != nil {
					return err
				}

				switch resp.StatusCode {
				case http.StatusNoContent:
					return nil
				case http.StatusNotFound:
					return errors.New("No such container")
				default:
					res, err := modokiClient.DecodeErrorResponse(resp)

					if err != nil {
						return errors.Wrap(err, resp.Status)
					}

					return errors.Wrap(res, resp.Status)
				}
			},
		},
		cli.Command{
			Name:      "inspect",
			ArgsUsage: "[id or name]",
			Usage:     "Show the inspection of a container",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "json",
					Usage: "in json format",
				},
			},
			Action: func(ctx *cli.Context) error {
				if ctx.NArg() < 1 {
					return errors.New("ID or name is not specified")
				}

				resp, err := modokiClient.InspectContainer(context.Background(), modoki.InspectContainerPath(), ctx.Args()[0])

				if err != nil {
					return err
				}

				switch resp.StatusCode {
				case http.StatusOK:
					res, err := modokiClient.DecodeGoaContainerInspect(resp)

					if err != nil {
						return err
					}

					if ctx.Bool("json") {
						json.NewEncoder(os.Stdout).Encode(res)
					} else {
						pp.Println(res)
					}

					return nil
				case http.StatusNotFound:
					return errors.New("No such container")
				default:
					res, err := modokiClient.DecodeErrorResponse(resp)

					if err != nil {
						return errors.Wrap(err, resp.Status)
					}

					return errors.Wrap(res, resp.Status)
				}
			},
		},
		cli.Command{
			Name:      "ps",
			ArgsUsage: "[options]",
			Usage:     "Show a list of containers",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "json",
					Usage: "in json format",
				},
			},
			Action: func(ctx *cli.Context) error {
				resp, err := modokiClient.ListContainer(context.Background(), modoki.ListContainerPath())

				if err != nil {
					return err
				}

				switch resp.StatusCode {
				case http.StatusOK:
					res, err := modokiClient.DecodeGoaContainerListEachCollection(resp)

					if err != nil {
						return err
					}

					if ctx.Bool("json") {
						json.NewEncoder(os.Stdout).Encode(res)

						return nil
					}

					table := tablewriter.NewWriter(os.Stdout)
					table.SetBorder(true)
					table.SetHeader([]string{"Name", "ID", "Image", "Status", "Command/Msg"})

					for i := range res {
						table.Append([]string{
							res[i].Name,
							strconv.Itoa(res[i].ID),
							res[i].Image,
							res[i].Status,
							res[i].Command,
						})
					}

					table.Render()

					return nil
				default:
					res, err := modokiClient.DecodeErrorResponse(resp)

					if err != nil {
						return errors.Wrap(err, resp.Status)
					}

					return errors.Wrap(res, resp.Status)
				}
			},
		},
		cli.Command{
			Name:      "logs",
			ArgsUsage: "[options...] [id or name]",
			Usage:     "Show logs of a container",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "follow, f",
				},
				cli.BoolTFlag{
					Name: "stdout",
				},
				cli.BoolTFlag{
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
			Action: func(ctx *cli.Context) error {
				if ctx.NArg() < 1 {
					return errors.New("Image name is not specified")
				}

				var follow, timestamps *bool
				var stdout, stderr bool
				var since, until *time.Time
				var tail *string
				if ctx.IsSet("follow") {
					b := ctx.Bool("follow")
					follow = &b
				}
				stdout = ctx.Bool("stdout")
				stderr = ctx.Bool("stderr")

				if ctx.IsSet("timestamps") {
					b := ctx.Bool("timestamps")
					timestamps = &b
				}
				if ctx.IsSet("since") {
					b := time.Unix(int64(ctx.Int("since")), 0)
					since = &b
				}
				if ctx.IsSet("until") {
					b := time.Unix(int64(ctx.Int("until")), 0)
					until = &b
				}
				if ctx.IsSet("tail") {
					b := ctx.String("tail")
					tail = &b
				}

				prevScheme := modokiClient.Scheme

				var scheme string
				switch prevScheme {
				case "http":
					scheme = "ws"
				case "https":
					scheme = "wss"
				}
				modokiClient.Scheme = scheme
				defer func() {
					modokiClient.Scheme = prevScheme
				}()

				conn, err := modokiLogsContainer(modokiClient, context.Background(), modoki.LogsContainerPath(), ctx.Args()[0], follow, since, &stderr, &stdout, tail, timestamps, until)

				if err != nil {
					return err
				}

				defer conn.Close()
				io.Copy(os.Stdout, conn)

				return nil
			},
		},
		cli.Command{
			Name:        "cp",
			ArgsUsage:   `(container id/name:)source_path (container id/name:)dest_path`,
			Usage:       "Upload or download files",
			Description: "You cannot set 'container id/name' to both parameters",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "archive, a",
					Usage: "Archive mode (copy all uid/gid information)",
				},
				cli.BoolFlag{
					Name:  "follow-link, L",
					Usage: "Always follow symbol link in SRC_PATH",
				},
			},
			Action: func(ctx *cli.Context) error {
				copyUIDGID := ctx.Bool("archive")
				followLink := ctx.Bool("follow-link")

				if ctx.NArg() != 2 {
					cli.ShowSubcommandHelp(ctx)

					return errors.New("The source and destination paths must be specified")
				}
				fromContainer, from := splitCpArg(ctx.Args()[0])
				toContainer, to := splitCpArg(ctx.Args()[1])

				if fromContainer != "" && toContainer != "" {
					return errors.New("Copying between containers is not supported")
				}

				if fromContainer == "" && toContainer == "" {
					return errors.New("Use 'cp' command instead")
				}

				if fromContainer != "" {
					return copyFromContainer(context.Background(), modokiClient, cpConfig{
						followLink: followLink,
						copyUIDGID: copyUIDGID,
						sourcePath: from,
						destPath:   to,
						container:  fromContainer,
					})
				}

				return copyToContainer(context.Background(), modokiClient, cpConfig{
					followLink: followLink,
					copyUIDGID: copyUIDGID,
					sourcePath: from,
					destPath:   to,
					container:  toContainer,
				})

			},
		},

		cli.Command{
			Name:  "config",
			Usage: "Change the config",
			Subcommands: []cli.Command{
				cli.Command{
					Name:  "signin",
					Usage: "Set token in the config file",
					Action: func(ctx *cli.Context) error {
						fmt.Print("Token: ")
						var token string
						fmt.Scan(&token)
						config.Token = token

						fmt.Println("OK")

						return nil
					},
				},
				cli.Command{
					Name:        "endpoint",
					Usage:       "Set scheme and host in the config file",
					Description: "Websocket API: http->ws, https->wss",
					ArgsUsage:   " [scheme(http/https)] [host]",
					Action: func(ctx *cli.Context) error {
						if ctx.NArg() < 2 {
							return cli.ShowSubcommandHelp(ctx)
						}

						config.Scheme = ctx.Args()[0]
						config.Host = ctx.Args()[1]

						return nil
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	if configPath != "" {
		b, _ := json.Marshal(config)

		if err := ioutil.WriteFile(configPath, b, 0660); err != nil {
			log.Fatal(err)
		}
	}
}
