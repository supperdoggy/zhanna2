package main

import (
	"strings"
)

func split(text, pattern string) []string {
	return strings.Split(text, pattern)
}
