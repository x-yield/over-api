-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TYPE env_enum AS ENUM ('prod', 'stg', 'dev');
ALTER TABLE collections ALTER COLUMN env TYPE env_enum USING env::env_enum;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

ALTER TABLE collections ALTER COLUMN env TYPE TEXT;
DROP TYPE env_enum;
