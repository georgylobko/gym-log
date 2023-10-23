package mappers

import (
	"reflect"
	"testing"

	"github.com/georgylobko/gym-log/internal/database"
)

func TestDatabaseExerciseToExercise(t *testing.T) {
	dbExercise := database.GetExircisesRow{
		ID:           1,
		Name:         "Bench press",
		PhotoUrl:     "http://im.com/benchpress",
		MuscleGroups: []string{"Chest"},
	}

	exercise := ExerciseRow{
		ID:           dbExercise.ID,
		Name:         dbExercise.Name,
		PhotoUrl:     dbExercise.PhotoUrl,
		MuscleGroups: dbExercise.MuscleGroups,
	}

	if !reflect.DeepEqual(DatabaseExerciseRowToExercise(dbExercise), exercise) {
		t.Fatalf("db exercise does not equal exercise after mapping")
	}
}
