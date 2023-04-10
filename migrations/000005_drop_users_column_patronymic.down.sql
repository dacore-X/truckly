DROP TABLE IF EXISTS meta; 
DROP TABLE IF EXISTS deliveries;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
  id bigserial PRIMARY KEY,
  surname varchar NOT NULL,
  name varchar NOT NULL,
  patronymic varchar,
  email varchar UNIQUE NOT NULL,
  phone_number varchar NOT NULL,
  hash_password varchar NOT NULL,
  created_at timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE meta (
  id bigserial PRIMARY KEY,
  user_id bigint,
  is_admin bool NOT NULL DEFAULT (FALSE),
  is_courier bool NOT NULL DEFAULT (FALSE),
  is_banned bool NOT NULL DEFAULT (FALSE),
  rating numeric(3, 2) NOT NULL DEFAULT (5.00)
);

CREATE TABLE deliveries (
  id bigserial PRIMARY KEY,
  client_id bigint NOT NULL,
  courier_id bigint,
  status_id bigint NOT NULL,
  type_id bigint,
  geo_id bigint NOT NULL,
  price float8 NOT NULL,
  has_loader bool NOT NULL DEFAULT (FALSE),
  created_at timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE meta ADD FOREIGN KEY (user_id) REFERENCES users (id);

ALTER TABLE deliveries ADD FOREIGN KEY (client_id) REFERENCES users (id);

ALTER TABLE deliveries ADD FOREIGN KEY (courier_id) REFERENCES users (id);

ALTER TABLE deliveries ADD FOREIGN KEY (status_id) REFERENCES statuses (id);

ALTER TABLE deliveries ADD FOREIGN KEY (type_id) REFERENCES delivery_types (id);

ALTER TABLE deliveries ADD FOREIGN KEY (geo_id) REFERENCES geo (id);