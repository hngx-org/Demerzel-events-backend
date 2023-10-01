package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
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

func VerifyToken(signedToken string, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error: malformed jwt token")
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("error: invalid jwt token")
}

func DecodeToken(signedToken string) (jwt.MapClaims, error) {
	parser := jwt.Parser{}
	token, _, err := parser.ParseUnverified(signedToken, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		return claims, nil
	}

	return nil, fmt.Errorf("error: invalid jwt token")
}

func VerifyFromBearer(value string) (jwt.MapClaims, error) {
	if value == "" {
		return nil, fmt.Errorf("authentication header is requires")
	}

	token := value[len("Bearer "):]
	if token == "" {
		return nil, fmt.Errorf("authentication header is required")
	}

	claims, err := VerifyToken(token, os.Getenv("JWT_SECRET"))
	if err != nil {
		return nil, err
	}

	return claims, nil
}
