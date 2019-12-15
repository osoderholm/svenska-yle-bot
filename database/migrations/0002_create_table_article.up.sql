create table if not exists bot.articles
(
    id        serial                   not null primary key,
    title     text                     not null,
    link      varchar(512)             not null unique,
    image     varchar(512)             not null,
    published timestamp with time zone not null
);