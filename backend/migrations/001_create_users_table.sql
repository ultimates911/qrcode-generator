-- +goose Up
CREATE TABLE "users" (
  "id" serial PRIMARY KEY,
  "name" varchar(50) NOT NULL,
  "email" varchar NOT NULL UNIQUE,
  "hashed_password" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

-- +goose Down
DROP TABLE IF EXISTS "users";