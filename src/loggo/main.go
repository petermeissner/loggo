package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/gofiber/websocket/v2"
)

func main() {

	// configs
	configs := map[string]string{
		"template_path": "./assets/views",
		"template_ext":  ".html",
		"static_path":   "./assets/public",
	}

	// Startup checks
	// ...

	// check if template path exists
	if File_exists(configs["template_path"]) {
		// ok
	} else {
		log.Println("template_path does not exist: ", configs["template_path"])
		os.Exit(1)
	}

	// check if static path exists
	if File_exists(configs["static_path"]) {
		// ok
	} else {
		log.Println("static_path does not exist: ", configs["static_path"])
		os.Exit(1)
	}

	// initialize template engine
	engine := html.New(configs["template_path"], configs["template_ext"])
	engine.Reload(true)
	// engine.Debug(true)
	engine.Layout("embed")
	engine.Delims("{{", "}}")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// route: root, for static content
	app.Static("/", configs["static_path"])

	// route: root
	app.Get("/", func(c *fiber.Ctx) error {

		routes := app.GetRoutes()
		rout_paths := []string{}

		for i := 0; i < len(routes); i++ {
			if routes[i].Method != "HEAD" {
				rout_paths = append(rout_paths, "["+routes[i].Method+"] "+routes[i].Path)
			}
		}

		return c.Render("index", fiber.Map{
			"timestamp": time.Now().Format("2006-01-02 15:04:05"),
			"log_files": list_log_files(),
			"routes":    rout_paths,
		})

	})

	app.Get("/apis", func(c *fiber.Ctx) error {

		routes := app.GetRoutes()
		rout_paths := []string{}

		for i := 0; i < len(routes); i++ {
			if routes[i].Method != "HEAD" {
				rout_paths = append(rout_paths, "["+routes[i].Method+"] "+routes[i].Path)
			}
		}

		return c.Render("apis", fiber.Map{
			"timestamp": time.Now().Format("2006-01-02 15:04:05"),
			"log_files": list_log_files(),
			"routes":    rout_paths,
		})

	})

	// return simple list of available log files
	app.Get("/log_files", func(c *fiber.Ctx) error {

		// get list of log files
		log_files := list_log_files()

		// return log files as json
		return c.JSON(log_files)

	})

	app.Get("/ws/echo", websocket.New(func(c *websocket.Conn) {
		// access websocket via ws://localhost:3000/ws/echo

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

	app.Get("/ws/test_stream", websocket.New(func(c *websocket.Conn) {

		var (
			messageType int
			msg         []byte
			err         error
		)

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

	// start and listen
	log.Fatal(app.Listen("localhost:3000"))
}

func File_exists(s string) {
	panic("unimplemented")
}
