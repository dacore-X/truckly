DROP TABLE IF EXISTS deliveries;
DROP TABLE IF EXISTS geo;

CREATE TABLE deliveries (
  id bigserial PRIMARY KEY,
  client_id bigint NOT NULL,
  courier_id bigint,
  status_id bigint NOT NULL,
  type_id bigint,
  from_longitude float8 NOT NULL,
  from_latitude float8 NOT NULL,
  to_longitude float8 NOT NULL,
  to_latitude float8 NOT NULL,
  distance float8 NOT NULL,
  price float8 NOT NULL,
  has_loader bool NOT NULL DEFAULT (FALSE),
  created_at timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE deliveries ADD FOREIGN KEY (client_id) REFERENCES users (id);

ALTER TABLE deliveries ADD FOREIGN KEY (courier_id) REFERENCES users (id);

ALTER TABLE deliveries ADD FOREIGN KEY (status_id) REFERENCES statuses (id);

ALTER TABLE deliveries ADD FOREIGN KEY (type_id) REFERENCES delivery_types (id);