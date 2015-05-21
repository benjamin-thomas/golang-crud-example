package main

type city struct {
	Id        int    `json:"id"`
	CountryId int    `json:"countryId"`
	Name      string `json:"name"`
}

func (c *city) create() error {
	_, err := db.Exec("INSERT INTO cities (country_id, name) VALUES ($1, $2)", c.CountryId, c.Name)
	return err
}

func (c *city) read() error {
	return db.QueryRow("SELECT name FROM cities WHERE id = $1", &c.Id).Scan(&c.Name)
}
