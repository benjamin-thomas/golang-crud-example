package main

import (
	"fmt"
	"log"
	"net/http"
)

func apiIndexAddresses(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()

	per := values.Get("per")
	page := values.Get("page")

	q := values.Get("q")
	op := values.Get("op")

	as := addresses{}

	err := as.index(per, page, q, op)
	if err != nil {
		log.Println(err)
		if _, ok := err.(*syntaxErr); ok {
			fmt.Println("ok =", ok)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			fmt.Println("ok =", ok)
			httpGenericErr(w)
		}
		return
	}

	renderJSON(w, as)
}

func indexAddresses(w http.ResponseWriter, r *http.Request) {
	renderHTML(w, nil, "addresses/index")
}
