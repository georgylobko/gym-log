package helpers

import (
	"net/http"
	"os"

	"github.com/georgylobko/gym-log/internal/mappers"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	User mappers.User `json:"user"`
	jwt.RegisteredClaims
}

func GetParsedToken(w http.ResponseWriter, r *http.Request) (mappers.User, error) {
	tokenCookie, err := r.Cookie("token")
	if err != nil {
		return mappers.User{}, err
	}
	secretString := os.Getenv("JWT_SECRET")
	if secretString == "" {
		return mappers.User{}, err
	}

	token, err := jwt.ParseWithClaims(tokenCookie.Value, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretString), nil
	})
	if err != nil {
		return mappers.User{}, err
	}

	claims := token.Claims.(*Claims)

	return claims.User, nil
}
