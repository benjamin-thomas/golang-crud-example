SELECT 'db/migrations/down/02_drop_table_countries.sql' AS filepath;

DROP TRIGGER countries_touch_trg ON countries;

DROP TABLE countries;

DELETE FROM migrations WHERE version = 2;
