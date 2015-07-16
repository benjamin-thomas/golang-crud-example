package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/k0kubun/pp"
)

func showCountry(w http.ResponseWriter, r *http.Request, id int) {
	c := &country{Id: id}
	err := c.read()
	if err != nil {
		log.Println(err)
		httpGenericErr(w)
		return
	}
	if isAPIPath(r.URL.Path) {
		renderJSON(w, c)
	} else {
		renderHTML(w, c, "countries/show")
	}
}

func showCountryStats(w http.ResponseWriter, r *http.Request, id int) {
	c := &country{Id: id}
	err := c.stats()
	if err != nil {
		log.Println(err)
		httpGenericErr(w)
		return
	}
	if isAPIPath(r.URL.Path) {
		renderJSON(w, c)
	} else {
		renderHTML(w, c, "countries/stats")
	}
}
func updateCountry(w http.ResponseWriter, r *http.Request, id int) {
	var name = r.FormValue("name")

	c := &country{Id: id, Name: name}
	err := c.update()
	if err != nil {
		log.Println(err)
		httpGenericErr(w)
		return
	}
}

func createCountry(w http.ResponseWriter, r *http.Request) {
	var name = r.FormValue("name")

	c := &country{Name: name}
	err := c.create()
	if err != nil {
		log.Println(err)
		httpGenericErr(w)
		return
	}
}

func deleteCountry(w http.ResponseWriter, r *http.Request, id int) {
	c := &country{Id: id}
	err := c.delete()
	if err != nil {
		log.Println(err)
		httpGenericErr(w)
		return
	}
}

func newCountry(w http.ResponseWriter, r *http.Request) {
	renderHTML(w, nil, "countries/new")
}

func editCountry(w http.ResponseWriter, r *http.Request, id int) {
	var c = country{Id: id}
	err := c.read()
	if err != nil {
		log.Println(err)
		httpGenericErr(w)
		return
	}

	renderHTML(w, c, "countries/edit")
}

func indexCountries(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	page := values.Get("page")
	per := values.Get("per")

	cs := countries{}

	var tmplData struct {
		Path       string
		Pagination pagination
		Countries  countries
	}
	tmplData.Path = r.URL.Path

	count, err := cs.Count()
	if err != nil {
		log.Println(err)
		httpGenericErr(w)
		return
	}

	p, err := newPagination(per, page, count)
	if err != nil {
		log.Println(err)
		httpGenericErr(w)
		return
	}
	tmplData.Pagination = p

	err = cs.Index(p.Per, p.Page)
	if err != nil {
		log.Println(err)
		httpGenericErr(w)
		return
	}
	tmplData.Countries = cs

	if isAPIPath(r.URL.Path) {
		renderJSON(w, tmplData.Countries)
	} else {
		renderHTML(w, tmplData, "countries/index")
	}
}

func asArray(cs countries) [][]interface{} {
	rows := [][]interface{}{}
	for _, c := range cs {
		rows = append(rows, []interface{}{c.Id, c.Name, c.HasStats, "todo"})
	}
	return rows
}

func datatableCountries(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	pp.Println(values)
	start := values.Get("start")
	length := values.Get("length")

	iStart, err := strconv.Atoi(start)
	if err != nil {
		log.Println(err)
		httpGenericErr(w)
	}

	per, err := strconv.Atoi(length)
	if err != nil {
		log.Println(err)
		httpGenericErr(w)
	}
	page := (iStart / per) + 1 // start=0 => page=1; start=10 => page=2; etc.

	draw := values.Get("draw")
	pp.Println("draw =", draw)
	println("page =", page)
	println("per =", per)
	q := values.Get("search[value]")
	if q == "" {
		// Try column search
		q = values.Get("columns[1][search][value]")
	}
	q = "%" + q + "%"
	pp.Println("q =", q)
	var tmplData struct {
		Draw            int             `json:"draw"`
		RecordsTotal    int             `json:"recordsTotal"`
		RecordsFiltered int             `json:"recordsFiltered"`
		Data            [][]interface{} `json:"data"`
	}

	iDraw, err := strconv.Atoi(draw)
	if err != nil {
		log.Println(err)
		httpGenericErr(w)
		return
	}
	tmplData.Draw = iDraw

	cs := countries{}
	count, err := cs.Count()
	if err != nil {
		log.Println(err)
		httpGenericErr(w)
		return
	}

	p, err := newPagination(strconv.Itoa(per), strconv.Itoa(page), count)
	if err != nil {
		log.Println(err)
		httpGenericErr(w)
		return
	}
	tmplData.RecordsTotal = p.Count

	err = cs.searchByName(q, per, page)
	if err != nil {
		log.Println(err)
		httpGenericErr(w)
		return
	}

	cnt, err := cs.countByName(q)
	if err != nil {
		log.Println(err)
		httpGenericErr(w)
		return
	}

	tmplData.Data = asArray(cs)
	tmplData.RecordsFiltered = cnt
	pp.Println("tmplData =", tmplData)

	renderJSON(w, tmplData)
}

func indexCountryCities(w http.ResponseWriter, r *http.Request, id int) {
	values := r.URL.Query()
	page := values.Get("page")
	per := values.Get("per")

	var tmplData struct {
		Path       string
		Pagination pagination
		Country    *country
	}
	tmplData.Path = r.URL.Path

	c := &country{Id: id}

	count, err := c.CitiesCount()
	if err != nil {
		log.Println(err)
		httpGenericErr(w)
		return
	}

	p, err := newPagination(per, page, count)
	if err != nil {
		log.Println(err)
		httpGenericErr(w)
		return
	}
	tmplData.Pagination = p

	err = c.indexCities(p.Per, p.Page)
	if err != nil {
		log.Println(err)
		httpGenericErr(w)
		return
	}
	tmplData.Country = c

	if isAPIPath(r.URL.Path) {
		renderJSON(w, tmplData)
	} else {
		renderHTML(w, tmplData, "countries/indexCities")
	}
}
