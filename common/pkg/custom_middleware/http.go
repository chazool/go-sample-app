package custommiddleware

import (
	"fmt"
	"runtime"
	"time"

	"github.com/chazool/go-sample-app/common/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"go.uber.org/zap"
)

const defaultStackTraceBufLengrth uint = 2048

// RequestMiddleware using add fiber middleware
func RequestMiddleware(app *fiber.App, pProfEnabled bool) {

	//Nedd a global for our zapStackTraceHandler
	app.Use(requestTimer(),
		recover.New(recover.Config{
			EnableStackTrace:  true,
			StackTraceHandler: zapStackTraceHandler,
		}),
		requestid.New(),
	)

	if pProfEnabled {
		app.Use(pprof.New())
	}

}

// requestTimer will measure how long it takes before a response is returned
func requestTimer() fiber.Handler {
	return func(c *fiber.Ctx) error {
		//start time
		start := time.Now()
		//next routes
		err := c.Next()
		// stop timer
		stop := time.Now()
		// do somthing with response
		c.Append("Server-Timing", fmt.Sprintf("app;dur=%v", stop.Sub(start).String()))
		// return stack error if exist
		return err
	}
}

func zapStackTraceHandler(c *fiber.Ctx, e interface{}) {
	buf := make([]byte, defaultStackTraceBufLengrth)
	buf = buf[:runtime.Stack(buf, false)]
	utils.Logger.Error("Recoverer middleware cought panic", zap.ByteString("stsck", buf))
}
