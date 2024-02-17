package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	util "github.com/petermeissner/loggo/src/util"
	"github.com/urfave/cli/v2"
)

func main() {

	// global vars
	var log_path string
	var pattern string
	var port string

	// logger
	var log = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	// cli application
	cli_app := &cli.App{

		// name
		Name:    "loggo-serve",
		Version: "0.0.1",

		// description
		Usage: "Serve log files over websocket",

		// Default Action
		Action: func(cli_context *cli.Context) error {
			cli.ShowAppHelp(cli_context)
			println("\n--------------------------------------------------------------------\n")
			cli.ShowCommandHelpAndExit(cli_context, "serve", 1)
			return nil
		},

		Commands: []*cli.Command{

			// Command: serve
			{
				Name:  "serve",
				Usage: "Serve log files in folder over websocket",

				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "logpath",
						Aliases:  []string{"l"},
						Usage:    "path to log files to server",
						Required: true,
					},

					&cli.StringFlag{
						Name:    "pattern",
						Aliases: []string{"r"},
						Value:   ".*log$",
						Usage:   "path to log files to server",
					},

					&cli.StringFlag{
						Name:    "simulate",
						Aliases: []string{"s"},
						Value:   "false",
						Usage:   "run but without actually starting the server",
					},

					&cli.StringFlag{
						Name:    "port",
						Aliases: []string{"p"},
						Value:   "3000",
						Usage:   "port to run server on",
					},
				},

				Action: func(cli_context *cli.Context) error {

					// store parameter in global variables
					log_path = cli_context.String("logpath")
					pattern = cli_context.String("pattern")
					port = cli_context.String("port")

					// check if path exists
					if !util.File_exists(log_path) {
						log.Fatal("logpath does not exist: '" + log_path + "'. Use '--logpath ...' to configure.")
					}

					// current log files
					log_files := util.File_list(log_path, ".*log$")

					// report on log path and log files found
					log.Println("LogPath: " + log_path)
					for i := 0; i < len(log_files); i++ {
						log.Println("  LogFile:   " + log_files[i])
					}

					// Server
					// - if simulate is set to true, do not start the server
					// - configure server
					// - define routes
					// - start server
					if !cli_context.Bool("simulate") {

						// server instance
						server_app := fiber.New(fiber.Config{})

						// route:
						// - /
						server_app.Get("/", func(server_context *fiber.Ctx) error {

							routes := server_app.GetRoutes()
							rout_paths := []string{}

							for i := 0; i < len(routes); i++ {
								if routes[i].Method != "HEAD" {
									rout_paths = append(rout_paths, "["+routes[i].Method+"] "+routes[i].Path)
								}
							}

							return server_context.JSON(fiber.Map{
								"timestamp": time.Now().Format("2006-01-02 15:04:05"),
								"routes":    rout_paths,
								"whoami":    "loggo-serve",
								"version":   cli_context.App.Version,
							})

						})

						// route:
						// - /log_files
						server_app.Get("/log_files", func(server_context *fiber.Ctx) error {

							return server_context.JSON(fiber.Map{
								"timestamp": time.Now().Format("2006-01-02 15:04:05"),
								"log_files": util.File_list(log_path, pattern),
							})

						})

						// route:
						// - /ws/echo
						// - access websocket via ws://localhost:3000/ws/echo
						server_app.Get("/ws/echo", websocket.New(func(c *websocket.Conn) {

							var (
								messageType int
								msg         []byte
								err         error
							)

							for {
								//  Read
								messageType, msg, err = c.ReadMessage()
								if err != nil {
									log.Println("read:", err)
									break
								}
								log.Printf("type: %d, msg: %s", messageType, msg)

								// Write
								msg = []byte("Echo from server @ " + time.Now().Format("2006-01-02 15:04:05") + " : " + string(msg))
								err = c.WriteMessage(messageType, []byte(msg))
								if err != nil {
									log.Println("write:", err)
									break
								}
							}
						}))

						// route:
						// - /ws/test_stream
						server_app.Get("/ws/test_stream", websocket.New(func(c *websocket.Conn) {

							// route specific variables
							var messageType int
							var msg []byte
							var err error

							messageType = 1
							for {
								// wait for 1 second
								time.Sleep(1 * time.Second)

								// Write
								msg = []byte("Echo from server @ " + time.Now().Format("2006-01-02 15:04:05"))
								err = c.WriteMessage(messageType, []byte(msg))
								if err != nil {
									log.Println("write:", err)
									break
								}
							}
						}))

						// route:
						// - /ws/log/<filename>
						server_app.Get("/ws/log/:filename", websocket.New(func(c *websocket.Conn) {

							// report connection
							log.Println("Incoming connection: " + c.Conn.RemoteAddr().Network() + ", " + c.Conn.RemoteAddr().String() + ", to route: /ws/log/" + c.Params("filename"))

							// route specific variables
							var messageType int
							var msg []byte
							var err error

							// open file
							file, err := os.Open("./logs/" + c.Params("filename"))
							if err != nil {
								log.Fatal(err)
							}
							defer file.Close()

							// file reader
							scanner := bufio.NewScanner(file)
							buf := make([]byte, 1e9)
							scanner.Buffer(buf, 1e9)

							// read file line by line
							i := 0
							messageType = 1
							for scanner.Scan() {
								scanner.Text()

								// Write
								msg = []byte(scanner.Text())
								err = c.WriteMessage(messageType, []byte(strconv.Itoa(i)+string(msg)))
								if err != nil {
									log.Println("write:", err)
									break
								}
							}

							if err := scanner.Err(); err != nil {
								log.Fatal(err)
							}

							for {
								time.Sleep(10 * time.Second)
								err = c.WriteMessage(messageType, []byte("Standby: "+time.Now().Format("2006-01-02 15:04:05")))
								if err != nil {
									log.Println("write:", err)
									break
								}
							}

						},
							websocket.Config{EnableCompression: true}))

						// start server instance
						log.Fatal(server_app.Listen("localhost:" + port))
					}

					return nil
				},
			},

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
		},
	}

	if err := cli_app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
