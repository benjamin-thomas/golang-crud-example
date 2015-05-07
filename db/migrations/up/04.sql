SELECT 'db/migrations/up/4.sql' AS filepath;

CREATE TABLE zip_codes (
  id SERIAL PRIMARY KEY,
  city_id INT NOT NULL REFERENCES cities(id),

  code VARCHAR(50),

  created_on TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  updated_on TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX zip_codes_city_id_code_uidx ON zip_codes(city_id, code);

CREATE TRIGGER zip_codes_touch_trg
  BEFORE UPDATE
  ON zip_codes
  FOR EACH ROW
    EXECUTE PROCEDURE touch();

INSERT INTO migrations VALUES (4);
