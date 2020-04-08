// Package triangle determines if the 3 input variables can represent a triangle and if so which type.
package triangle

import (
	//	"fmt"
	"math"
)

//Inputs are float64 same as data type output
type Kind float64

const (
	// Pick values for the following identifiers used by the test program.
	NaT = iota // not a triangle
	Equ        // equilateral
	Iso        // isosceles
	Sca        // scalene
)

// KindFromSides recieves 3 variables as input establishes if it's a triangle or not and what kind of them.
func KindFromSides(a, b, c float64) Kind {
	//fmt.Println(a, b, c)
	var k Kind
	if (a <= 0 || b <= 0 || c <= 0) || (a+b < c) || (a+c < b) || (c+b < a) || math.IsNaN(a) || math.IsNaN(b) || math.IsNaN(c) || math.IsInf(a, 0) || math.IsInf(b, 0) || math.IsInf(c, 0) {
		k = 0 //not a triangle
	} else if a == b && b == c {
		k = 1 // equilateral
	} else if (a == b && a != c) || (a == c && a != b) || (b == c && b != a) {
		k = 2 // isosceles
	} else if a != b && b != c && a != c {
		k = 3 // scalene
	}
	return k
}
