package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/georgylobko/gym-log/internal/database"
	"github.com/georgylobko/gym-log/internal/helpers"
	"github.com/georgylobko/gym-log/internal/mappers"
)

func (apiCfg *ApiConfig) HandlerCreateMuscleGroup(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name     string `json:"name"`
		PhotoUrl string `json:"photo_url"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}

	muscleGroup, err := apiCfg.DB.CreateMuscleGroup(r.Context(), database.CreateMuscleGroupParams{
		Name:     params.Name,
		PhotoUrl: params.PhotoUrl,
	})
	if err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprintf("Could not create muscle group: %s", err))
		return
	}

	helpers.RespondWithJSON(w, 200, mappers.DatabaseMuscleGroupToMuscleGroup(muscleGroup))
}
