package main

import (
	"log"
	"net/http"
)

func indexAddresses(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()

	per := values.Get("per")
	page := values.Get("page")

	as := addresses{}

	err := as.index(per, page)
	if err != nil {
		log.Println(err)
		httpGenericErr(w)
		return
	}

	renderJSON(w, as)
}
