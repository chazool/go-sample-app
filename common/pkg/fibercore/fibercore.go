package fibercore

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func SettupFiber(IdleTimeout time.Duration) *fiber.App {

	app := fiber.New(fiber.Config{
		Prefork:               false,
		IdleTimeout:           IdleTimeout,
		DisableStartupMessage: false,
	})

	return app
}
