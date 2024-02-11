package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/petermeissner/loggo/src/util"
	"github.com/urfave/cli/v2"
	"golang.org/x/net/websocket"
)

func main() {

	// logger
	var log = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	// cli application
	cli_app := &cli.App{

		// name
		Name:    "loggo-connect",
		Version: "0.0.1",

		// description
		Usage: "Read log files from loggo-serve instance over websocket",

		// Default Action
		Action: func(cli_context *cli.Context) error {
			cli.ShowAppHelp(cli_context)
			println("\n--------------------------------------------------------------------\n")
			cli.ShowCommandHelp(cli_context, "ping")
			return nil
		},

		Commands: []*cli.Command{
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
					// retrieve general info from server
					log.Println("Connecting to server at: ", cli_context.String("host")+":"+cli_context.String("port"))
					resp, err := http.Get("http://" + cli_context.String("host") + ":" + cli_context.String("port"))
					if err != nil {
						log.Fatal(err)
					}
					defer resp.Body.Close()

					// response info
					resp_dump, err := httputil.DumpResponse(resp, true)
					if err != nil {
						log.Fatal(err)
					}
					for _, line := range strings.Split(string(resp_dump), "\n") {
						log.Println("  " + line)
					}

					// check server is loggo-serve instance
					var loggo_serve_info util.Loggo_serve_info
					json.NewDecoder(resp.Body).Decode(&loggo_serve_info)
					if loggo_serve_info.Whoami != "loggo-serve" {
						log.Fatal("Server is not a loggo-serve instance")
					}
					log.Printf("Server info: %s %s", loggo_serve_info.Whoami, loggo_serve_info.Version)

					// connect to server via websocket
					ws_url := "ws://" + cli_context.String("host") + ":" + cli_context.String("port") + "/ws/test_stream"
					ws_protocol := "ws"
					ws_origin := "http://localhost"
					log.Println("Connecting to server at: ", ws_url, ", protocol: ", ws_protocol, ", origin: ", ws_origin)
					conn, err := websocket.Dial(ws_url, ws_protocol, ws_origin)
					if err != nil {
						log.Fatal(err)
					}

					// receive data
					var message string
					log.Println("Receiving messages from server")
					for i := 1; i <= 3; i++ {
						websocket.Message.Receive(conn, &message)
						log.Printf("Message %d: "+message, i)
					}

					return nil
				},
			},

			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "print version",

				Action: func(cli_context *cli.Context) error {
					cli.ShowVersion(cli_context)
					return nil
				},
			},
		},
	}

	if err := cli_app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
