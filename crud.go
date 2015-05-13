// Data Access Layer

package main

import "database/sql"

type countriesCrud struct {
}

type crudDefs struct {
	countries countriesCrud
}

var crud = crudDefs{
	countriesCrud{},
}

func (c countriesCrud) create(name string) (sql.Result, error) {
	stmt := mustPrepare("INSERT INTO countries (name) VALUES ($1)")
	return stmt.Exec(name)
}

func (c countriesCrud) update(id, name string) (sql.Result, error) {
	stmt := mustPrepare("UPDATE countries SET name = $2 WHERE id = $1")
	return stmt.Exec(id, name)
}
