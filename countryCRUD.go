package main

import (
	"database/sql"
)

type country struct {
	Id              int           `json:"id"`
	Name            string        `json:"name"`
	PopulationCount sql.NullInt64 `json:"population_count"`
	Cities          []city        `json:"cities"`
}

func (c *country) create() error {
	_, err := db.Exec("INSERT INTO countries (name) VALUES ($1)", c.Name)
	return err
}

func (c *country) read() error {
	return db.QueryRow(`
		SELECT c.name
		  FROM countries AS c
		  LEFT
		  JOIN country_stats AS cs
		    ON cs.country_id = c.id
		 WHERE c.id          = $1
	`, &c.Id).Scan(&c.Name)
}

func (c *country) update() error {
	_, err := db.Exec("UPDATE countries SET name = $2 WHERE id = $1", c.Id, c.Name)
	return err
}

func (c *country) delete() error {
	_, err := db.Exec("DELETE FROM countries WHERE id = $1", c.Id)
	return err
}

func (c *country) stats() error {
	return db.QueryRow(`
		SELECT c.name
				 , cs.population_count
			FROM countries AS c
			LEFT
			JOIN country_stats AS cs
				ON cs.country_id = c.id
		 WHERE c.id          = $1
		`, &c.Id).Scan(&c.Name, &c.PopulationCount)
}

func (c *country) CitiesCount() (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM cities WHERE id = $1", &c.Id).Scan(&count)
	return count, err
}

func (c *country) indexCities(per, page int) error {
	offset := per * (page - 1)

	err := db.QueryRow("SELECT name FROM countries WHERE id = $1", &c.Id).Scan(&c.Name)
	if err != nil {
		return err
	}

	rows, err := db.Query(`SELECT id, name FROM cities WHERE country_id = $1  LIMIT $2 OFFSET $3`, &c.Id, &per, &offset)
	if err != nil {
		return err
	}

	defer rows.Close()

	var ct city
	for rows.Next() {
		err = rows.Scan(&ct.Id, &ct.Name)
		if err != nil {
			return err
		}
		c.Cities = append(c.Cities, ct)
	}

	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}
