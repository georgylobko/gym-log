-- +goose Up

CREATE TABLE workouts (
    id SERIAL NOT NULL PRIMARY KEY,
    user_id SERIAL NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL 
);

-- +goose Down
DROP TABLE workouts;