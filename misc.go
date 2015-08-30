package main

import (
	"strings"
)

func isAPIPath(path string) bool {
	if strings.HasPrefix(path, "/api/") {
		return true
	}
	return false
}
