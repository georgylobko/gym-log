-- name: CreateExercise :one
INSERT INTO exercises (name, photo_url)
VALUES ($1, $2)
RETURNING *;

-- name: GetExerciseById :one
SELECT * FROM exercises
WHERE id = $1;