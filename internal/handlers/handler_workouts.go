package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/georgylobko/gym-log/internal/database"
	"github.com/georgylobko/gym-log/internal/helpers"
)

func (apiCfg *ApiConfig) HandlerCreateWorkout(w http.ResponseWriter, r *http.Request, userID string) {
	parsedUserId, _ := strconv.ParseInt(userID, 10, 32)
	workout, err := apiCfg.DB.CreateWorkout(r.Context(), database.CreateWorkoutParams{
		UserID:    int32(parsedUserId),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprintf("Could not create the db entity: %s", err))
		return
	}

	helpers.RespondWithJSON(w, 200, workout)
}

func (apiCfg *ApiConfig) HandlerUpdateWorkout(w http.ResponseWriter, r *http.Request, userID string) {
	type parameters struct {
		ID int32 `json:"id"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}

	_, err = apiCfg.DB.GetWorkoutById(r.Context(), params.ID)
	if err != nil {
		helpers.RespondWithError(w, 404, fmt.Sprintf("Entity does not exist: %s", err))
		return
	}

	err = apiCfg.DB.UpdateWorkout(r.Context(), database.UpdateWorkoutParams{
		ID:        params.ID,
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprintf("Could not update the db entity: %s", err))
		return
	}

	helpers.RespondWithJSON(w, 200, struct{}{})
}

func (apiCfg *ApiConfig) HandlerGetWorkouts(w http.ResponseWriter, r *http.Request, userID string) {
	parsedUserId, _ := strconv.ParseInt(userID, 10, 32)
	workouts, err := apiCfg.DB.GetWorkoutsByUserId(r.Context(), int32(parsedUserId))
	if err != nil {
		helpers.RespondWithError(w, 400, fmt.Sprintf("Could not get the db entity: %s", err))
		return
	}

	helpers.RespondWithJSON(w, 200, workouts)
}
