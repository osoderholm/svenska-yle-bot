create table if not exists bot.subscribers
(
    id              serial                   not null primary key,
    chat_id         varchar(512)             not null unique,
    update_interval integer                  not null,
    last_article_id integer                  not null default 0,
    last_notified   timestamp with time zone null
);