-- +goose Up
-- SQL in this section is executed when the migration is applied.

ALTER TABLE collections ALTER COLUMN name SET NOT NULL, ALTER COLUMN project SET NOT NULL, ALTER COLUMN env SET NOT NULL, ALTER COLUMN ref SET NOT NULL;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

ALTER TABLE collections ALTER COLUMN name DROP NOT NULL, ALTER COLUMN project DROP NOT NULL, ALTER COLUMN env DROP NOT NULL, ALTER COLUMN ref DROP NOT NULL;