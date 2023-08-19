-- name: CreateMuscleGroup :one
INSERT INTO muscle_groups (name, photo_url)
VALUES ($1, $2)
RETURNING *;