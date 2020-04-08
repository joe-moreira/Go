// Package bob establish a dialoge with a user
package scale

import (
	"regexp"
	"strings"
)
// Hey delivers a response based on the question provided
func Scale(tonic string, interval string) string {
	switch {
	case tonic == "C":
		return ([]"C","C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B")
	default:
		 return ""
}
}