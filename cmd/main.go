package main

import (
    "fmt"
	"regexp"

)

func main() {
	substringQuery := NewSubString("hello", "hii")
	regexQuery := NewRegex(regexp.MustCompile(`(tridip|hello)`), "fuck")
	subvertQuery := NewSubvert("foo_bar", "spam_eggs")

	input := "hello my name is foo_bar and my freind name is FooBar and his friend name is foo-bar"

	fmt.Println("Substring Query:")
	replacement := Replace(input, substringQuery)
	if replacement != nil {
		fmt.Println("Output:", replacement.ReplacementOutput())
	} else {
		fmt.Println("No replacement made.")
	}

	fmt.Println("\nRegex Query:")
	replacement = Replace(input, regexQuery)
	if replacement != nil {
		fmt.Println("Output:", replacement.ReplacementOutput())
	} else {
		fmt.Println("No replacement made.")
	}

	fmt.Println("\nSubvert Query:")
	replacement = Replace(input, subvertQuery)
	if replacement != nil {
		fmt.Println("Output:", replacement.ReplacementOutput())
	} else {
		fmt.Println("No replacement made.")
	}
}