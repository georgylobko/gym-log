package mappers

import (
	"time"

	"github.com/georgylobko/gym-log/internal/database"
)

type Workout struct {
	ID        int32     `json:"id"`
	UserID    int32     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func DatabaseWorkoutToWorkout(dbWorkout database.Workout) Workout {
	return Workout{
		ID:        dbWorkout.ID,
		UserID:    dbWorkout.UserID,
		CreatedAt: dbWorkout.CreatedAt,
		UpdatedAt: dbWorkout.UpdatedAt,
	}
}

func DatabaseWorkoutsToWorkouts(dbWorkouts []database.Workout) []Workout {
	workouts := []Workout{}

	for _, dbWorkout := range dbWorkouts {
		workouts = append(workouts, DatabaseWorkoutToWorkout(dbWorkout))
	}

	return workouts
}
