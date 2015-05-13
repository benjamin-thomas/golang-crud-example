// Data Access Layer

package main

import (
	"errors"
)

type countriesCrud struct {
}

type crudDefs struct {
	countries countriesCrud
}

var crud = crudDefs{
	countriesCrud{},
}

func (c countriesCrud) create(name string) error {
	stmt := mustPrepare("INSERT INTO countries (name) VALUES ($1)")
	res, err := stmt.Exec(name)
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

func (c countriesCrud) update(id, name string) error {
	stmt := mustPrepare("UPDATE countries SET name = $2 WHERE id = $1")
	res, err := stmt.Exec(id, name)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return errors.New(funcName() + ": Failed to update country id: " + id)
	}
	return nil
}

func (c countriesCrud) delete(id string) error {
	stmt := mustPrepare("DELETE FROM countries WHERE id = $1")

	res, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return errors.New(funcName() + ": Failed to delete country id: " + id)
	}
	return nil
}
