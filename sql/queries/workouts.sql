-- name: CreateWorkout :one
INSERT INTO workouts (user_id, created_at, updated_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetWorkoutsByUserId :many
SELECT * FROM workouts
WHERE user_id = $1;

-- name: GetWorkoutById :one
SELECT * FROM workouts
WHERE id = $1;

-- name: UpdateWorkout :exec
UPDATE workouts 
SET updated_at = $2
WHERE id = $1;