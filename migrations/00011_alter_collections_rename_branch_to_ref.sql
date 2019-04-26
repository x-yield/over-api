-- +goose Up
-- SQL in this section is executed when the migration is applied.

ALTER TABLE collections RENAME COLUMN branch TO ref;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

ALTER TABLE collections RENAME COLUMN ref TO branch;
