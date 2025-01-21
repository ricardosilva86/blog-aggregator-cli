-- +goose Up
CREATE TABLE feed_follows (
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null,
    feed_id uuid not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    UNIQUE (feed_id, user_id), -- this makes sure that user_id and feed_id combinations are always unique
    constraint fk_users_feed_follows
        foreign key (user_id)
        references users(id)
        on delete cascade,
    constraint fk_feed_feed_follows
        foreign key (feed_id)
        references feeds(id)
        on delete cascade
);

-- +goose Down
drop table feed_follows;