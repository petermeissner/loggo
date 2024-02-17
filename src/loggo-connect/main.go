package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/petermeissner/loggo/src/util"
	"github.com/urfave/cli/v2"
)

func main() {

	// cli application
	cli_app := &cli.App{

		// name
		Name:    "loggo-connect",
		Version: "0.0.1",

		// description
		Usage: "Read log files from loggo-serve instance over websocket",

		// Default Action
		Action: func(cli_context *cli.Context) error {

			// show main help
			cli.ShowAppHelp(cli_context)

			// show help for ping
			println("\n--------------------------------------------------------------------\n")
			cli.ShowCommandHelp(cli_context, "ping")

			// show help for read
			println("\n--------------------------------------------------------------------\n")
			cli.ShowCommandHelp(cli_context, "read")
			return nil
		},

		// Sub commands (actions) available
		Commands: []*cli.Command{

			// Command: version
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "print version",

				Action: func(cli_context *cli.Context) error {
					cli.ShowVersion(cli_context)
					return nil
				},
			},

			// Command: ping
			{
				Name:  "ping",
				Usage: "Ping server",

				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "host",
						Value: "localhost",
						Usage: "host address/name to connect to",
					},

					&cli.StringFlag{
						Name:    "port",
						Aliases: []string{"p"},
						Value:   "3000",
						Usage:   "host port to connect to",
					},
				},

				Action: func(cli_context *cli.Context) error {

					// setup logger
					logger := util.Setup_logger("info", "text")

					// retrieve general info from server
					logger.Info(
						"Connecting to server",
						"host", cli_context.String("host"),
						"port", cli_context.String("port"),
						"url", cli_context.String("host")+":"+cli_context.String("port"))
					resp, err := http.Get("http://" + cli_context.String("host") + ":" + cli_context.String("port"))
					if err != nil {
						logger.Error(err.Error())
					}
					defer resp.Body.Close()

					// response info
					resp_dump, err := httputil.DumpResponse(resp, true)
					if err != nil {
						logger.Error(err.Error())
					}
					for _, line := range strings.Split(string(resp_dump), "\n") {
						logger.Info("  " + line)
					}

					// check server is loggo-serve instance
					var loggo_serve_info util.Loggo_serve_info
					json.NewDecoder(resp.Body).Decode(&loggo_serve_info)
					if loggo_serve_info.Whoami != "loggo-serve" {
						logger.Error("Server is not a loggo-serve instance")
					}
					logger.Info("Server info:", "whoami", loggo_serve_info.Whoami, "version", loggo_serve_info.Version)

					// connect to server via websocket
					host_port := cli_context.String("host") + ":" + cli_context.String("port")
					ws_url := url.URL{Scheme: "ws", Host: host_port, Path: "/ws/test_stream"}
					logger.Info("Connecting to server at: " + ws_url.String())
					conn, _, err := websocket.DefaultDialer.Dial(ws_url.String(), nil)
					if err != nil {
						log.Fatal("dial:", err)
					}
					defer conn.Close()
					done := make(chan struct{})
					defer close(done)

					// receive data
					logger.Info("Receiving messages from server")
					for i := 1; i <= 3; i++ {
						_, msg, _ := conn.ReadMessage()
						logger.Info(string(msg), "i", i)
					}

					return nil
				},
			},

			// Command: read
			{
				Name:  "read",
				Usage: "Read log file from server",

				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "host",
						Value: "localhost",
						Usage: "host address/name to connect to",
					},

					&cli.StringFlag{
						Name:    "port",
						Aliases: []string{"p"},
						Value:   "3000",
						Usage:   "host port to connect to",
					},
				},

				Action: func(cli_context *cli.Context) error {

					// setup logger
					logger := util.Setup_logger("info", "text")

					// connect to server via websocket
					host_port := cli_context.String("host") + ":" + cli_context.String("port")
					ws_url := url.URL{Scheme: "ws", Host: host_port, Path: "/ws/log/act.log"}
					logger.Info("Connecting to server at: " + ws_url.String())
					conn, _, err := websocket.DefaultDialer.Dial(ws_url.String(), nil)
					if err != nil {
						log.Fatal("dial:", err)
					}
					defer conn.Close()
					done := make(chan struct{})
					defer close(done)

					// receive data
					i := 1
					logger.Info("Receiving messages from server")
					for {
						_, msg, _ := conn.ReadMessage()
						logger.Info(string(msg), "i", i)
						i++
					}

					return nil
				},
			},
		},
	}

	if err := cli_app.Run(os.Args); err != nil {
		// setup logger
		logger := util.Setup_logger("info", "text")

		logger.Error(err.Error())
	}
}
