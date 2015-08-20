SELECT 'db/migrations/down/03_drop_table_cities.sql' AS filepath;

DROP TRIGGER cities_touch_trg ON cities;

DROP TABLE cities;

DELETE FROM migrations WHERE version = 3;
