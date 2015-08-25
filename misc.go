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
