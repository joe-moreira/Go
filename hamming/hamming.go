// Package hamming calculates the Hamming distance between two DNA strands
package hamming

import (
	"errors"
)

// Distance delivers the distance number between two DNA strands
func Distance(a, b string) (c int, d error) {
	distance := 0

	ar := []rune(a)
	br := []rune(b)

	if len(ar) != len(br) {
		//		fmt.Println("DNA strands have different number of components.")
		return 0, errors.New("strands have different number of components")
	}
	for i := range ar {
		if ar[i] != br[i] {
			distance++
		}

	}
	return distance, d
}
