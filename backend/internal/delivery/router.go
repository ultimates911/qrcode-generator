package delivery

import (
	"qrcodegen/config"
	_ "qrcodegen/docs"
	"qrcodegen/internal/delivery/http"
	"qrcodegen/internal/delivery/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/clode-labs/gofiber-swagger/v2"
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

	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Get("/redirect/:hash", r.linkHandler.Redirect)

	apiV1 := app.Group("/api/v1")

	apiV1.Post("/register", r.userHandler.Register)
	apiV1.Post("/login", r.userHandler.Login)

	authenticated := apiV1.Group("/", middleware.Auth(r.cfg))
	authenticated.Post("/qrcode", r.qrHandler.Generate)

	links := authenticated.Group("/links")
	links.Post("/create", r.linkHandler.CreateLink)
	links.Get("/", r.linkHandler.GetAllLinks)
	links.Get("/:id<int>", r.linkHandler.GetLink)
	links.Patch("/:id<int>", r.linkHandler.EditLink)
	links.Get("/:id<int>/download", r.linkHandler.DownloadQR)
	links.Get("/:id<int>/transitions", r.linkHandler.GetTransitionsByLink)
}
