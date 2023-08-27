-- name: CreateSet :one
INSERT INTO sets (workout_id, exercise_id, reps, weight)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetSets :many
SELECT * FROM sets;

-- name: GetSetsByWorkout :many
SELECT * FROM sets
WHERE workout_id = $1;