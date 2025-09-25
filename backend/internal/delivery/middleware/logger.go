package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"time"
)

func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		chainErr := c.Next()

		latency := time.Since(start)
		status := c.Response().StatusCode()
		method := c.Method()
		path := c.Path()

		event := log.Info()
		if chainErr != nil {
			if e, ok := chainErr.(*fiber.Error); ok {
				status = e.Code
			}
			event = log.Error().Err(chainErr)
		}

		if logErr, ok := c.Locals("logError").(error); ok {
			event = log.Error().Err(logErr)
		}

		event.
			Str("method", method).
			Str("path", path).
			Int("status", status).
			Dur("latency", latency).
			Str("ip", c.IP()).
			Str("user_agent", c.Get("User-Agent")).
			Msg("Request")

		return nil
	}
}