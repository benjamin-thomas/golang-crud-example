package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
)

type country struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func showCountry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Fatal(err)
	}

	var (
		name string
	)

	err = stmtGetCountry.QueryRow(&id).Scan(&name)
	if err == sql.ErrNoRows {
		fmt.Fprintf(w, "No country with id: %d", id)
		return
	}
	if err != nil {
		log.Println("Failed on db fetch:", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Oops, something went wrong")
		return
	}

	c := country{id, name}
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

	_, err := crud.countries.update(id, name)
	if err != nil {
		log.Println("Err:", err)
		http.Error(w, fmt.Sprintf("Could not create country with name: %s", name), http.StatusInternalServerError)
	}
}

func createCountry(w http.ResponseWriter, r *http.Request) {
	var name = r.FormValue("name")

	_, err := crud.countries.create(name)
	if err != nil {
		fmt.Println("Err:", err)
		http.Error(w, fmt.Sprintf("Could not create country with name: %s", name), http.StatusInternalServerError)
	}
}

func deleteCountry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]
	stmt := mustPrepare("DELETE FROM countries WHERE id = $1")
	res, err := stmt.Exec(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not delete country with id: '%s'", id), http.StatusInternalServerError)
		return
	}
	affected, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	if affected == 0 {
		http.Error(w, fmt.Sprintf("Country '%s' already deleted", id), http.StatusInternalServerError)
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

func editCountry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var id = vars["id"]

	var c = country{}
	var stmt = mustPrepare("SELECT id, name FROM countries WHERE id = $1 ")
	var err = stmt.QueryRow(id).Scan(&c.Id, &c.Name)

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
			Countries     []country
			From          int
			To            int
			Total         int
			LastPage      int
			Page          int
			Path          string
			HideFirstLink bool
			HideLastLink  bool
			NextPage      int
			PrevPage      int
			ValidNextPage bool
			Per           int
		}
		c country
	)

	stmt := mustPrepare("SELECT COUNT(*) FROM countries;")
	err := stmt.QueryRow().Scan(&data.Total)
	if err != nil {
		log.Fatal(err)
	}

	data.Per = mustAtoi(per)
	data.Page = mustAtoi(page)
	if data.Page < 1 {
		data.Page = 1
	}
	data.LastPage = (data.Total / data.Per) + 1
	offset := data.Per * (data.Page - 1)

	stmt = mustPrepare(`
	SELECT  c.id
				, c.name
	FROM countries AS c
	ORDER BY c.name ASC
	LIMIT $1 OFFSET $2
	`)
	rows, err := stmt.Query(&data.Per, &offset)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	data.PrevPage = data.Page - 1
	if data.PrevPage < 1 {
		data.PrevPage = 1
	}
	data.NextPage = data.Page + 1
	if data.NextPage > data.LastPage {
		data.NextPage = data.LastPage
	}

	data.ValidNextPage = data.NextPage <= data.LastPage
	data.HideFirstLink = data.Page <= 1
	data.HideLastLink = data.Page >= data.LastPage

	data.Path = r.URL.Path

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
