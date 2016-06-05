# golang CRUD example


## Prerequisite

1. Postgresql
2. bash (for the dev scripts in ./bin)

## Setup

### 1. Setup the database

Follow the instructions in ./config/docker/pg/README.md to setup a docker container.

In a nutshell, the app requires a few env vars. One could store those key-values in a text file.

```bash
set -a; source ~/.env/golang-crud-example; set +a
env | grep ^PG
# PGPASSWORD=my_password
# PGHOST=localhost
# PGPORT=5432
# PGUSER=postgres
# PGDATABASE=golang_crud_example_dev

psql -d postgres -c "CREATE DATABASE $PGDATABASE;"
```

NOTE: `./bin/pg_dump`, `./bin/psql` and `./bin/psql-pipe` are wrapper scripts that launch docker containers. Those are used by the dev scripts, in `./bin/*`. If not using docker you will have to tweak those scripts for now, and make them point to your system's binaries.

### 2. Run migrations

#### Run all migrations
```bash
find ./db/migrations/up/ -type f -name "*.sql" | sort -n | xargs cat | ./bin/psql-pipe
```

### 3 Seed the db
```bash
cat ./db/seed.sql | ./bin/psql-pipe
```

### 4. Start the app
```bash
CRUD_USER=my_user CRUD_PW=my_password DEV=1 go run *.go
```
or (reload with control-c)
```bash
while true; do CRUD_USER=my_user CRUD_PW=my_password DEV=1 go run *.go; sleep 1; done

Then go to http://localhost:8080
```

---

## Miscellaneous

### Migrations

#### Reset the db
```bash
./bin/reset_db
```

#### Run one migration
```bash
./bin/psql-pipe -a < ./db/migrations/up/1.sql
```

#### Revert one migration
```bash
./bin/psql-pipe -a < ./db/migrations/down/1.sql
```

#### Run multiple migrations
```bash
cat ./db/migrations/down/{02..01}.sql | ./bin/psql-pipe -a
cat ./db/migrations/up/{01..02}.sql | ./bin/psql-pipe -a
```
or
```bash
find ./db/migrations/down/ -type f -name "*.sql" | sort -rn | xargs cat | ./bin/psql-pipe
find ./db/migrations/up/ -type f -name "*.sql" | sort -n | xargs cat | ./bin/psql-pipe
```

#### Get current migration version
```bash
./bin/psql-pipe -c "SELECT max(version) FROM migrations"
```

### 3. Dump the schema
```bash
./bin/pg_dump -s $PGDATABASE >./db/schema.sql
```

### Develop queries
cat ./queries/addresses.sql | ./bin/psql-pipe


### Easier setup from scratch

0. ./manage/destroy_db_container (optional)
1. ./manage/create_db_container
2. ./manage/create_db
3. ./manage/seed_db
4. ./manage/psql # test the connection
5. [HTTP_SOCKET=localhost:8080] ./manage/run
