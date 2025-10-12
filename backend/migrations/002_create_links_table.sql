-- +goose Up
CREATE TABLE "links" (
  "id" serial PRIMARY KEY,
  "original_url" varchar NOT NULL,
  "hash" varchar NOT NULL UNIQUE,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "user_id" integer NOT NULL
);

ALTER TABLE "links" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
CREATE INDEX ON "links" ("user_id");

-- +goose Down
DROP TABLE IF EXISTS "links";