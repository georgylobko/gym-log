package helpers

import (
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func GetParsedToken(w http.ResponseWriter, r *http.Request) (string, error) {
	tokenCookie, err := r.Cookie("token")
	if err != nil {
		return "", err
	}
	secretString := os.Getenv("JWT_SECRET")
	if secretString == "" {
		return "", err
	}
	token, err := jwt.ParseWithClaims(tokenCookie.Value, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretString), nil
	})
	if err != nil {
		return "", err
	}

	claims := token.Claims
	issuer, err := claims.GetIssuer()
	if err != nil {
		return "", err
	}

	return issuer, nil
}
