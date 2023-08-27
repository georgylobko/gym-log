-- +goose Up

CREATE TABLE sets (
    id SERIAL NOT NULL PRIMARY KEY,
    workout_id SERIAL NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,
    exercise_id SERIAL NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
    reps INTEGER NOT NULL,
    weight INTEGER NOT NULL
);

-- +goose Down
DROP TABLE sets;