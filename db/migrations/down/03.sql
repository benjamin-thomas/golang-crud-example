SELECT 'db/migrations/down/3.sql' AS filepath;

DROP TRIGGER cities_touch_trg ON cities;

DROP TABLE cities;

DELETE FROM migrations WHERE version = 3;
