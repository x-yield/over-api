-- +goose Up
-- SQL in this section is executed when the migration is applied.

ALTER TABLE IF EXISTS jobs ADD COLUMN autostop_time DOUBLE PRECISION;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

ALTER TABLE IF EXISTS jobs DROP COLUMN autostop_time;
