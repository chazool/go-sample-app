package config

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chazool/go-sample-app/common/pkg/fibercore"
	"github.com/chazool/go-sample-app/common/pkg/utils"
	"github.com/chazool/go-sample-app/common/pkg/utils/constant"
	"github.com/gofiber/fiber/v2"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type service struct {
	_             struct{}
	App           *fiber.App
	Ctx           context.Context
	Cancel        context.CancelFunc
	TraceProvider *tracesdk.TracerProvider
}

func Start() {

	appConfig := GetConfig()
	app := fibercore.SettupFiber(appConfig.ChildFiberProcessIdleTimeout)

	ctx, cancel := context.WithCancel(context.Background())

	//initialize the datadog or jaeger, opentelementry config
	traceprovider, err := appConfig.setOpentelementry(app, ctx)

	if err != nil {
		log.Panic(err.Error())
	}

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

	defer shutdown(service{App: app, Ctx: ctx, Cancel: cancel, TraceProvider: traceprovider})

}

func shutdown(service service) {
	defer service.Cancel()

	// cleanly shutdown and flush telementry when the application exits
	defer func(ctx context.Context) {
		// do not make the application hang when it is shutdown
		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := service.TraceProvider.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}(service.Ctx)

	if appConfig.Tracesink == constant.DatadogTracingSink {
		//shutdown datadog
		defer tracer.Stop()
	}

	defer utils.Logger.Sync()

	err := fibercore.Shutdown(service.App)
	if err != nil {
		utils.Logger.Fatal("Error during shudown", zap.Error(err))
	}
}
