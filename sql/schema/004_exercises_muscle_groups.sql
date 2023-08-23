-- +goose Up

CREATE TABLE exercises_muscle_groups (
    id SERIAL NOT NULL PRIMARY KEY,
    exercise_id SERIAL NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
    muscle_group_id SERIAL NOT NULL REFERENCES muscle_groups(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE exercises_muscle_groups;