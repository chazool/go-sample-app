package fibercore

import (
	"time"

	"github.com/chazool/go-sample-app/common/pkg/utils"
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

func Shutdown(app *fiber.App) error {
	utils.Logger.Info("Shutting down Fiber...")
	err := app.Shutdown()
	utils.Logger.Info("Shutdown complete")
	return err
}
