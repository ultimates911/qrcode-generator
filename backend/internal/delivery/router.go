package delivery

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"qrcodegen/internal/delivery/http"
	"qrcodegen/internal/delivery/middleware"
)

type Router struct {
	userHandler *http.UserHandler
}

func NewRouter(userHandler *http.UserHandler) *Router {
	return &Router{userHandler: userHandler}
}

func (r *Router) Register(app *fiber.App) {
	app.Use(middleware.Recovery())
	app.Use(middleware.Logger())
	app.Use(cors.New())

	apiV1 := app.Group("/api/v1")

	userRoutes := apiV1.Group("/users")
	userRoutes.Post("/register", r.userHandler.Register)
}