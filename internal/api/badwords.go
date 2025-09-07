package api

import "strings"

func StripBadWords(s string, r string, badWords []string) string {
	words := strings.Split(s, " ")
	for i := range words {
		for j := range badWords {
			if strings.ToLower(words[i]) == strings.ToLower(badWords[j]) {
				words[i] = r
			}
		}
	}

	return strings.Join(words, " ")
}
