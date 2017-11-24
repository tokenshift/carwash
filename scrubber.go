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
		scanners: make([]Scanner, 0),
		vals:     NewStringSet(),
		mask:     mask,
	}
}

type scrubber struct {
	scanners []Scanner
	vals     StringSet
	mask     string
}

func (s *scrubber) AddSensitiveKey(key string) {
	// s.scanners = append(s.scanners, NewKeyScanner(key))
}

func (s *scrubber) AddSensitiveValue(val string) {
	s.vals.Add(strings.ToLower(val))
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
	s.vals.Each(func(val string) {
		input = maskLine(input, val, s.mask)
	})

	return input
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
