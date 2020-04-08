package acronym

import (
	"fmt"
	//	"unicode"
	"regexp"
)

// Abbreviate returns acronym for given phrase
func Abbreviate(s string) string {
	regex := regexp.MustCompile(`[a-zA-Z']+`).FindAllString(s, -1)
	fmt.Println(regex)
	return ""
}
