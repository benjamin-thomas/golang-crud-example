SELECT 'db/migrations/up/02_create_table_countries.sql' AS filepath;

CREATE TABLE countries (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,

  created_on TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  updated_on TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX countries_name_uidx ON countries(name);

CREATE TRIGGER countries_touch_trg
  BEFORE UPDATE
  ON countries
  FOR EACH ROW
    EXECUTE PROCEDURE touch();

INSERT INTO migrations VALUES (2);
