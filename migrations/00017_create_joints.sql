-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE IF NOT EXISTS joints (
    id                   SERIAL,
    jobs                 INTEGER[],
    name                 TEXT,
    created_at           TIMESTAMP           DEFAULT     CURRENT_TIMESTAMP
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE joints;
