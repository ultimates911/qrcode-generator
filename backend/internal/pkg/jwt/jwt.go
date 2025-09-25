package jwt

import (
	"time"

	"qrcodegen/config"
	sqldb "qrcodegen/sqlc/generated"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func Sign(user sqldb.User, cfg *config.Config) (string, int64, error) {
	now := time.Now()
	exp := now.Add(time.Duration(cfg.JWTTTLMinutes) * time.Minute)

	claims := Claims{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   toString(user.ID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := t.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", 0, err
	}
	return signed, int64(exp.Sub(now).Seconds()), nil
}

func toString(id int64) string {
	return fmtInt(id)
}

func fmtInt(v int64) string {
	b := [20]byte{}
	i := len(b)
	neg := v < 0
	if neg {
		v = -v
	}
	for {
		i--
		b[i] = byte('0' + v%10)
		v /= 10
		if v == 0 {
			break
		}
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}
