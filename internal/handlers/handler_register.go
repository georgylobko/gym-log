package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/georgylobko/gym-log/internal/database"
	"github.com/georgylobko/gym-log/internal/helpers"
	"github.com/georgylobko/gym-log/internal/mappers"
	"golang.org/x/crypto/bcrypt"
)

var validate *validator.Validate

func (apiCfg *ApiConfig) HandlerRegister(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name     string `json:"name" validate:"required,min=3,max=12"`
		Gender   string `json:"gender" validate:"required"`
		Role     string `json:"role" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}

	validate = validator.New()
	err = validate.Struct(params)

	if err != nil {
		helpers.RespondWithError(w, 400, err.Error())
		return
	}

	fmt.Println(err)

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
