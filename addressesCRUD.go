package main

import (
	"errors"
	"fmt"
	"strings"
)

type address struct {
	Id      int     `json:"id"`
	Name    string  `json:"name"`
	Line1   *string `json:"line1"`
	Line2   *string `json:"line2"`
	Line3   *string `json:"line3"`
	City    string  `json:"city"`
	ZipCode string  `json:"zip_code"`
	Country string  `json:"country"`
}

type addresses []address

func (as *addresses) index(sPer, sPage, q, op string, cols []string) error {
	per, offset := paginateParams(sPer, sPage)

	if op == "" {
		op = "AND"
	}

	if op != "AND" && op != "OR" {
		return errors.New(fmt.Sprintf("index: invalid operator '%s'", op))
	}

	qry := `
		SELECT a.id
				 , a.name
				 , a.line1
				 , a.line2
				 , a.line3
				 , c.name AS city
				 , c.zip_code
				 , cn.name AS country
			FROM addresses AS a
		 INNER
			JOIN cities AS c
				ON c.id = a.city_id
		 INNER
			JOIN countries AS cn
				ON cn.id = c.country_id
	`

	if q != "" {
		if len(cols) == 0 {
			cols = []string{"name"}
		}

		for i, c := range cols {
			if i == 0 {
				qry += "WHERE "
			} else {
				qry += "\n" + op + " "
			}
			switch c {
			case "name":
				c = "a.name"
			case "city":
				c = "c.name"
			}
			qry += fmt.Sprint(c, " ILIKE '%"+q+"%' ")
		}
	}

	qry += "\nORDER BY id LIMIT $1 OFFSET $2"

	fmt.Println("\033[1;33m", strings.Replace(qry, "\t", "  ", -1), "\n\033[1;m")
	rows, err := db.Query(qry, &per, &offset)
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		a := address{}
		err = rows.Scan(&a.Id, &a.Name, &a.Line1, &a.Line2, &a.Line3, &a.City, &a.ZipCode, &a.Country)
		if err != nil {
			return err
		}
		*as = append(*as, a)
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	return nil

}
