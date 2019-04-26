-- +goose Up
-- SQL in this section is executed when the migration is applied.

ALTER TABLE collections ADD COLUMN "type" TEXT;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

ALTER TABLE collections DROP COLUMN "type";
