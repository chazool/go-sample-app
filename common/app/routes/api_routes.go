package routes

import (
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/chazool/go-sample-app/common/app/routes/handler"
	"github.com/chazool/go-sample-app/common/pkg/utils/constant"
	"github.com/gofiber/fiber/v2"
)

func APIRoutes(app *fiber.App, livez, readyz, testConnection func(*fiber.Ctx) error) fiber.Router {

	//document routes
	app.Get("/docs", handler.DocumentHandler(constant.Doc))
	// static document route
	app.Get("/docs", handler.DocumentHandler(constant.Static))
	// swagger route
	app.Get("/docs/*", swagger.HandlerDefault)

	route := app.Group("/app/v1")
	route.Get("livez", livez)
	route.Get("readyz", readyz)
	route.Get("readyz", testConnection)

	return route
}
