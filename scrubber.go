package main

import (
	"bytes"
	"fmt"
	"strings"
)

type Scrubber interface {
	AddSensitiveKey(string)
	AddSensitiveValue(string)
	Discover(string)
	Obscure(string) string
	Scrub(string) string
}

func NewScrubber(mask string) Scrubber {
	return &scrubber{
		scanners: make([]Scanner, 0, 64),
		keys:     NewStringSet(),
		vals:     NewSensitiveValues(),
		mask:     mask,
	}
}

type scrubber struct {
	scanners []Scanner
	keys     StringSet
	vals     SensitiveValues
	mask     string
}

func (s *scrubber) AddSensitiveKey(key string) {
	lower := strings.ToLower(key)
	if s.keys.Add(lower) {
		s.scanners = append(s.scanners, NewKeyScanner(lower))
	}
}

func (s *scrubber) AddSensitiveValue(val string) {
	s.vals.Add(val)
}

func (s *scrubber) Scrub(input string) string {
	s.Discover(input)
	return s.Obscure(input)
}

func (s *scrubber) Discover(input string) {
	for _, scanner := range s.scanners {
		scanner.Each(func(matched string) {
			s.AddSensitiveValue(matched)
		}, input)
	}
}

func (s *scrubber) Obscure(input string) string {
	pattern := s.vals.Pattern()
	if pattern == nil {
		return input
	}

	var result bytes.Buffer

	index := 0
	for index < len(input) {
		next := pattern.FindStringIndex(input[index:])
		if next == nil {
			fmt.Fprint(&result, input[index:])
			return result.String()
		} else {
			nextStart, nextEnd := next[0], next[1]
			fmt.Fprint(&result, input[index:index+nextStart])
			fmt.Fprint(&result, s.mask)
			index = index + nextEnd
		}
	}

	return result.String()
}

// TODO: Turn the list of sensitive values into a regex, so this doesn't require
// a bunch of iterations over the input string.
func maskLine(input, match, mask string) string {
	inputLower := strings.ToLower(input)
	matchLower := strings.ToLower(match)

	var result bytes.Buffer

	index := 0
	for index < len(inputLower) {
		next := strings.Index(inputLower[index:], matchLower)
		if next < 0 {
			fmt.Fprint(&result, input[index:])
			return result.String()
		} else {
			fmt.Fprint(&result, input[index:index+next])
			fmt.Fprint(&result, mask)
			index = index + next + len(match)
		}
	}

	return result.String()
}
