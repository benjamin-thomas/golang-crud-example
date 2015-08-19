package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"text/template"
)

func renderJSON(w http.ResponseWriter, data interface{}) {
	var j []byte
	var err error

	if os.Getenv("DEV") == "1" {
		j, err = json.MarshalIndent(data, "", " ")
	} else {
		j, err = json.Marshal(data)
	}

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

	// Avoid partial rendering
	b := &bytes.Buffer{}
	err = t.Execute(b, data)
	if err != nil {
		log.Println(err)
		httpGenericErr(w)
		return
	}
	_, err = b.WriteTo(w)
	if err != nil {
		log.Fatalf("renderHTML: %s\n", err)
	}
}
