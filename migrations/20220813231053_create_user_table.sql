-- +goose Up
create table users
(
    id         uuid
        constraint pk_id
            primary key,
    first_name text      not null,
    last_name  text      not null,
    birthday   timestamp not null,
    created_at timestamp not null,
    updated_at timestamp
);

-- +goose Down
drop table users;
