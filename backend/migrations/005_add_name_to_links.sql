-- +goose Up
ALTER TABLE "links" ADD COLUMN "name" varchar NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE "links" DROP COLUMN IF EXISTS "name";


