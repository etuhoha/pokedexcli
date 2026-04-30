package main

import "strings"

func cleanInput(text string) []string {
	var result []string
	for _, s := range strings.Split(strings.TrimSpace(text), " ") {
		result = append(result, strings.ToLower(s))
	}
	return result
}
