-- +goose Up
-- SQL in this section is executed when the migration is applied.

ALTER TABLE collections ALTER COLUMN ref SET DEFAULT 'master';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

ALTER TABLE collections ALTER COLUMN ref DROP DEFAULT;