package qrcode

import (
	"bytes"
	"fmt"
	"image/png"
	"regexp"
	"strings"

	"github.com/quickqr/gqr"
	export "github.com/quickqr/gqr/export/image"
	"github.com/quickqr/gqr/export/image/shapes"
)

func normalizeHex(s string) (string, error) {
	s = strings.TrimSpace(strings.TrimPrefix(s, "#"))
	if ok, _ := regexp.MatchString("^[0-9a-fA-F]{6}$", s); !ok {
		return "", fmt.Errorf("invalid hex color: %q", s)
	}
	return strings.ToUpper(s), nil
}

func GeneratePNG(url, colorHex, bgHex string, smoothing float64) ([]byte, error) {
	fg, err := normalizeHex(colorHex)
	if err != nil {
		fg = "5EC8FF"
	}

	bg, err := normalizeHex(bgHex)
	if err != nil {
		bg = "FFFFFF"
	}

	qr, err := gqr.NewWith(
		url,
		gqr.WithErrorCorrectionLevel(gqr.ErrorCorrectionHighest),
		gqr.WithEncodingMode(gqr.EncModeAuto),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create qr matrix: %w", err)
	}

	exp := export.NewExporter(
		export.WithImageSize(1024),
		export.WithQuietZone(48),
		export.WithModuleGap(0.14),
		export.WithBgColorHex("#"+bg),

		export.WithModuleShape(shapes.RoundedModuleShape(0.25, true)),
		export.WithFinderShape(shapes.RoundedFinderShape(0.5)),

		export.WithGradient(
			export.GradientDirectionLTR,
			export.ParseFromHex("#"+fg),
			export.ParseFromHex("#"+fg),
		),
	)

	img := exp.Export(*qr)

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("failed to encode png: %w", err)
	}
	return buf.Bytes(), nil
}
