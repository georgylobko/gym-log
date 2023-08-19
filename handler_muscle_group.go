package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/georgylobko/gym-log/internal/database"
)

func (apiCfg *apiConfig) handlerCreateMuscleGroup(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name     string `json:"name"`
		PhotoUrl string `json:"photo_url"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}

	muscleGroup, err := apiCfg.DB.CreateMuscleGroup(r.Context(), database.CreateMuscleGroupParams{
		Name:     params.Name,
		PhotoUrl: params.PhotoUrl,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not create muscle group: %s", err))
		return
	}

	respondWithJSON(w, 200, databaseMuscleGroupToMuscleGroup(muscleGroup))
}
