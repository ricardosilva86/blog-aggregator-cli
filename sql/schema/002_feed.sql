-- +goose Up
CREATE TABLE feeds (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL,
    url text not null unique,
    user_id uuid,
    created_at timestamp NOT NULL default now(),
    updated_at timestamp NOT NULL DEFAULT now(),
    constraint fk_users_feeds
        foreign key (user_id)
        references users(id)
        on delete cascade
);
-- +goose Down
drop table feeds;