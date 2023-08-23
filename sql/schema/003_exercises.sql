-- +goose Up

CREATE TABLE exercises (
    id SERIAL NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    photo_url TEXT NOT NULL
);

-- +goose Down
DROP TABLE exercises;