package main

import (
	"sort"
	"strings"
)

func main() {

}

func newAnnoDict(words []string) map[string][]string {
	dict := make(map[string][]string)

	for _, word := range words {
		word = strings.ToLower(word)
		key := convertStringToKey(word)
		if _, ok := dict[key]; !ok {
			dict[key] = []string{word}
		} else {
			if !isStringInSet(word, dict[key]) {
				dict[key] = append(dict[key], word)
			}
		}
	}

	return sortMapValues(removeOneValueKeys(dict))
}

func isStringInSet(str string, strs []string) bool {
	for _, v := range strs {
		if v == str {
			return true
		}
	}
	return false
}

func convertStringToKey(word string) string {
	letters := strings.Split(word, "")
	sort.Strings(letters)

	return strings.Join(letters, "")
}

func sortMapValues(dict map[string][]string) map[string][]string {
	for k := range dict {
		sort.Strings(dict[k])
	}
	return dict
}

func removeOneValueKeys(dict map[string][]string) map[string][]string {
	for k, v := range dict {
		if len(v) <= 1 {
			delete(dict, k)
		}
	}
	return dict
}
