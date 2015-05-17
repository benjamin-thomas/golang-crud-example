package main

import (
	"log"
	"net/http"
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

	var tmplData struct {
		Path       string
		Pagination pagination
		Countries  []country
	}
	tmplData.Path = r.URL.Path

	cs := countries{}
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

	tmplData.Countries, err = cs.Index(p.Per, p.Page)
	if err != nil {
		log.Println(err)
		httpGenericErr(w)
		return
	}

	if isAPIPath(r.URL.Path) {
		renderJSON(w, tmplData.Countries)
	} else {
		renderHTML(w, tmplData, "countries/index")
	}
}
