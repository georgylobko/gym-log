package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/georgylobko/gym-log/internal/database"
	"github.com/georgylobko/gym-log/internal/helpers"
	"github.com/georgylobko/gym-log/internal/mappers"
	"github.com/go-chi/chi"
)

func (apiCfg *ApiConfig) HandlerCreateExercise(w http.ResponseWriter, r *http.Request, userID string) {
	type parameters struct {
		Name            string `json:"name"`
		PhotoUrl        string `json:"photo_url"`
		MuscleGroupsIds []int  `json:"muscle_groups_ids"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}

	exercise, err := apiCfg.DB.CreateExercise(r.Context(), database.CreateExerciseParams{
		Name:     params.Name,
		PhotoUrl: params.PhotoUrl,
	})
	if err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprintf("Could not create the db entity: %s", err))
		return
	}

	for _, muscleGroupId := range params.MuscleGroupsIds {
		_, err := apiCfg.DB.CreateExerciseMuscleGroup(r.Context(), database.CreateExerciseMuscleGroupParams{
			ExerciseID:    exercise.ID,
			MuscleGroupID: int32(muscleGroupId),
		})
		if err != nil {
			helpers.RespondWithError(w, 400, fmt.Sprintf("Could not create the db entity: %s", err))
			return
		}
	}

	helpers.RespondWithJSON(w, 200, struct{}{})
}

func (apiCfg *ApiConfig) HandlerGetExercise(w http.ResponseWriter, r *http.Request, userID string) {
	exerciseIDStr := chi.URLParam(r, "exerciseID")
	exerciseID, err := strconv.Atoi(exerciseIDStr)
	if err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprintf("Could not parse: %s", err))
		return
	}

	exercise, err := apiCfg.DB.GetExerciseById(r.Context(), int32(exerciseID))
	if err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprintf("Could not get exercise: %s", err))
		return
	}
	muscleGroups, err := apiCfg.DB.GetMuscleGroupsByExercise(r.Context(), int32(exerciseID))
	if err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprintf("Could not get muscle groups: %s", err))
		return
	}

	helpers.RespondWithJSON(w, 200, mappers.DatabaseExerciseToExercise(exercise, muscleGroups))
}

func (apiCfg *ApiConfig) HandlerGetExercises(w http.ResponseWriter, r *http.Request, userID string) {
	exercises, err := apiCfg.DB.GetExircises(r.Context())
	if err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprintf("Could not get exercise: %s", err))
		return
	}

	helpers.RespondWithJSON(w, 200, exercises)
}
