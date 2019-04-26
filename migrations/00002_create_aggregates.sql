-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE IF NOT EXISTS aggregates (
    id                   SERIAL               PRIMARY KEY,
    job_id               INTEGER              REFERENCES jobs(id)          ON DELETE CASCADE,
    q50                  REAL,
    q75                  REAL,
    q80                  REAL,
    q85                  REAL,
    q90                  REAL,
    q95                  REAL,
    q98                  REAL,
    q99                  REAL,
    q100                 REAL,
    avg                  REAL,
    ok_count             INTEGER,
    err_count            INTEGER,
    response_code        TEXT,
    net_recv             REAL,
    net_send             REAL,
    label                TEXT
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE aggregates;
