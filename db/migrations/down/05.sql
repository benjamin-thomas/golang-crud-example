SELECT 'db/migrations/down/5.sql' AS filepath;

DROP TRIGGER addresses_touch_trg ON addresses;

DROP TABLE addresses;

DELETE FROM migrations WHERE version = 5;
