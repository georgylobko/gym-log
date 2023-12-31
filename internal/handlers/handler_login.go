package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/georgylobko/gym-log/internal/helpers"
	"github.com/georgylobko/gym-log/internal/mappers"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (apiCfg *ApiConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}

	user, err := apiCfg.DB.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprintf("Could not get user: %s", err))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprintf("Password is not valid: %s", err))
		return
	}

	claims := helpers.Claims{
		mappers.DatabaseUserToUser(user),
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			Issuer:    strconv.Itoa(int(user.ID)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretString := os.Getenv("JWT_SECRET")
	if secretString == "" {
		helpers.RespondWithError(w, 500, fmt.Sprintf("Something went wrong: %s", err))
		return
	}

	ss, err := token.SignedString([]byte(secretString))
	if err != nil {
		helpers.RespondWithError(w, 500, fmt.Sprintf("Something went wrong: %s", err))
		return
	}

	cookie := http.Cookie{
		Name:     "token",
		Value:    ss,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)

	helpers.RespondWithJSON(w, 200, mappers.DatabaseUserToUser(user))
}
