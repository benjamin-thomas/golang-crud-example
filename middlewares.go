package main

import (
	"net/http"
	"strings"
)

func redirectOnTrailingSlash(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path != "/" && strings.HasSuffix(path, "/") {
			http.Redirect(w, r, path[0:len(path)-1], http.StatusMovedPermanently)
			return
		}
		h.ServeHTTP(w, r)
	})
}
