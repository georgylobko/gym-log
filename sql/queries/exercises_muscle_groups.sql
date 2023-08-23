-- name: CreateExerciseMuscleGroup :one
INSERT INTO exercises_muscle_groups (exercise_id, muscle_group_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetMuscleGroupsByExercise :many
SELECT muscle_groups.id, muscle_groups.name, muscle_groups.photo_url 
FROM exercises_muscle_groups
JOIN muscle_groups ON muscle_groups.id = muscle_group_id
WHERE exercise_id = $1;
