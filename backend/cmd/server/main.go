package main

import (
	"os"
	"qrcodegen/internal/app"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// @title QR Code Generator API
// @BasePath /api/v1
func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	app.New().Run()
}