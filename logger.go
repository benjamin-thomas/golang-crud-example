// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// https://raw.githubusercontent.com/gin-gonic/gin/develop/logger.go
package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

// type Pipeline struct {
//   curr http.Handler
//   next http.Handler
// }

// func (this Pipeline) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//   curr(w, r)
//   next(w, r)
// }

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *loggingResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	// Status is 200, unless overridden
	return &loggingResponseWriter{w, 200}
}

func appLogger(out io.Writer, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Start timer
		start := time.Now()
		path := r.URL.Path

		// Process request
		// c.Next()
		lw := newLoggingResponseWriter(w)
		h.ServeHTTP(lw, r)

		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		clientIP := r.RemoteAddr
		method := r.Method
		statusCode := lw.status
		statusColor := colorForStatus(statusCode)
		methodColor := colorForMethod(method)

		fmt.Fprintf(out, "[APP] %v |%s %3d %s| %12v | %s |%s  %s %-7s %s\n",
			end.Format("2006-12-30 - 15:04:05"),
			statusColor, statusCode, reset,
			latency,
			clientIP,
			methodColor, reset, method,
			path,
		)
	})
}

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code <= 299:
		return green
	case code >= 300 && code <= 399:
		return white
	case code >= 400 && code <= 499:
		return yellow
	default:
		return red
	}
}

func colorForMethod(method string) string {
	switch {
	case method == "GET":
		return blue
	case method == "POST":
		return cyan
	case method == "PUT":
		return yellow
	case method == "DELETE":
		return red
	case method == "PATCH":
		return green
	case method == "HEAD":
		return magenta
	case method == "OPTIONS":
		return white
	default:
		return reset
	}
}
