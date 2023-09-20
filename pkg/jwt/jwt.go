package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func CreateToken(data map[string]interface{}, secret string, expiry int) (string, error) {
	exp := time.Now().Add(time.Hour * time.Duration(expiry)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"data": data, "exp": exp})

	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return t, err
}
