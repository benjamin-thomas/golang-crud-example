package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
)

func newCity(w http.ResponseWriter, r *http.Request, countryId int) {
	var tmplData struct {
		Country    country
		CancelPath string
	}

	c := country{Id: countryId}
	err := c.read()
	if err != nil {
		log.Println(err)
		httpGenericErr(w)
		return
	}

	tmplData.Country = c
	tmplData.CancelPath = strings.Replace(r.URL.Path, "/new", "", -1)
	renderHTML(w, tmplData, "cities/new")
}

func createCity(w http.ResponseWriter, r *http.Request) {
	strCountryId := r.FormValue("countryId")
	name := r.FormValue("name")

	countryId, err := strconv.Atoi(strCountryId)
	if err != nil {
		log.Println(err)
		httpGenericErr(w)
		return
	}

	c := &city{
		CountryId: countryId,
		Name:      name,
	}

	err = c.create()
	if err != nil {
		log.Println(err)
		httpGenericErr(w)
		return
	}

}

func showCity(w http.ResponseWriter, r *http.Request, id int) {
	c := &city{Id: id}
	err := c.read()
	if err != nil {
		log.Println(err)
		httpGenericErr(w)
		return
	}

	if isAPIPath(r.URL.Path) {
		renderJSON(w, c)
	} else {
		renderHTML(w, c, "cities/show")
	}
}
