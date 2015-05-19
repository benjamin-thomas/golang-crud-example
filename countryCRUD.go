package main

import (
	"database/sql"
	"errors"
	"fmt"
)

type country struct {
	Id              int           `json:"id"`
	Name            string        `json:"name"`
	HasStats        bool          `json:hasStats`
	PopulationCount sql.NullInt64 `json:populationCount`
	Cities          []city        `json:cities`
}

func (c *country) create() error {
	res, err := db.Exec("INSERT INTO countries (name) VALUES ($1)", c.Name)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return errors.New(funcName() + ": Failed to create country")
	}
	return nil
}

func (c *country) read() error {
	err := db.QueryRow(`
SELECT c.name
     , (cs.id IS NOT NULL) AS has_stats
  FROM countries AS c
  LEFT
  JOIN country_stats AS cs
    ON cs.country_id = c.id
 WHERE c.id          = $1 ;
	`, &c.Id).Scan(&c.Name, &c.HasStats)
	if err == sql.ErrNoRows {
		return err
	}
	if err != nil {
		return err
	}
	return nil
}

func (c *country) update() error {
	res, err := db.Exec("UPDATE countries SET name = $2 WHERE id = $1", c.Id, c.Name)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return fmt.Errorf("%s: Failed to update country: %+v", c)
	}
	return nil
}

func (c *country) delete() error {
	res, err := db.Exec("DELETE FROM countries WHERE id = $1", c.Id)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return fmt.Errorf("%s: Failed to update country: %+v", c)
	}
	return nil
}

func (c *country) stats() error {
	err := db.QueryRow(`
		SELECT c.name
				 , cs.population_count
			FROM countries AS c
			LEFT
			JOIN country_stats AS cs
				ON cs.country_id = c.id
		 WHERE c.id          = $1;
		`, &c.Id).Scan(&c.Name, &c.PopulationCount)
	if err == sql.ErrNoRows {
		return err
	}
	if err != nil {
		return err
	}
	return nil
}

type city struct {
	Id   int    `json: id`
	Name string `json: name`
}

func (c *country) indexCities() error {
	rows, err := db.Query(`
		SELECT c.name
		     , ct.id AS city_id
		     , ct.name AS city_name
		  FROM countries AS c
		  LEFT
		  JOIN cities AS ct
		    ON ct.country_id = c.id
		 WHERE c.id          = $1;
		`, &c.Id)

	defer rows.Close()

	var ct city
	for rows.Next() {
		err = rows.Scan(&c.Name, &ct.Id, &ct.Name) // FIXME: must I override &c.Name over and over? Or use 2 queries?
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
