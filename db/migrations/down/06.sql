SELECT 'db/migrations/down/06.sql' AS filepath;

DROP TRIGGER country_stats_touch_trg ON country_stats;

DROP TABLE country_stats;

DELETE FROM migrations WHERE version = 6;
