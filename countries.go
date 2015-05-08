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
		t, _ := template.ParseFiles("tmpl/countries/show.html")
		t.Execute(w, c)
	}
}

func updateCountry(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "updateCountry:", r.URL.Path, "r.Method =", r.Method)
	fmt.Println("r.URL.Query() =", r.URL.Query())
}

func createCountry(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "createCountry:", r.URL.Path)
	fmt.Println("r.URL.Query() =", r.URL.Query())
}

func deleteCountry(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "deleteCountry:", r.URL.Path)
	fmt.Println("r.URL.Query() =", r.URL.Query())
}

func destroyCountry(w http.ResponseWriter, r *http.Request, key string) {
	fmt.Fprintln(w, "destroyCountry:", r.URL.Path, "key=", key)
	fmt.Println("r.URL.Query() =", r.URL.Query())
}

func newCountry(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "newCountry:", r.URL.Path)
	fmt.Println("r.URL.Query() =", r.URL.Query())
	fmt.Fprintf(w, "NOOP new")
}

func listCountry(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "listCountry:", r.URL.Path)
	fmt.Println("r.URL.Query() =", r.URL.Query())
}

func editCountry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprintln(w, "vars =", vars)
	fmt.Fprintf(w, "editCountry: path=%s\n", r.URL.Path)
	fmt.Println("r.URL.Query() =", r.URL.Query())
	fmt.Fprintf(w, "NOOP edit\n")
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
			Path          string
			Page          int
			NextPage      int
			PrevPage      int
			ValidNextPage bool
			ValidPrevPage bool
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

	data.Path = r.URL.Path

	data.PrevPage = data.Page - 1
	data.NextPage = data.Page + 1

	data.ValidPrevPage = data.PrevPage > 0
	data.ValidNextPage = data.NextPage <= data.LastPage

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

		// t, err := template.ParseFiles(
		//   "tmpl/layout/app.html",
		//   "tmpl/layout/pagination.html",
		//   "tmpl/countries/index.html",
		// )
		if err != nil {
			log.Fatal(err)
		}
		t.Execute(w, data)
	}
}
