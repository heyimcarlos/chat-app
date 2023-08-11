package api

import (
	"fmt"
	// "log"
	// "net/http"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/heyimcarlos/chat-app/backend/internal"
)

type Server struct {
	listenAddr string
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
	}
}

func (s *Server) FiberRun() error {
	app, err := buildServer()
	if err != nil {
		return err
	}

	fmt.Println("Server running on: ", s.listenAddr)

	// start server
	return app.Listen("0.0.0.0:" + s.listenAddr)
}

func buildServer() (*fiber.App, error) {
	// create fiber app
	app := fiber.New()

	// add middleware
	app.Use(cors.New())
	app.Use(logger.New())

	// add websocket middleware
	app.Use("/ws", func(ctx *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(ctx) {
			ctx.Locals("allowed", true)
			return ctx.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// gofiber websocket handler
	app.Get("/gofiber/:room/:name", websocket.New(internal.NewRoomStore().WsHandler))

	return app, nil
}
