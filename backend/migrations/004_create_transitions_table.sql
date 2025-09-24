-- +goose Up
CREATE TABLE "transitions" (
  "id" serial PRIMARY KEY,
  "link_id" integer NOT NULL,
  "country" varchar,
  "city" varchar,
  "referer" varchar,
  "user_agent" varchar,
  "browser" varchar,
  "os" varchar,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "transitions" ADD FOREIGN KEY ("link_id") REFERENCES "links" ("id");
CREATE INDEX ON "transitions" ("link_id");

-- +goose Down
DROP TABLE IF EXISTS "transitions";