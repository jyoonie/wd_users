-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS wdiet;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS wdiet.users
(
    user_uuid uuid not null default gen_random_uuid()
        constraint users_primary_key
            primary key,
    hashed_password varchar(128)    not null,
    active          boolean         not null default true,
    first_name      varchar(64)     not null,
    last_name       varchar(64)     not null,
    email_address   varchar(128)    not null UNIQUE,
    created_at      timestamp       not null default now(),
    updated_at      timestamp       not null default now()
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS wdiet;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS wdiet.users;
-- +goose StatementEnd