SELECT 'db/migrations/up/5.sql' AS filepath;

CREATE TABLE addresses (
  id SERIAL PRIMARY KEY,
  city_id INT NOT NULL REFERENCES cities(id),

  name VARCHAR(100) NOT NULL,

  created_on TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  updated_on TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX addresses_name_uidx ON addresses(name);

CREATE TRIGGER addresses_touch_trg
  BEFORE UPDATE
  ON addresses
  FOR EACH ROW
    EXECUTE PROCEDURE touch();

INSERT INTO migrations VALUES (5);
