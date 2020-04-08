package proverb

import "fmt"

func Proverb(rhyme []string) []string {
	fmt.Println(len(rhyme))
	if len(rhyme) == 1 {
		proverb := []string{"And all for the want of a nail."}
		return proverb
	} else if len(rhyme) > 1 {
		for i := 0; i < len(rhyme); i++ {
			//	fmt.Println(rhyme[i])
			fmt.Printf("For want of a %v the %v was lost.", rhyme[i], rhyme[i+1])
			fmt.Printf("And all for the want of a %v.", rhyme[i])
		}
	} else {
		return []string{}
	}
	return []string{}
}
