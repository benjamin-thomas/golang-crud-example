# golang CRUD example

## Setup

1. Setup the db
2. Run the migrations
3. Dump the schema
4. Seed the db (dev)

## 1. Setup the database

Follow the instructions in ./config/docker/pg/README.md

## 2. Run migrations

### Run one migration
```bash
./bin/psql-pipe -a < ./db/migrations/up/1.sql
```

### Revert one migration
```bash
./bin/psql-pipe -a < ./db/migrations/down/1.sql
```

### Run multiple migrations
```bash
cat ./db/migrations/up/{01..02}.sql | ./bin/psql-pipe -a
cat ./db/migrations/down/{02..01}.sql | ./bin/psql-pipe -a
```

### Get current migration version
```bash
./bin/psql-pipe -c "SELECT max(version) FROM migrations"
```

## 3. Dump the schema
```bash
./bin/pg_dump -s $PGDATABASE >./db/schema.sql
```

## 4 Seed the db
```bash
./bin/psql-pipe -c ./db/seed.sql
```
