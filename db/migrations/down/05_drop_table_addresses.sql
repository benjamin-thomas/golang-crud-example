SELECT 'db/migrations/down/05_drop_table_addresses.sql' AS filepath;

DROP TRIGGER addresses_touch_trg ON addresses;

DROP TABLE addresses;

DELETE FROM migrations WHERE version = 5;
