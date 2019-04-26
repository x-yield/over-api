-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE IF NOT EXISTS jobs (
    id                   SERIAL               PRIMARY KEY,
    test_start           DOUBLE PRECISION,
    test_stop            DOUBLE PRECISION,
    config               TEXT,
    author               VARCHAR(256),
    regression_id        TEXT,
    description          TEXT,
    tank                 VARCHAR(256),
    target               VARCHAR(256),
    environment_details  TEXT,
    status               VARCHAR(256)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE jobs;
