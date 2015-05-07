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

	limit := r.URL.Query().Get("limit")
	fmt.Println("limit =", limit)
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

func indexCountries(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	limit := values.Get("limit")
	if limit == "" {
		limit = defaultLimit
	}
	offset := values.Get("offset")
	if offset == "" {
		offset = defaultOffset
	}

	var (
		data struct { // data for the templates
			Countries  []country
			Offset     int
			PrevOffset int
			NextOffset int
			Limit      int
			From       int
			To         int
			Total      int
			Path       string
		}
		c country
	)

	stmt := mustPrepare("SELECT COUNT(*) FROM countries;")
	err := stmt.QueryRow().Scan(&data.Total)
	if err != nil {
		log.Fatal(err)
	}

	stmt = mustPrepare(`
	SELECT  c.id
				, c.name
	FROM countries AS c
	ORDER BY c.name ASC
	LIMIT $1 OFFSET $2
	`)
	rows, err := stmt.Query(&limit, &offset)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	data.Path = r.URL.Path
	data.Limit, err = strconv.Atoi(limit)
	if err != nil {
		log.Fatal(err)
	}
	data.Offset, err = strconv.Atoi(offset)
	if err != nil {
		log.Fatal(err)
	}
	data.NextOffset = data.Offset + data.Limit

	if data.Offset == 0 {
		data.PrevOffset = 0
	} else {
		data.PrevOffset = data.Offset - data.Limit
	}

	data.From = data.Offset + 1
	data.To = data.NextOffset
	if data.NextOffset > data.Total {
		data.NextOffset = data.Offset
		data.To = data.Total
	}

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
		t, err := template.ParseFiles(
			"tmpl/layout/app.html",
			"tmpl/layout/pagination.html",
			"tmpl/countries/index.html",
		)
		if err != nil {
			log.Fatal(err)
		}
		t.Execute(w, data)
	}
}
