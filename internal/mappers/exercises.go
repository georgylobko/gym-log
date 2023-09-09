package mappers

import "github.com/georgylobko/gym-log/internal/database"

type Exercise struct {
	ID           int32         `json:"id"`
	Name         string        `json:"name"`
	PhotoUrl     string        `json:"photo_url"`
	MuscleGroups []MuscleGroup `json:"muscle_groups"`
}

type ExerciseRow struct {
	ID           int32    `json:"id"`
	Name         string   `json:"name"`
	PhotoUrl     string   `json:"photo_url"`
	MuscleGroups []string `json:"muscle_groups"`
}

func DatabaseExerciseToExercise(dbExercise database.Exercise, dbMuscleGroups []database.MuscleGroup) Exercise {
	return Exercise{
		ID:           dbExercise.ID,
		Name:         dbExercise.Name,
		PhotoUrl:     dbExercise.PhotoUrl,
		MuscleGroups: DatabaseMuscleGroupsToMuscleGroups(dbMuscleGroups),
	}
}

func DatabaseExerciseRowToExercise(dbExercise database.GetExircisesRow) ExerciseRow {
	return ExerciseRow{
		ID:           dbExercise.ID,
		Name:         dbExercise.Name,
		PhotoUrl:     dbExercise.PhotoUrl,
		MuscleGroups: dbExercise.MuscleGroups,
	}
}

func DatabaseExercisesToExercises(dbExercises []database.GetExircisesRow) []ExerciseRow {
	exercises := []ExerciseRow{}

	for _, dbExercise := range dbExercises {
		exercises = append(exercises, DatabaseExerciseRowToExercise(dbExercise))
	}

	return exercises
}
