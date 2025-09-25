package delivery

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"qrcodegen/config"
	"qrcodegen/internal/delivery/http"
	"qrcodegen/internal/delivery/middleware"
)

type Router struct {
	userHandler *http.UserHandler
	cfg         *config.Config
}

func NewRouter(userHandler *http.UserHandler, cfg *config.Config) *Router {
	return &Router{userHandler: userHandler, cfg: cfg}
}

func (r *Router) Register(app *fiber.App) {
	app.Use(middleware.Recovery())
	app.Use(middleware.Logger())
	app.Use(cors.New())

	apiV1 := app.Group("/api/v1")

	apiV1.Post("/register", r.userHandler.Register)
	apiV1.Post("/login", r.userHandler.Login)
}