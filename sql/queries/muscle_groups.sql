-- name: CreateMuscleGroup :one
INSERT INTO muscle_groups (name, photo_url)
VALUES ($1, $2)
RETURNING *;

-- name: GetMuscleGroups :many
SELECT * FROM muscle_groups;