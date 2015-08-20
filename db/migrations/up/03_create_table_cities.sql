SELECT 'db/migrations/up/03_create_table_cities.sql' AS filepath;

CREATE TABLE cities (
  id SERIAL PRIMARY KEY,
  country_id INT NOT NULL REFERENCES countries(id) ON DELETE CASCADE,

  name VARCHAR(100) NOT NULL CHECK (trim(name) != ''),

  created_on TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  updated_on TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX cities_name_uidx ON cities(name);

CREATE TRIGGER cities_touch_trg
  BEFORE UPDATE
  ON cities
  FOR EACH ROW
    EXECUTE PROCEDURE touch();

INSERT INTO migrations VALUES (3);
