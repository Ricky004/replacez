package main

import (
	rx "regexp"
)

type Query interface {
	isQuery()
}

type QueryType int

type SubString struct {
	Old string
	New string
}

func (s SubString) isQuery() {}

type Regex struct {
	Pattern     *rx.Regexp
	Replacement string
}

func (r Regex) isQuery() {}

type Subvert struct {
	Pattern     string
	Replacement string
}

func (s Subvert) isQuery() {}

func NewSubString(old, new string) Query {
	return SubString{Old: old, New: new}
}

func NewRegex(pattern *rx.Regexp, replacement string) Query {
	return Regex{Pattern: pattern, Replacement: replacement}
}

func NewSubvert(pattern, replacement string) Query {
	return Subvert{Pattern: pattern, Replacement: replacement}
}
