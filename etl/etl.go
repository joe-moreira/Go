package etl

import (
	"strings"
)

// Transform recieves map[int][]string and returns map[string]int
func Transform(input map[int][]string) map[string]int {

	output := make(map[string]int)

	for k, v := range input {
		for _, vv := range v {
			output[strings.ToLower(vv)] = k
		}
	}
	return output
}
