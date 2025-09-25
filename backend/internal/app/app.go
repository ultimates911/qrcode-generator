package app

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
	"qrcodegen/config"
	"qrcodegen/internal/delivery"
	"qrcodegen/internal/delivery/http"
	"qrcodegen/internal/pkg/database"
	"qrcodegen/internal/repository/postgres"
	"qrcodegen/internal/usecase"
)

func New() *fx.App {
	return fx.New(
		fx.Provide(
			config.New,

			database.NewDBPool,

			postgres.NewRepository,

			usecase.NewUserUseCase,

			http.NewUserHandler,
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
	return fiber.New()
}