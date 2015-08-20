package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func indexAddresses(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()

	per := values.Get("per")
	page := values.Get("page")

	q := values.Get("q")
	fmt.Println("q =", q)

	sCols := values.Get("cols")
	op := values.Get("op")

	var cols []string
	json.Unmarshal([]byte(sCols), &cols)

	fmt.Printf("cols = %#v\n", cols)
	for _, c := range cols {
		fmt.Println("c =", c)
	}

	as := addresses{}

	err := as.index(per, page, q, op, cols)
	if err != nil {
		log.Println(err)
		httpGenericErr(w)
		return
	}

	renderJSON(w, as)
}
