package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/georgylobko/gym-log/internal/database"
	"github.com/georgylobko/gym-log/internal/helpers"
	"github.com/georgylobko/gym-log/internal/mappers"
	"github.com/go-playground/validator/v10"
)

func (apiCfg *ApiConfig) HandlerCreateMuscleGroup(w http.ResponseWriter, r *http.Request, user mappers.User) {
	type parameters struct {
		Name     string `json:"name" validate:"required"`
		PhotoUrl string `json:"photo_url" validate:"required"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}

	validate := validator.New()

	err = validate.Struct(params)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		helpers.RespondWithError(w, 400, fmt.Sprintf("Bad user input: %s", errors))
		return
	}

	muscleGroup, err := apiCfg.DB.CreateMuscleGroup(r.Context(), database.CreateMuscleGroupParams{
		Name:     params.Name,
		PhotoUrl: params.PhotoUrl,
	})
	if err != nil {
		helpers.RespondWithError(w, 500, fmt.Sprintf("Could not create muscle group: %s", err))
		return
	}

	helpers.RespondWithJSON(w, 200, mappers.DatabaseMuscleGroupToMuscleGroup(muscleGroup))
}

func (apiCfg *ApiConfig) HandlerGetMuscleGroups(w http.ResponseWriter, r *http.Request, user mappers.User) {
	muscleGroups, err := apiCfg.DB.GetMuscleGroups(r.Context())
	if err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprintf("Could not get muscle groups: %s", err))
		return
	}

	helpers.RespondWithJSON(w, 200, mappers.DatabaseMuscleGroupsToMuscleGroups(muscleGroups))
}
