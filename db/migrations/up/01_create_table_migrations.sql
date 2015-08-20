SELECT 'db/migrations/up/01_create_table_migrations.sql' AS filepath;

CREATE TABLE migrations (
  version INT NOT NULL
);

CREATE UNIQUE INDEX migrations_version_uidx ON migrations(version);

-- Inspired from: http://www.depesz.com/2012/11/14/how-i-learned-to-stop-worrying-and-love-the-triggers/
CREATE FUNCTION touch() RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'UPDATE' THEN
        NEW.updated_on := now();
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

INSERT INTO migrations VALUES (1);
