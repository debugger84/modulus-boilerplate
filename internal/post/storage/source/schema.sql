
create table "post"."post"
(
    id           uuid                     not null
        constraint user_pk
            primary key,
    author_id    uuid                     not null,
    title        varchar(255)             not null,
    previewImage varchar(255),
    content      jsonb,
    created_at   timestamp with time zone not null
);

