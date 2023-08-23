-- name: CreateExercise :one
INSERT INTO exercises (name, photo_url)
VALUES ($1, $2)
RETURNING *;

-- name: GetExerciseById :one
SELECT * FROM exercises
WHERE id = $1;

-- name: GetExircises :many
SELECT
    e.id,
    e.name,
    e.photo_url,
    ARRAY_AGG(m.name)::TEXT[] AS muscle_groups
FROM
    exercises AS e
LEFT JOIN
    exercises_muscle_groups AS emg ON e.id = emg.exercise_id
LEFT JOIN
    muscle_groups AS m ON emg.muscle_group_id = m.id
GROUP BY
    e.id, e.name, e.photo_url;