package zprompt

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Prompt
// Show a text message and wait for a response
func Prompt(text string) (res string) {
	r := bufio.NewReader(os.Stdin)
	fmt.Fprint(os.Stderr, text+" ")
	res, _ = r.ReadString('\n')
	return strings.TrimSpace(res)
}

// Confirm
// Ask a Y/N question
func Confirm(text string, defaultYes bool) (answeredYes bool) {
	if defaultYes {
		text += " (Y/n)"
	} else {
		text += " (y/N)"
	}
	res := Prompt(text)
	res = strings.ToLower(res)
	return res == "y" || res == "yes"
}
