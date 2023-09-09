package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/georgylobko/gym-log/internal/database"
	"github.com/georgylobko/gym-log/internal/helpers"
	"github.com/georgylobko/gym-log/internal/mappers"
)

func (apiCfg *ApiConfig) HandlerCreateSet(w http.ResponseWriter, r *http.Request, user mappers.User) {
	type parameters struct {
		WorkoutID  int32 `json:"workout_id"`
		ExerciseID int32 `json:"exercise_id"`
		Reps       int32 `json:"reps"`
		Weight     int32 `json:"weight"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}

	set, err := apiCfg.DB.CreateSet(r.Context(), database.CreateSetParams{
		WorkoutID:  params.WorkoutID,
		ExerciseID: params.ExerciseID,
		Reps:       params.Reps,
		Weight:     params.Weight,
	})
	if err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprintf("Could not create the db entity: %s", err))
		return
	}

	helpers.RespondWithJSON(w, 200, set)
}

func (apiCfg *ApiConfig) HandlerGetSets(w http.ResponseWriter, r *http.Request, user mappers.User) {
	query := r.URL.Query()
	workoutID := query.Get("workout_id")

	if workoutID != "" {
		workoutID, _ := strconv.Atoi(workoutID)
		sets, err := apiCfg.DB.GetSetsByWorkout(r.Context(), int32(workoutID))
		if err != nil {
			helpers.RespondWithError(w, 400, fmt.Sprintf("Could not get the db entity: %s", err))
			return
		}

		helpers.RespondWithJSON(w, 200, sets)
	} else {
		sets, err := apiCfg.DB.GetSets(r.Context())
		if err != nil {
			helpers.RespondWithError(w, 400, fmt.Sprintf("Could not get the db entity: %s", err))
			return
		}

		helpers.RespondWithJSON(w, 200, sets)
	}
}
