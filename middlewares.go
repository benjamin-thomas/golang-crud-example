package main

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/handlers"
)

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
