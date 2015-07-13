# Source: `https://registry.hub.docker.com/_/postgres/`

## 1. Source the env vars
```bash
set -a; source ~/.env/golang-crud-example; set +a
```

## 2. Ensure the following env vars are loaded
```bash
env | grep ^PG
# PGPASSWORD=my_password
# PGHOST=pg
# PGPORT=5432
# PGUSER=postgres
# PGDATABASE=mydb_dev
```

## 3. Launch the db container
# Run with the following command
```bash
docker run --name golang_crud_example-pg -e POSTGRES_PASSWORD=$PGPASSWORD -d postgres
```

## 4. Find the containrs running IP
 - Override the PGHOST env var (set to "pg", for linked containers)

 ```bash
 PGHOST=$(docker inspect --format '{{ .NetworkSettings.IPAddress  }}' golang_crud_example-pg)
 # create the db
 psql
 ```
 - Or type the IP manually

 ```bash
 docker inspect --format '{{ .NetworkSettings.IPAddress  }}' golang_crud_example-pg
 # create the db
 psql -h ip
 ```

## 5. Create the database
```bash
./bin/psql -d postgres -c "CREATE DATABASE $PGDATABASE;"
```
