package raindrops

import (
	"strconv"
)

// Convert receives a number (int) and returns a result (string) based on its factor.
func Convert(number int) (result string) {
	if number%3 == 0 {
		result += "Pling"
	}
	if number%5 == 0 {
		result += "Plang"
	}
	if number%7 == 0 {
		result += "Plong"
	} else if number%3 != 0 && number%5 != 0 && number%7 != 0 {
		return strconv.Itoa(number)
	}
	return result
}
