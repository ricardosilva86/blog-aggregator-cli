-- +goose Up
CREATE TABLE users
(
    id         uuid PRIMARY KEY unique default gen_random_uuid(),
    name       text      NOT NULL,
    created_at timestamp NOT NULL default now(),
    updated_at timestamp NOT NULL default now()
);

-- +goose Down
DROP TABLE users;