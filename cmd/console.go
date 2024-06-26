package main

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

type Verbosity int

const (
	Quiet Verbosity = iota
	Normal
)

type Console struct {
	verbosity Verbosity
}

func NewConsoleWithVerbosity(verbosity Verbosity) *Console {
	return &Console{
		verbosity: verbosity,
	}
}

func NewConsole() *Console {
	return &Console{
		verbosity: Normal,
	}
}

func (c *Console) PrintMessage(message string) {
	if c.verbosity == Quiet {
		return
	}
	fmt.Print(message)
}

func (c *Console) PrintError(errorMessage string) {
	fmt.Print(errorMessage)
}

func (c *Console) PrintReplacement(prefix string, replacement *Replacement) {
	redUnderline := color.New(color.FgRed, color.Underline).SprintFunc()
	greenUnderline := color.New(color.FgGreen, color.Underline).SprintFunc()

	fragments := replacement.GetFragments().data

	inputFragments := make([]Fragment, len(fragments))
	outputFragments := make([]Fragment, len(fragments))

	for i, frag := range fragments {
		inputFragments[i] = frag[0]
		outputFragments[i] = frag[1]
	}

	redPrefix := fmt.Sprintf("%s%s", prefix, color.New(color.FgRed).Sprint(" - "))
	c.printFragments(redPrefix, redUnderline, replacement.ReplacementInput(), inputFragments)

	greenPrefix := fmt.Sprintf("%s%s", prefix, color.New(color.FgGreen).Sprint(" + "))
	c.printFragments(greenPrefix, greenUnderline, replacement.ReplacementOutput(), outputFragments)
}

func (c *Console) printFragments(prefix string, colorFunc func(a ...interface{}) string, line string, fragments []Fragment) {
	c.PrintMessage(prefix)
	currentIndex := 0
	for i, fragment := range fragments {
		index := fragment.Index
		text := fragment.Text

		if index < 0 || index > len(line) {
			continue
		}

		if i == 0 {
			c.PrintMessage(strings.TrimLeft(line[currentIndex:index], " "))
		} else {
			c.PrintMessage(line[currentIndex:index])
		}

		endIndex := index + len(text)
		if endIndex > len(line) {
			endIndex = len(line)
		}

		c.PrintMessage(colorFunc(text))
		currentIndex = endIndex
	}

	if currentIndex < len(line) {
		c.PrintMessage(line[currentIndex:])
	}

	if !strings.HasSuffix(line, "\n") {
		c.PrintMessage("\n")
	}
}
