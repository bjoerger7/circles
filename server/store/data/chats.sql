CREATE TABLE chats
(
    id          TEXT    NOT NULL    PRIMARY KEY,
    name        TEXT    NOT NULL,
    created_at  TEXT    DEFAULT     CURRENT_TIMESTAMP
);