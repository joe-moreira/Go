// Package twofer distributes the responses between two participants based on their names.
package twofer

import "fmt"

//ShareWith recieves the participant's name as the input and prints on screen the distribution as the output.
func ShareWith(name string) string {
	if name == "" {
		name = "you"
	}
	return fmt.Sprintf("One for %s, one for me.", name)
}
