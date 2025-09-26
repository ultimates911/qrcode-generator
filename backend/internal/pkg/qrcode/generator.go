package qrcode

import (
	"bytes"
	"fmt"
	"image/png"
	"math"
	"regexp"
	"strings"

	"github.com/quickqr/gqr"
	export "github.com/quickqr/gqr/export/image"
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
		return nil, err
	}
	bg, err := normalizeHex(bgHex)
	if err != nil {
		return nil, err
	}

	gap := math.Max(0, math.Min(0.5, smoothing))

	mat, err := gqr.NewWith(
		url,
		gqr.WithErrorCorrectionLevel(gqr.ErrorCorrectionMedium),
		gqr.WithEncodingMode(gqr.EncModeAuto),
	)
	if err != nil {
		return nil, fmt.Errorf("create qr: %w", err)
	}

	img := export.NewExporter(
		export.WithFgColorHex("#"+fg),
		export.WithBgColorHex("#"+bg),
		export.WithImageSize(1024),
		export.WithQuietZone(60),
		export.WithModuleGap(gap),
	).Export(*mat)

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("encode png: %w", err)
	}
	return buf.Bytes(), nil
}
