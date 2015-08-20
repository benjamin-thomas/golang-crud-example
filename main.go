package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/k0kubun/pp"
	_ "github.com/lib/pq"
)

var (
	db               *sql.DB
	defaultPer       int = 10
	stmtGetCountry   *sql.Stmt
	stmtGetCountries *sql.Stmt
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/countries", http.StatusFound)
}

// func contractsHandler(w http.ResponseWriter, r *http.Request, key string) {
//   fmt.Println("contracts!!, key = ", key)
//   w.Write([]byte("contracts handler"))
//   keys := strings.Split(key, "/contracts/")
//   clientID := keys[0]
//   contractID := keys[1]
//   fmt.Println("keys =", keys)
//   fmt.Println("clientID =", clientID)
//   fmt.Println("contractID =", contractID)
// }

// func showContract(w http.ResponseWriter, r *http.Request) {
//   fmt.Println("r.URL.Path =", r.URL.Path)
//   fmt.Println("r.URL.Query() =", r.URL.Query())
//   fmt.Fprintln(w, "showContract:", r.URL.Path)
//   client_id := r.URL.Query().Get(":client_id")
//   fmt.Println("client_id =", client_id)

//   contract_id := r.URL.Query().Get(":contract_id")
//   fmt.Println("contract_id =", contract_id)
// }

func listReports(w http.ResponseWriter, r *http.Request) {
	var (
		id       int
		filename string
	)
	rows, err := db.Query("SELECT id, filename FROM reports")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &filename)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintln(w, "id =", id, ", ", "filename =", filename)
	}
}

func setupDB() {
	dsn := fmt.Sprintf("host=%s dbname=%s user=%s port=%s sslmode=disable",
		os.Getenv("PGHOST"),
		os.Getenv("PGDATABASE"),
		os.Getenv("PGUSER"),
		os.Getenv("PGPORT"),
	)
	log.Println("dsn =", dsn)

	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func injectKey(fn func(http.ResponseWriter, *http.Request, string), path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "*****path=%s******\r\n", path)
		key := r.URL.Path[len(path):]
		fn(w, r, key)
	}
}

func route(fn func(http.ResponseWriter, *http.Request, string), basePath string) http.HandlerFunc {
	pp.Println("basePath =", basePath)

	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Path[len(basePath):]
		fn(w, r, key)
	}
}

func timerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//fmt.Println("Before:", time.Now())
		next.ServeHTTP(w, r)
		//fmt.Println("After:", time.Now())
	})
}

func mustGetenv(env string) string {
	s := os.Getenv(env)
	if s == "" {
		panic(fmt.Sprintf("%s env var not set", env))
	}
	return s
}

func stringKeyProvider(key string, fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		k, ok := vars[key]
		if !ok {
			log.Printf("stringKeyProvider: key '%s' not found\n", key)
			httpGenericErr(w)
		} else {
			fn(w, r, k)
		}
	}
}

func intKeyProvider(key string, fn func(http.ResponseWriter, *http.Request, int)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		k, ok := vars[key]
		if !ok {
			log.Printf("intKeyProvider: key '%s' not found\n", key)
			httpGenericErr(w)
		} else {
			n, err := strconv.Atoi(k)
			if err != nil {
				log.Println(err)
				httpGenericErr(w)
				return
			}

			fn(w, r, n)
		}
	}
}

// Injects middlewares while routing
type routerHelper struct {
	mux             *mux.Router
	middlewareChain func(h http.Handler) http.Handler
}

func (r *routerHelper) HandleFunc(pat string, h http.HandlerFunc) *mux.Route {
	return r.mux.Handle(pat, r.middlewareChain(h))
}

func main() {
	setupDB()
	defer db.Close()

	mux := mux.NewRouter()
	r := &routerHelper{
		mux:             mux,
		middlewareChain: middlewares,
	}

	/*
		rr := &routerHelper{
			mux: mux,
			middlewareChain: func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprintln(w, "rr middleware chain!!")
					next.ServeHTTP(w, r)
				})
			},
		}

		r.mux.HandleFunc("/countries/new", newCountry).Methods("GET") // no middlewares
		rr.mux.HandleFunc("/countries/new", newCountry).Methods("GET") // another middleware
	*/

	// http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	// http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("public/assets"))))
	http.Handle("/assets/", middlewares(http.FileServer(http.Dir("public"))))
	// http.Handle("/assets/", http.StripPrefix("public/assets", http.FileServer(http.Dir("public/assets"))))

	if os.Getenv("DEV") == "1" {
		http.Handle("/assets_dev/", http.StripPrefix("/assets_dev/", http.FileServer(http.Dir("assets"))))
	}

	r.HandleFunc("/", rootHandler).Methods("GET")

	// r.HandleFunc("/countries", indexCountries).Methods("GET")
	// r.HandleFunc("/countries/{id}", intKeyProvider("id", showCountry)).Methods("GET")
	// r.HandleFunc("/countries", createCountry).Methods("POST")
	// r.HandleFunc("/countries/{id}", intKeyProvider("id", updateCountry)).Methods("PUT", "PATCH")

	r.HandleFunc("/addresses", indexAddresses).Methods("GET")

	// r.HandleFunc("/countries/{id}/contracts/new", newCountry).Methods("GET")
	// r.HandleFunc("/countries/{id}/stats", intKeyProvider("id", showCountryStats)).Methods("GET")

	// r.HandleFunc("/countries/{id}/cities", intKeyProvider("id", indexCountryCities)).Methods("GET")

	// r.HandleFunc("/countries/{country_id}/cities/new", intKeyProvider("country_id", newCity)).Methods("GET")
	// r.HandleFunc("/cities/{id}/edit", intKeyProvider("id", editCity)).Methods("GET")
	// r.HandleFunc("/cities/{id}", intKeyProvider("id", showCity)).Methods("GET")
	// r.HandleFunc("/cities/{id}", intKeyProvider("id", updateCity)).Methods("PUT", "PATCH")
	// r.HandleFunc("/cities", createCity).Methods("POST")

	http.Handle("/", redirectOnTrailingSlash(r.mux))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
