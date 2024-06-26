package main

import (
    "fmt"
	"regexp"
    
	"github.com/fatih/color"
)

func main() {
	substringQuery := NewSubString("hello", "hi")
	regexQuery := NewRegex(regexp.MustCompile(`old`), "new")
	subvertQuery := NewSubvert("case", "casecase")

	input := "This is a dummy text [hello] [old] []"

	fmt.Println("Substring Query:")
	replacement := Replace(input, substringQuery)
	if replacement != nil {
		c := color.New(color.FgHiGreen).Add(color.Bold)
		c.Printf("Output: %s\n", replacement.ReplacementOutput())
	} else {
		c := color.New(color.FgRed).Add(color.Bold)
		c.Println("No replacement made.")
	}

	fmt.Println("\nRegex Query:")
	replacement = Replace(input, regexQuery)
	if replacement != nil {
		c := color.New(color.FgHiMagenta).Add(color.Bold)
		c.Printf("Output: %s\n", replacement.ReplacementOutput())
	} else {
		c := color.New(color.FgRed).Add(color.Bold)
		c.Println("No replacement made.")
	}

	fmt.Println("\nSubvert Query:")
	replacement = Replace(input, subvertQuery)
	if replacement != nil {
		c := color.New(color.FgHiYellow).Add(color.Bold)
		c.Printf("Output: %s\n", replacement.ReplacementOutput())
	} else {
		c := color.New(color.FgRed).Add(color.Bold)
		c.Println("No replacement made.")
	}

	// console := NewConsoleWithVerbosity(Normal)

	// fragments := &Fragments{
	// 	data: [][2]Fragment{
	// 		{Fragment{20, "colored"}, Fragment{20, "colourful"}},
	// 		{Fragment{28, "text"}, Fragment{30, "words"}},
	// 	},
	// }

	// replacement1 := NewReplacement("This is a line with colored text fragments.", fragments, "This is a line with colourful words fragments.")

	// console.PrintReplacement("", replacement1)

}