-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE IF NOT EXISTS tank_sessions (
    id                   SERIAL,
    tank                 VARCHAR(256),
    conf                 TEXT,
    "name"               VARCHAR(256),
    failures             TEXT[],
    stage                VARCHAR(256),
    "status"             VARCHAR(256),
    external_id          VARCHAR(256),
    overload_id          INT,
    external_joint       VARCHAR(256),
    overload_joint       INT,
    author               VARCHAR(256),
    created_at           TIMESTAMP           DEFAULT     CURRENT_TIMESTAMP,
    CONSTRAINT tank_session_entry UNIQUE (tank, "name")
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE tank_sessions;
