-- +goose Up
-- SQL in this section is executed when the migration is applied.

ALTER TABLE ammo RENAME TO ammos;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

ALTER TABLE ammos RENAME TO ammo;
