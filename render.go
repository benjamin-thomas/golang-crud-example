package main

import (
	"encoding/json"
	"log"
	"net/http"
	"text/template"
)

func renderJSON(w http.ResponseWriter, data interface{}) {
	j, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		httpGenericErr(w)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf8")
	w.Write(j)
}

func renderHTML(w http.ResponseWriter, data interface{}, tmplPathName string) {
	t, err := template.ParseFiles(
		"tmpl/layout/app.html",
		"tmpl/layout/pagination.html",
		"tmpl/"+tmplPathName+".html",
	)
	if err != nil {
		log.Println(err)
		httpGenericErr(w)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		log.Println(err)
		httpGenericErr(w)
		return
	}
}
