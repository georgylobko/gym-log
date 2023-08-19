package main

import "github.com/georgylobko/gym-log/internal/database"

type MuscleGroup struct {
	ID       int32  `json:"id"`
	Name     string `json:"name"`
	PhotoUrl string `json:"photo_url"`
}

func databaseMuscleGroupToMuscleGroup(dbMuscleGroup database.MuscleGroup) MuscleGroup {
	return MuscleGroup{
		ID:       dbMuscleGroup.ID,
		Name:     dbMuscleGroup.Name,
		PhotoUrl: dbMuscleGroup.PhotoUrl,
	}
}
