package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

func showCountry(w http.ResponseWriter, r *http.Request, id string) {
	c := &country{Id: id}
	err := c.read()
	if err != nil {
		log.Println(err)
		httpGenericErr(w)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/api/") {
		j, err := json.Marshal(c)
		if err != nil {
			log.Fatal(err)
		}
		w.Header().Set("Content-Type", "application/json; charset=utf8")
		w.Write(j)
	} else {
		t, _ := template.ParseFiles(
			"tmpl/layout/app.html",
			"tmpl/countries/show.html",
		)
		err = t.Execute(w, c)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func updateCountry(w http.ResponseWriter, r *http.Request, id string) {
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

func deleteCountry(w http.ResponseWriter, r *http.Request, id string) {
	c := &country{Id: id}
	err := c.delete()
	if err != nil {
		log.Println(err)
		httpGenericErr(w)
		return
	}
}

func destroyCountry(w http.ResponseWriter, r *http.Request, key string) {
	fmt.Fprintln(w, "destroyCountry:", r.URL.Path, "key=", key)
	fmt.Println("r.URL.Query() =", r.URL.Query())
}

func newCountry(w http.ResponseWriter, r *http.Request) {
	var t, err = template.ParseFiles(
		"tmpl/layout/app.html",
		"tmpl/countries/new.html",
	)
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func listCountry(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "listCountry:", r.URL.Path)
	fmt.Println("r.URL.Query() =", r.URL.Query())
}

func editCountry(w http.ResponseWriter, r *http.Request, id string) {
	var c = country{Id: id}
	err := c.read()

	var t *template.Template
	t, err = template.ParseFiles(
		"tmpl/layout/app.html",
		"tmpl/countries/edit.html",
	)
	if err != nil {
		log.Fatal(err)
	}

	err = t.Execute(w, c)
	if err != nil {
		log.Fatal(err)
	}
}

func mustAtoi(s string) int {
	res, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}
	return res
}

// hideFirstLink, hideLastLink bool, path string

func indexCountries(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	page := values.Get("page")
	if page == "" {
		page = "1"
	}
	per := values.Get("per")
	if per == "" {
		per = defaultPer
	}

	var tmplData struct {
		Path       string
		Pagination pagination
		Countries  []country
	}
	tmplData.Path = r.URL.Path

	cs := countries{}

	count, err := cs.Count()
	if err != nil {
		httpGenericErr(w)
		return
	}

	p, err := newPagination(per, page, count)
	if err != nil {
		httpGenericErr(w)
		return
	}
	tmplData.Pagination = p

	tmplData.Countries, err = cs.Index(p.Per, p.Page)
	if err != nil {
		httpGenericErr(w)
		return
	}

	if strings.HasPrefix(r.URL.Path, "/api/") {
		j, err := json.Marshal(tmplData.Countries)
		if err != nil {
			log.Fatal(err)
		}
		w.Header().Set("Content-Type", "application/json; charset=utf8")
		w.Write(j)
	} else {

		t, err := template.ParseFiles(
			"tmpl/layout/app.html",
			"tmpl/layout/pagination.html",
			"tmpl/countries/index.html",
		)
		if err != nil {
			log.Fatal(err)
		}

		err = t.Execute(w, tmplData)
		if err != nil {
			log.Fatal(err)
		}
	}
}
