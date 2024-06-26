package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type FilePatcher struct {
	Path            string
	NewContents     string
	NumReplacements int
	NumLines        int
}

func NewFilePatcher(console *Console, path string, query Query) (*FilePatcher, error) {
	var numReplacements, numLines int

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("not open %s: %w", path, err)
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	var newContents strings.Builder

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		replacement := Replace(line, query)
		if replacement == nil {
			newContents.WriteString(line)
			newContents.WriteString("\n")
		} else {
			numLines++
			numReplacements += replacement.NumFragments()
			lineNo := numLines
			prefix := fmt.Sprintf("%s:%d", path, lineNo)
			console.PrintReplacement(prefix, replacement)
			newContents.WriteString(replacement.ReplacementOutput())
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error while reading %s: %w", path, err)
	}

	if numReplacements == 0 {
		return nil, nil
	}

	return &FilePatcher{
		Path:            path,
		NewContents:     newContents.String(),
		NumLines:        numLines,
		NumReplacements: numReplacements,
	}, nil
}

type LineIterator struct {
	delimeter byte
	reader *bufio.Reader
}

func NewLineIterator(delimeter byte, reader io.Reader) *LineIterator {
	return &LineIterator{
		delimeter: delimeter,
		reader: bufio.NewReader(reader),
	}
}

func (lt *LineIterator) Next() ([]byte, error) {
	line, err := lt.reader.ReadBytes(lt.delimeter)
	if err != nil {
		if err == io.EOF && len(line) > 0 {
			return line, nil
		}
		return nil, err
	}
	return line, nil
}