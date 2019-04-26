-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE IF NOT EXISTS ammo (
    id                   SERIAL,
    url                  VARCHAR(256),
    bucket               VARCHAR(63),
    "key"                VARCHAR(1024),
    last_used            DOUBLE PRECISION,
    "type"               VARCHAR(256),
    author               VARCHAR(256)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE ammo;
