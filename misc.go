package main

import (
	"io/ioutil"
	"log"
	"strings"
)

func isAPIPath(path string) bool {
	if strings.HasPrefix(path, "/api/") {
		return true
	}
	return false
}

func mustReadFile(path string) []byte {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return b
}
