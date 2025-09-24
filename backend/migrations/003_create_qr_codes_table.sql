-- +goose Up
CREATE TABLE "qr_codes" (
  "id" serial PRIMARY KEY,
  "link_id" integer NOT NULL UNIQUE,
  "color" varchar(6) NOT NULL,
  "background" varchar(6) NOT NULL,
  "smoothing" float
);

ALTER TABLE "qr_codes" ADD FOREIGN KEY ("link_id") REFERENCES "links" ("id");

-- +goose Down
DROP TABLE IF EXISTS "qr_codes";