package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func main() {
	addr := flag.String("addr", ":8080", "http service address")
	debug := flag.Bool("debug", false, "run in debug mode")
	flag.Parse()

	fiberConf := fiber.Config{
		Prefork: true,
	}

	wsConf := websocket.Config{
		HandshakeTimeout: 100 * time.Second,
		Origins: []string{
			fmt.Sprintf("http://localhost%s", *addr),
			fmt.Sprintf("http://127.0.0.1%s", *addr),
		},
	}

	if *debug {
		fiberConf.Prefork = false
		wsConf.Origins = []string{"*"}
	}

	app := fiber.New(fiberConf)

	hub := NewHub()
	hub.Defaults()

	app.Use("/ws/chat", hub.Upgrade)

	app.Get("/ws/chat", websocket.New(hub.Handler, wsConf))

	app.Static("/", "./client/dist", fiber.Static{
		Compress: true,
	})
	app.Get("/*", func(ctx *fiber.Ctx) error {
		return ctx.SendFile("./client/dist/index.html", true)
	})

	go hub.Run()

	log.Fatal(app.Listen(*addr))
}
