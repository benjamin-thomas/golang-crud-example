SELECT 'db/migrations/up/06.sql' AS filepath;

CREATE TABLE country_stats (
  id SERIAL PRIMARY KEY,
  country_id INT NOT NULL REFERENCES countries(id),

  population_count INT NOT NULL,

  created_on TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  updated_on TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX country_stats_country_id_uidx ON country_stats(country_id);

CREATE TRIGGER country_stats_touch_trg
  BEFORE UPDATE
  ON country_stats
  FOR EACH ROW
    EXECUTE PROCEDURE touch();

INSERT INTO migrations VALUES(6);
