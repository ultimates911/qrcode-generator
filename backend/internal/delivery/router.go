package delivery

import (
	"qrcodegen/config"
	"qrcodegen/internal/delivery/http"
	"qrcodegen/internal/delivery/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Router struct {
	userHandler *http.UserHandler
	linkHandler *http.LinkHandler
	qrHandler   *http.QRHandler
	cfg         *config.Config
}

func NewRouter(userHandler *http.UserHandler, linkHandler *http.LinkHandler, qrHandler *http.QRHandler, cfg *config.Config) *Router {
	return &Router{
		userHandler: userHandler,
		linkHandler: linkHandler,
		qrHandler:   qrHandler,
		cfg:         cfg,
	}
}

func (r *Router) Register(app *fiber.App) {
	app.Use(middleware.Recovery())
	app.Use(middleware.Logger())
	app.Use(cors.New())

	apiV1 := app.Group("/api/v1")

	apiV1.Post("/register", r.userHandler.Register)
	apiV1.Post("/login", r.userHandler.Login)

	authenticated := apiV1.Group("/", middleware.Auth(r.cfg))
	authenticated.Post("/qrcode", r.qrHandler.Generate)

	links := authenticated.Group("/links")
	links.Post("/create", r.linkHandler.CreateLink)
	links.Get("/:id<int>", r.linkHandler.GetLink)
}