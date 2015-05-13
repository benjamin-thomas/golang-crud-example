package main

import "database/sql"

type countryCrud struct {
}

type crudDefs struct {
	country countryCrud
}

var crud = crudDefs{
	countryCrud{},
}

func (c countryCrud) updateName(id, name string) (sql.Result, error) {
	stmt := mustPrepare("UPDATE countries SET name = $2 WHERE id = $1")
	return stmt.Exec(id, name)
}
