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

	var (
		data struct { // data for the templates
			Countries  []country
			From       int
			To         int
			Total      int
			Path       string
			Pagination pagination
		}
		c country
	)

	var err error
	cs := countries{}

	data.Total, err = cs.Count()
	if err != nil {
		httpGenericErr(w)
		return
	}

	p, err := newPagination(per, page, data.Total)
	if err != nil {
		httpGenericErr(w)
		return
	}
	data.Pagination = p

	rows, err := cs.Index(p.Per, p.Page)
	if err != nil {
		httpGenericErr(w)
		return
	}
	defer rows.Close()

	data.Path = r.URL.Path

	// push that into cs.Index()
	for rows.Next() {
		err := rows.Scan(&c.Id, &c.Name)
		if err != nil {
			log.Fatal(err)
		}
		data.Countries = append(data.Countries, c)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	if strings.HasPrefix(r.URL.Path, "/api/") {
		j, err := json.Marshal(data.Countries)
		if err != nil {
			log.Fatal(err)
		}
		w.Header().Set("Content-Type", "application/json; charset=utf8")
		w.Write(j)
	} else {

		t := template.New("app.html")
		t.Funcs(template.FuncMap{
			"dec": func(a, b int) int {
				return b - a
			},
			"inc": func(a, b int) int {
				return a + b
			},
		})

		t, err = t.ParseFiles(
			"tmpl/layout/app.html",
			"tmpl/layout/pagination.html",
			"tmpl/countries/index.html",
		)
		if err != nil {
			log.Fatal(err)
		}

		err = t.Execute(w, data)
		if err != nil {
			log.Fatal(err)
		}
	}
}
