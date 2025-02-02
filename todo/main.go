package main

import (
	"log"
	"todo/db"
	"todo/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
)

func main() {
	db.Connect()
	engine := html.New("./tmplts", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Static("/","./static")
	app.Use(cors.New())
	routes.Setup(app)
    log.Fatal(app.Listen(":4444"))
}
