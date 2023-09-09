package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/georgylobko/gym-log/internal/database"
	"github.com/georgylobko/gym-log/internal/helpers"
	"github.com/georgylobko/gym-log/internal/mappers"
	"golang.org/x/crypto/bcrypt"
)

func (apiCfg *ApiConfig) HandlerRegister(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name     string `json:"name"`
		Gender   string `json:"gender"`
		Role     string `json:"role"`
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

	passwordByteArr, _ := bcrypt.GenerateFromPassword([]byte(params.Password), 14)
	passwordHash := string(passwordByteArr)

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Name:     params.Name,
		Gender:   sql.NullString{String: params.Gender, Valid: true},
		Role:     params.Role,
		Email:    params.Email,
		Password: passwordHash,
	})
	if err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprintf("Could not create user: %s", err))
		return
	}

	helpers.RespondWithJSON(w, 200, mappers.DatabaseUserToUser(user))
}
