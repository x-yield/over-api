-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE IF NOT EXISTS collections (
    id                   SERIAL,
    service              TEXT,
    "name"               TEXT,
    env                  TEXT,
    author               VARCHAR(256),
    branch               TEXT,
    created_at           TIMESTAMP           DEFAULT     CURRENT_TIMESTAMP,
    CONSTRAINT collection_entry UNIQUE (env, service, branch, "name")
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE collections;
