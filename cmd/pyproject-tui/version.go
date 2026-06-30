package main

import (
	"os"
	"strings"
)

func versionFromFile() string {
	data, err := os.ReadFile("VERSION")
	if err != nil {
		return "dev"
	}
	v := strings.TrimSpace(string(data))
	if v == "" {
		return "dev"
	}
	return v
}
