package mappers

import "github.com/georgylobko/gym-log/internal/database"

type Set struct {
	ID         int32 `json:"id"`
	WorkoutID  int32 `json:"workout_id"`
	ExerciseID int32 `json:"exercise_id"`
	Reps       int32 `json:"reps"`
	Weight     int32 `json:"weight"`
}

func DatabaseSetToSet(dbSet database.Set) Set {
	return Set{
		ID:         dbSet.ID,
		WorkoutID:  dbSet.WorkoutID,
		ExerciseID: dbSet.ExerciseID,
		Reps:       dbSet.Reps,
		Weight:     dbSet.Weight,
	}
}

func DatabaseSetsToSets(dbSets []database.Set) []Set {
	sets := []Set{}

	for _, dbSet := range dbSets {
		sets = append(sets, DatabaseSetToSet(dbSet))
	}

	return sets
}
