package main

type countries []country

func (cs countries) Count() (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM countries").Scan(&count)
	return count, err
}

func (cs *countries) Index(per, page int) error {
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

	for rows.Next() {
		c := country{}
		err = rows.Scan(&c.Id, &c.Name, &c.HasStats)
		if err != nil {
			return err
		}
		*cs = append(*cs, c)
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	return nil
}

func (cs *countries) searchByName(name string, per, page int) error {
	// name = "%" + name + "%"
	offset := per * (page - 1)
	rows, err := db.Query(`
	SELECT c.id
			 , c.name
			 , (cs.id IS NOT NULL) AS has_stats
		FROM countries AS c
		LEFT
		JOIN country_stats AS cs
			ON cs.country_id = c.id
		WHERE c.name ILIKE $1 LIMIT $2 OFFSET $3
	`, &name, &per, &offset)
	defer rows.Close()

	for rows.Next() {
		c := country{}
		err = rows.Scan(&c.Id, &c.Name, &c.HasStats)
		if err != nil {
			return err
		}
		*cs = append(*cs, c)
	}
	err = rows.Err()
	if err != nil {
		return err
	}

	return nil
}

func (cs *countries) countByName(name string) (int, error) {
	var count int
	var err = db.QueryRow("SELECT COUNT(*) FROM countries WHERE name ILIKE $1", &name).Scan(&count)
	return count, err
}
