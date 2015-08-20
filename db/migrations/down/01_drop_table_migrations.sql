SELECT 'db/migrations/down/01_drop_table_migrations.sql' AS filepath;

DROP FUNCTION touch();

DROP TABLE migrations;
