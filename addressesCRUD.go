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

func (as *addresses) index(sPer, sPage, q, op string) error {
	per, offset := paginateParams(sPer, sPage)

	if op == "" {
		op = "OR"
	}
	if op != "AND" && op != "OR" {
		return &syntaxErr{fmt.Sprintf("index: invalid conditional operator '%s', Use: 'AND' or 'OR'", op)}
	}

	qryDefs := strings.Split(q, ",")
	fmt.Println("q =", q)
	fmt.Printf("qryDefs = %#v\n", qryDefs)

	type col struct {
		colName, searchString string
	}
	var cols []col

	allowedCols := []string{"id", "name", "line1", "line2", "line3", "city", "zip_code", "country"}
	for _, d := range qryDefs {
		if d == "" {
			continue
		}
		fmt.Println("d =", d)
		dd := strings.SplitN(d, ":", 2)
		if len(dd) != 2 {
			return &syntaxErr{fmt.Sprintf("index: Bad query string: '%s', use format: 'col_name:searched_value'", q)}
		}
		colName, term := dd[0], dd[1]
		fmt.Println("colName =", colName)
		fmt.Println("term =", term)
		for i, ac := range allowedCols {
			if colName == ac {
				break
			}
			if i == len(allowedCols)-1 {
				return &syntaxErr{fmt.Sprintf("index: Bad column name '%s'. Allowed column names: %#v", colName, allowedCols)}
			}
		}
		cols = append(cols, col{colName, term})
		fmt.Println("cols =", cols)
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

	for i, c := range cols {
		if i == 0 {
			qry += "WHERE "
		} else {
			qry += "\n" + op + " "
		}
		switch {
		case strings.HasPrefix(c.searchString, "="): // force equals if needed
			cols[i].searchString = strings.Replace(c.searchString, "=", "", 1)
			qry += c.colName + " " + "=" + " " + fmt.Sprintf("$%d", i+3)
		case strings.HasPrefix(c.searchString, "!="): // force equals if needed
			cols[i].searchString = strings.Replace(c.searchString, "!=", "", 1)
			qry += c.colName + " " + "!=" + " " + fmt.Sprintf("$%d", i+3)
		case strings.HasPrefix(c.searchString, ">="):
			cols[i].searchString = strings.Replace(c.searchString, ">=", "", 1)
			qry += c.colName + " " + ">=" + " " + fmt.Sprintf("$%d", i+3)
		case strings.HasPrefix(c.searchString, ">"):
			cols[i].searchString = strings.Replace(c.searchString, ">", "", 1)
			qry += c.colName + " " + ">" + " " + fmt.Sprintf("$%d", i+3)
		case strings.HasPrefix(c.searchString, "<="):
			cols[i].searchString = strings.Replace(c.searchString, "<=", "", 1)
			qry += c.colName + " " + "<=" + " " + fmt.Sprintf("$%d", i+3)
		case strings.HasPrefix(c.searchString, "<"):
			cols[i].searchString = strings.Replace(c.searchString, "<", "", 1)
			qry += c.colName + " " + "<" + " " + fmt.Sprintf("$%d", i+3)
		case strings.Contains(c.searchString, "%"):
			qry += c.colName + " " + "LIKE" + " " + fmt.Sprintf("$%d", i+3)
		case strings.Contains(c.searchString, "*"):
			cols[i].searchString = strings.Replace(c.searchString, "*", "%", -1)
			qry += c.colName + " " + "ILIKE" + " " + fmt.Sprintf("$%d", i+3)
		default:
			qry += c.colName + " " + "=" + " " + fmt.Sprintf("$%d", i+3)
		}
	}

	qry += "\nORDER BY id LIMIT $1 OFFSET $2"

	fmt.Println("\033[1;33m", strings.Replace(qry, "\t", "  ", -1), "\n\033[1;m")

	var qryParams []interface{}
	qryParams = append(qryParams, per, offset)
	fmt.Println("$1 =", per)
	fmt.Println("$2 =", offset)
	for i, c := range cols {
		fmt.Printf("$%d = %v\n", i+3, c.searchString)
		qryParams = append(qryParams, c.searchString)
	}

	fmt.Println("qryParams =", qryParams)
	var rows *sql.Rows
	var err error
	rows, err = db.Query(qry, qryParams...)
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
