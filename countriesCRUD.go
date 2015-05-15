package main

import "database/sql"

type countries struct {
}

func (c *countries) Count() (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM countries").Scan(&count)
	return count, err
}

func (c *countries) Index(per, page int) (*sql.Rows, error) {
	offset := per * (page - 1)
	return db.Query(`SELECT c.id
															, c.name
												 FROM countries AS c
												 ORDER BY c.name ASC
												 LIMIT $1 OFFSET $2`, &per, &offset)

}
