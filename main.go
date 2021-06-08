package main

import (
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func main() {
	app := fiber.New(fiber.Config{
		Prefork: true,
	})

	hub := NewHub()

	app.Use("/ws", hub.Upgrade)

	app.Get("/ws/:roomID", websocket.New(hub.Handler, websocket.Config{
		Origins: []string{
			"http://localhost:8080",
			"http://127.0.0.1:8080",
			// "http://localhost:3000",
			// "http://127.0.0.1:3000",
		},
	}))

	app.Static("/", "./dist", fiber.Static{
		Compress: true,
	})
	app.Get("/*", func(ctx *fiber.Ctx) error {
		return ctx.SendFile("./dist/index.html", true)
	})

	go hub.Run()

	addr := flag.String("addr", ":8080", "http service address")
	flag.Parse()
	log.Fatal(app.Listen(*addr))
}
