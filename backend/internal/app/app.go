package app

import (
	"context"

	"qrcodegen/config"
	"qrcodegen/internal/delivery"
	"qrcodegen/internal/delivery/http"
	"qrcodegen/internal/pkg/database"
	"qrcodegen/internal/pkg/geo"
	"qrcodegen/internal/repository/postgres"
	"qrcodegen/internal/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

func New() *fx.App {
	return fx.New(
		fx.Provide(
			config.New,

			database.NewDBPool,

			postgres.NewRepository,

			geo.NewGeoResolver,

			usecase.NewUserUseCase,
			usecase.NewLinkUseCase,
			usecase.NewQRUseCase,

			http.NewUserHandler,
			http.NewLinkHandler,
			http.NewQRHandler,

			delivery.NewRouter,

			validator.New,
			NewFiberApp,
		),
		fx.Invoke(
			func(lifecycle fx.Lifecycle, app *fiber.App, router *delivery.Router, cfg *config.Config) {
				lifecycle.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						router.Register(app)

						log.Info().Msgf("Starting server on %s", cfg.HTTPServerAddress)
						go func() {
							if err := app.Listen(cfg.HTTPServerAddress); err != nil {
								log.Fatal().Err(err).Msg("Failed to start server")
							}
						}()
						return nil
					},
					OnStop: func(ctx context.Context) error {
						log.Info().Msg("Stopping server")
						return app.Shutdown()
					},
				})
			},
		),
	)
}

func NewFiberApp() *fiber.App {
	app := fiber.New(fiber.Config{
		EnableTrustedProxyCheck: true,
		TrustedProxies:          []string{"172.16.0.0/12", "192.168.0.0/16", "10.0.0.0/8"},
		ProxyHeader:             fiber.HeaderXForwardedFor,
	})
	return app
}
