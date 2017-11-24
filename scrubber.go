package main

import (
	"bytes"
	"fmt"
	"sort"
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
		vals:     make([]string, 0, 1024),
		mask:     mask,
	}
}

type scrubber struct {
	scanners []Scanner
	keys     StringSet
	vals     []string
	mask     string
}

func (s *scrubber) AddSensitiveKey(key string) {
	lower := strings.ToLower(key)
	if s.keys.Add(lower) {
		s.scanners = append(s.scanners, NewKeyScanner(lower))
	}
}

func (s *scrubber) AddSensitiveValue(val string) {
	s.vals = append(s.vals, strings.ToLower(val))
	// Strings are sorted by length so that the longest matches are obscured first.
	sort.Sort(ByLongest(s.vals))
	// TODO: Move this into a value matcher class; make sure strings aren't added redundantly.
}

type ByLongest []string

func (ss ByLongest) Less(i, j int) bool {
	return len(ss[i]) > len(ss[j])
}

func (ss ByLongest) Len() int {
	return len(ss)
}

func (ss ByLongest) Swap(i, j int) {
	ss[i], ss[j] = ss[j], ss[i]
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
	for _, val := range s.vals {
		input = maskLine(input, val, s.mask)
	}

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
