-- +goose Up

CREATE TABLE muscle_groups (
    id SERIAL NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    photo_url TEXT NOT NULL
);

-- +goose Down
DROP TABLE muscle_groups;