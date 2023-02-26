package config

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/chazool/go-sample-app/common/pkg/fibercore"
	"github.com/gofiber/fiber/v2"
)

type service struct {
	_      struct{}
	App    *fiber.App
	Ctx    context.Context
	Cancel context.CancelFunc
}

func Start() {

	appConfig := GetConfig()
	app := fibercore.SettupFiber(appConfig.ChildFiberProcessIdleTimeout)

	ctx, cancel := context.WithCancel(context.Background())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// fiber will listen from a defferent gorouting
	go func() {
		if err := app.Listen(":" + appConfig.SrvListenPort); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)                    // create chanel to signal begin send
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // when an interupt or termination signal is send, notif the channel

	<-c // this blocks the main thread until an interrupt is received

	defer shutdown(service{App: app, Ctx: ctx, Cancel: cancel})

}

func shutdown(service service) {
	defer service.Cancel()
}
