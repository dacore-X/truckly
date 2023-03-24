CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "surname" varchar NOT NULL,
  "name" varchar NOT NULL,
  "patronymic" varchar,
  "email" varchar UNIQUE NOT NULL,
  "phone_number" varchar NOT NULL,
  "hash_password" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "meta" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint,
  "is_admin" bool NOT NULL DEFAULT (FALSE),
  "is_courier" bool NOT NULL DEFAULT (FALSE),
  "is_banned" bool NOT NULL DEFAULT (FALSE),
  "rating" numeric(3, 2) NOT NULL DEFAULT (5.00)
);

CREATE TABLE "statuses" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL
);

CREATE TABLE "trucks" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "description" varchar NOT NULL
);

CREATE TABLE "deliveries" (
  "id" bigserial PRIMARY KEY,
  "client_id" bigint NOT NULL,
  "courier_id" bigint NOT NULL,
  "status_id" bigint NOT NULL,
  "truck_id" bigint,
  "from_longitude" float8 NOT NULL,
  "from_latitude" float8 NOT NULL,
  "to_longitude" float8 NOT NULL,
  "to_latitude" float8 NOT NULL,
  "distance" float8 NOT NULL,
  "price" float8 NOT NULL,
  "has_loader" bool NOT NULL DEFAULT (FALSE),
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "meta" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "deliveries" ADD FOREIGN KEY ("client_id") REFERENCES "users" ("id");

ALTER TABLE "deliveries" ADD FOREIGN KEY ("courier_id") REFERENCES "users" ("id");

ALTER TABLE "deliveries" ADD FOREIGN KEY ("status_id") REFERENCES "statuses" ("id");

ALTER TABLE "deliveries" ADD FOREIGN KEY ("truck_id") REFERENCES "trucks" ("id");
