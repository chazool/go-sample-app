package custommiddleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// crosMiddleware provide Fiber's built-in middewares
// see: https://docs.gofiber.io/v1.x/api/middleware/
func CorsMiddleware(app *fiber.App) {
	app.Use(cors.New())
}
