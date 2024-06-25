package main

import (
	rx "regexp"
	"strings"

	sc "github.com/iancoleman/strcase"
)

func Replace(input string, query Query) *Replacement {
	fragments := getFragments(input, query)
	if fragments.IsEmpty() {
		return nil
	}
	output := getOutput(input, fragments)
	return &Replacement{
		Fragments: fragments,
		Input: input,
		Output: output,
	}
}

type Replacement struct {
	Fragments *Fragments
	Input string
	Output string
}

func NewReplacement(input string, fragments *Fragments, output string) *Replacement {
	return &Replacement{
		Fragments: fragments,
		Input: input,
		Output: output,
	}
}

func (r *Replacement) ReplacementOutput() string {
	return r.Output
}

func (r *Replacement) ReplacementInput() string {
	return r.Input
}

func (r *Replacement) NumFragments() int {
	return r.Fragments.Len()
}

func (r *Replacement) GetFragments() *Fragments {
	return r.Fragments
}

type Fragment struct {
	Index int
	Text  string
}

type Fragments struct {
	data [][2]Fragment
}

func NewFragments() *Fragments {
	return &Fragments{
		data: make([][2]Fragment, 0),
	}
}

func (f *Fragments) Len() int {
	return len(f.data)
}

func (f *Fragments) IsEmpty() bool {
	return len(f.data) == 0
}

func (f *Fragments) Add(inputIndex int, inputText string, outputIndex int, outputText string) {
	inputFragment := Fragment{Index: inputIndex, Text: inputText}
	outputFragment := Fragment{Index: outputIndex, Text: outputText}

	f.data = append(f.data, [2]Fragment{inputFragment, outputFragment})
}

type FragmentsIterator struct {
	fragments *Fragments
	index     int
}

func (f *Fragments) Iterator() *FragmentsIterator {
	return &FragmentsIterator{
		fragments: f,
		index:     0,
	}
}

func (ft *FragmentsIterator) Next() (*[2]Fragment, bool) {
	if ft.index > len(ft.fragments.data) {
		item := &ft.fragments.data[ft.index]
		ft.index++
		return item, true
	}

	return nil, false
}

type Replacer interface {
	Replace(buff string) *ReplacementResult
}

type ReplacementResult struct {
	Index      int
	InputText  string
	OutputText string
}

type SubStringReplacer struct {
	Pattern     string
	Replacement string
}

func NewSubStringReplacer(pattern, replacement string) *SubStringReplacer {
	return &SubStringReplacer{
		Pattern: pattern,
		Replacement: replacement,
	}
} 

func (sr *SubStringReplacer) Replace(buff string) *ReplacementResult {
	index := strings.Index(buff, sr.Pattern)
	if index == -1 {
		return nil
	}

	inputText := buff[index : index+len(sr.Pattern)]
	outputText := sr.Replacement
	return &ReplacementResult{
		Index:      index,
		InputText:  inputText,
		OutputText: outputText,
	}
}

type SubvertReplacer struct {
	Items [][2]string
}

func NewSubvertReplacer(items [][2]string) *SubvertReplacer {
	return &SubvertReplacer{
		Items: items,
	}
}

func (sr *SubvertReplacer) Replace(buff string) *ReplacementResult {
	bestIndex := len(buff)
	bestPattern := -1

	for i, item := range sr.Items {
		pattern := item[0]
		if index := strings.Index(buff, pattern); index != -1 {
			if index < bestIndex {
				bestIndex = index
				bestPattern = i
			}
		}
	}

	if bestPattern == -1 {
		return nil
	}

	bestItem := sr.Items[bestPattern]
	pattern := bestItem[0]
	replacement := bestItem[1]

	return &ReplacementResult{
		Index:      bestIndex,
		InputText:  pattern,
		OutputText: replacement,
	}
}

type RegexReplacer struct {
	Regex *rx.Regexp
	Replacement string
} 

func NewRegexReplacer(regex *rx.Regexp, replacement string) *RegexReplacer {
	return &RegexReplacer{
       Regex: regex,
	   Replacement: replacement,
	}
}

func (rr *RegexReplacer) Replace(buff string) *ReplacementResult {
	regexMatch := rr.Regex.FindStringIndex(buff)
	if regexMatch == nil {
		return nil
	}

	start, end := regexMatch[0], regexMatch[1]
	inputText := buff[start:end]
	outputText := rr.Regex.ReplaceAllString(inputText, rr.Replacement)
	return &ReplacementResult{
		Index: start,
		InputText: inputText,
		OutputText: outputText,
	}
}


func getFragments(input string, query Query) *Fragments {
	switch q := query.(type) {
	case SubString:
		finder := NewSubStringReplacer(q.Old, q.New)
		return getFragmentsWithFinder(input, finder)
	case Regex:
		finder := NewRegexReplacer(q.Pattern, q.Replacement)
		return getFragmentsWithFinder(input, finder)
	case Subvert:
		functions := []func(string) string {
		  sc.ToScreamingSnake,
          sc.ToCamel,
		  sc.ToKebab,
		  sc.ToSnake,
		  sc.ToLowerCamel,
		  sc.ToScreamingKebab,
		}

		items := make([][2]string, len(functions))
		for i, function := range functions {
			items[i] = [2]string{function(q.Pattern), function(q.Replacement)}
		}

		finder := NewSubvertReplacer(items)
		return getFragmentsWithFinder(input, finder)
	default:
		return NewFragments()
	}
}

func getFragmentsWithFinder(input string, finder Replacer) *Fragments {
	fragments := NewFragments()
	inputIndex := 0
	outputIndex := 0

	for {
		res := finder.Replace(input[inputIndex:])
		if res == nil {
			break
		}
		index := res.Index
		inputText := res.InputText
		outputText := res.OutputText

		inputIndex += index
		outputIndex += index

		fragments.Add(inputIndex, inputText, outputIndex, outputText)

		inputIndex += len(inputText)
		outputIndex += len(outputText)
	}

	return fragments
}

func getOutput(input string, fragments *Fragments) string {
	var output strings.Builder
	currentIndex := 0;
    
	for _, fragmentPair := range fragments.data {
		inputFragment := fragmentPair[0]
		outputFragment := fragmentPair[1]

		inputText := inputFragment.Text
		inputIndex := inputFragment.Index
		outputText := outputFragment.Text

		output.WriteString(input[currentIndex:inputIndex])
		output.WriteString(outputText)
		currentIndex = inputIndex + len(inputText)
	}
	
	output.WriteString(input[currentIndex:])
	return output.String()
}