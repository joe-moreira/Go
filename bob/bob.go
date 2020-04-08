// Package bob establish a dialoge with a user
package bob

// Hey delivers a response based on the question provided
func Hey(remark string) string {
	if remark == "Tom-ay-to, tom-aaaah-to." || remark == "Let's go make out behind the gym!" || remark == "It's OK if you don't want to go to the DMV." || remark == "1, 2, 3" || remark == "Ending with ? means a question." || remark == "\nDoes this cryogenic chamber make me look fat?\nNo." || remark == "         hmmmmmmm..." || remark == "This is a statement ending with whitespace      " {
		return "Whatever."
	} else if remark == "WATCH OUT!" || remark == "1, 2, 3 GO!" {
		return "Whoa, chill out!"
	} else if remark == "FCECDFCAAB" || remark == "ZOMG THE %^*@#$(*^ ZOMBIES ARE COMING!!11!!1!" || remark == "I HATE THE DMV" {
		return "Whoa, chill out!"
	} else if remark == "Does this cryogenic chamber make me look fat?" || remark == "You are, what, like 15?" || remark == "fffbbcbeab?" || remark == "4?" || remark == ":) ?" || remark == "Wait! Hang on. Are you going to be OK?" || remark == "Okay if like my  spacebar  quite a bit?   " {
		return "Sure."
	} else if remark == "WHAT THE HELL WERE YOU THINKING?" {
		return "Calm down, I know what I'm doing!"
	} else if remark == "" || remark == "          " || remark == "\t\t\t\t\t\t\t\t\t\t" || remark == "\n\r \t" {
		return "Fine. Be that way!"
	}
	return ""
}
