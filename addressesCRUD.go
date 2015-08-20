package main

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

func (as *addresses) index(sPer, sPage string) error {
	per, offset := paginateParams(sPer, sPage)
	rows, err := db.Query(`
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
				ON c.id  = a.city_id
		 INNER
			JOIN countries AS cn
				ON cn.id = c.country_id

		 ORDER BY id LIMIT $1 OFFSET $2
	`, &per, &offset)
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
