package handlers

import (
	"fmt"
	"net/http"

	"github.com/georgylobko/gym-log/internal/helpers"
	"github.com/golang-jwt/jwt/v5"
)

func (apiCfg *ApiConfig) HandlerSession(w http.ResponseWriter, r *http.Request) {
	tokenCookie, err := r.Cookie("token")
	if err != nil {
		helpers.RespondWithError(w, 401, "Unauthorized")
		return
	}
	token, err := jwt.ParseWithClaims(tokenCookie.Value, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		helpers.RespondWithError(w, 500, fmt.Sprintf("Something went wrong: %s", err))
		return
	}

	claims := token.Claims

	helpers.RespondWithJSON(w, 200, claims)
}
