-- +goose Up
-- +goose StatementBegin
create table users (
    id serial primary key,
    email text not null UNIQUE,
    hashedpassword bytea not null
);

create index users_email_idx on users(email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index users_email_idx;
drop table users;
-- +goose StatementEnd
