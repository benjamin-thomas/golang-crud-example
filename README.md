# golang CRUD example

## Prerequisite

Docker

# set up the env vars below in this file  ~/.env/golang-crud-example
```bash
env | grep ^PG

PGPASSWORD=my_password
PGHOST=localhost
PGPORT=5432
PGUSER=postgres
PGDATABASE=golang_crud_example_dev
```

### Easier setup from scratch

0. ./manage/destroy_db_container (optional)
1. ./manage/create_db_container
2. ./manage/create_db
3. ./manage/seed_db
4. ./manage/psql # test the connection
5. [HTTP_SOCKET=localhost:8080] ./manage/run

---

## Miscellaneous

### Migrations

#### Reset the db
```bash
./manage/reset_db
```

#### Run one migration
```bash
./manage/psql-pipe -a < ./db/migrations/up/1.sql
```

#### Revert one migration
```bash
./manage/psql-pipe -a < ./db/migrations/down/1.sql
```

#### Run multiple migrations
```bash
cat ./db/migrations/down/{02..01}.sql | ./manage/psql-pipe -a
cat ./db/migrations/up/{01..02}.sql | ./manage/psql-pipe -a
```
or
```bash
find ./db/migrations/down/ -type f -name "*.sql" | sort -rn | xargs cat | ./manage/psql-pipe
find ./db/migrations/up/ -type f -name "*.sql" | sort -n | xargs cat | ./manage/psql-pipe
```

#### Get current migration version
```bash
./manage/psql-pipe -c "SELECT max(version) FROM migrations"
```

### 3. Dump the schema
```bash
./manage/pg_dump -s $PGDATABASE >./db/schema.sql
```

### Develop queries
cat ./queries/addresses.sql | ./manage/psql-pipe
