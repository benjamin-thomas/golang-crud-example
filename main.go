package main

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/k0kubun/pp"
	_ "github.com/lib/pq"
)

var (
	db               *sql.DB
	defaultLimit     string = "10"
	defaultPer       string = "10"
	defaultOffset    string = "0"
	stmtGetCountry   *sql.Stmt
	stmtGetCountries *sql.Stmt
)

func mustPrepare(qry string) *sql.Stmt {
	stmt, err := db.Prepare(qry)
	if err != nil {
		log.Fatal(err)
	}
	return stmt
}

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

func setupStmts() {
	stmtGetCountry = mustPrepare("SELECT name FROM countries WHERE id = $1")
	stmtGetCountries = mustPrepare("SELECT id, name FROM countries LIMIT $1 OFFSET $2")
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

func BasicAuth(user, pass []byte, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const basicAuthPrefix string = "Basic "

		// Get the Basic Authentication credentials
		auth := r.Header.Get("Authorization")
		if strings.HasPrefix(auth, basicAuthPrefix) {
			// Check credentials
			payload, err := base64.StdEncoding.DecodeString(auth[len(basicAuthPrefix):])
			if err == nil {
				pair := bytes.SplitN(payload, []byte(":"), 2)
				if len(pair) == 2 &&
					bytes.Equal(pair[0], user) &&
					bytes.Equal(pair[1], pass) {

					// Delegate request to the given handle
					h.ServeHTTP(w, r)
					return
				}
			}
		}

		// Request Basic Authentication otherwise
		w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	})
}

func mustGetenv(env string) string {
	s := os.Getenv(env)
	if s == "" {
		panic(fmt.Sprintf("%s env var not set", env))
	}
	return s
}

func middlewares(final http.Handler) http.Handler {
	chain := func(h http.Handler) http.Handler {
		return h
	}

	user := []byte(mustGetenv("CRUD_USER"))
	pw := []byte(mustGetenv("CRUD_PW"))
	authenticate := func(h http.Handler) http.Handler {
		return BasicAuth(user, pw, h)
	}

	dev := true
	logger := func(h http.Handler) http.Handler {
		if dev {
			// gin logger
			return appLogger(os.Stdout, h)
		} else {
			return handlers.LoggingHandler(os.Stdout, h)
		}
	}

	return chain(
		timerMiddleware(
			authenticate(
				logger(final),
			),
		),
	)
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

	setupStmts()

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
	r.HandleFunc("/", rootHandler).Methods("GET")
	r.HandleFunc("/countries/new", newCountry).Methods("GET")
	r.HandleFunc("/countries/{id}/edit", editCountry).Methods("GET")

	r.HandleFunc("/api/countries", indexCountries).Methods("GET")
	r.HandleFunc("/countries", indexCountries).Methods("GET")
	r.HandleFunc("/countries/{id}", showCountry).Methods("GET")
	r.HandleFunc("/api/countries/{id}", showCountry).Methods("GET")
	r.HandleFunc("/countries", createCountry).Methods("POST")
	r.HandleFunc("/countries/{id}", updateCountry).Methods("PUT", "PATCH")
	r.HandleFunc("/api/countries/{id}", deleteCountry).Methods("DELETE")

	r.HandleFunc("/countries/{id}/contracts", listCountry).Methods("GET")
	r.HandleFunc("/countries/{id}/contracts/new", newCountry).Methods("GET")
	r.HandleFunc("/countries/{id}/stats", showCountry).Methods("GET")

	http.Handle("/", r.mux)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
