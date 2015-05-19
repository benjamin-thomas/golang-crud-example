package main

type countries struct {
}

func (c countries) Name() string {
	return "Countries"
}

func (c countries) Count() (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM countries").Scan(&count)
	return count, err
}

func (c *countries) Index(per, page int) ([]country, error) {
	offset := per * (page - 1)
	rows, err := db.Query(`
										SELECT c.id
												 , c.name
												 , (cs.id IS NOT NULL) AS has_stats
											FROM countries AS c
											LEFT
											JOIN country_stats AS cs
												ON cs.country_id = c.id
										 ORDER BY c.name ASC LIMIT $1 OFFSET $2
												 `, &per, &offset)
	defer rows.Close()

	var cs []country

	for rows.Next() {
		c := country{}
		err = rows.Scan(&c.Id, &c.Name, &c.HasStats)
		if err != nil {
			return cs, err
		}
		cs = append(cs, c)
	}

	err = rows.Err()
	if err != nil {
		return cs, err
	}

	return cs, nil
}
