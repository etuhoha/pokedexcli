package main

import "strings"

func cleanInput(text string) []string {
	var result []string
	for s := range strings.SplitSeq(text, " ") {
		if len(s) > 0 {
			result = append(result, strings.ToLower(s))
		}
	}
	return result
}
