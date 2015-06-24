package main

type zipCode struct {
	Id   int `json: "id"`
	Code string
}

type city struct {
	Id int `json: "id"`

	Country country
	ZipCode zipCode
	Name    string `json:"name"`
}

func (c *city) create() error {
	_, err := db.Exec("INSERT INTO cities (country_id, name) VALUES ($1, $2)", c.Country.Id, c.Name)
	return err
}

func (c *city) read() error {
	return db.QueryRow(`
	SELECT c.name
	     , co.id
	     , co.name
			 , z.id
			 , z.code
	  FROM cities c
	  LEFT
	  JOIN countries co
	    ON co.id = c.country_id
		LEFT
		JOIN zip_codes z
		  ON z.city_id = c.id
	 WHERE c.id    = $1
	`, &c.Id).Scan(&c.Name, &c.Country.Id, &c.Country.Name, &c.ZipCode.Id, &c.ZipCode.Code)
}

func (c *city) update() error {
	_, err := db.Exec("UPDATE cities SET name = $2 WHERE id = $1", c.Id, c.Name)
	return err
}

func (c *city) delete() error {
	_, err := db.Exec("DELETE FROM cities WHERE id = $1", c.Id)
	return err
}
