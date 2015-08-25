package main

import (
	"database/sql"
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
	ZipCode *string `json:"zip_code"`
	Country string  `json:"country"`
}

type addresses []address

func (as *addresses) index(sPer, sPage, q, condOp, matchOp string, cols []string) error {
	per, offset := paginateParams(sPer, sPage)

	if condOp != "AND" && condOp != "OR" {
		return &syntaxErr{fmt.Sprintf("index: invalid conditional operator '%s'", condOp, "Use: 'AND' or 'OR'")}
	}

	if matchOp != "=" && matchOp != "LIKE" && matchOp != "ILIKE" {
		return &syntaxErr{fmt.Sprintf("index: invalid match operator '%s', Use: '=' or 'LIKE' or 'ILIKE'", matchOp)}
	}

	if len(cols) == 0 || cols[0] == "" {
		cols = []string{"name"}
	}

	allowedCols := []string{"id", "name", "line1", "line2", "line3", "city", "zip_code", "country"}
	for _, c := range cols {
		for i, ac := range allowedCols {
			if c == ac {
				break
			}
			if i == len(allowedCols)-1 {
				return &syntaxErr{fmt.Sprintf("index: Bad column name '%s'. Allowed column names: %#v", c, allowedCols)}
			}
		}
	}

	// Using a CTE to avoid name collisions on filtering
	qry := `
	WITH qry AS (
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
	)

	SELECT * FROM qry
	`

	if q != "" {
		for i, c := range cols {
			if i == 0 {
				qry += "WHERE "
			} else {
				qry += "\n" + condOp + " "
			}
			qry += c + " " + matchOp + " $3"
		}
	}

	qry += "\nORDER BY id LIMIT $1 OFFSET $2"

	fmt.Println("\033[1;33m", strings.Replace(qry, "\t", "  ", -1), "\n\033[1;m")

	var rows *sql.Rows
	var err error
	if q == "" {
		rows, err = db.Query(qry, per, offset)
	} else {
		fmt.Println("$1 =", per)
		fmt.Println("$2 =", offset)
		fmt.Println("$3 =", q)
		rows, err = db.Query(qry, per, offset, q)
	}
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
