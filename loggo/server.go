package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
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
	if file_exists(configs["template_path"]) {
		// ok
	} else {
		log.Println("template_path does not exist: ", configs["template_path"])
		os.Exit(1)
	}
	// check if static path exists
	if file_exists(configs["static_path"]) {
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

		return c.Render("index", fiber.Map{
			"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		})

	})

	app.Get("/log_files", func(c *fiber.Ctx) error {

		// get list of log files
		log_files := list_log_files()

		// return log files as json
		return c.JSON(log_files)

	})

	// start and listen
	log.Fatal(app.Listen("localhost:3000"))
}
