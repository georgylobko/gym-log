-- +goose Up

CREATE TABLE users (
    id SERIAL NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    gender VARCHAR(5),
    role VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password TEXT NOT NULL
);

-- +goose Down
DROP TABLE users;