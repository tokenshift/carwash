package main

import (
// "fmt"
// "regexp"
)

type Scanner interface {
	Each(func(string), string)
}

// type keyScanner struct {
// 	key      string
// 	patterns []*regexp.Regexp
// }

// // Looks for the following patterns of key-value pairs:
// // * KEY=value
// // * KEY value
// // * KEY: value
// // * KEY => value
// // Where the value may be single or double quoted.
// func NewKeyScanner(key string) Scanner {
// 	escapedKey := regexp.QuoteMeta(key)
// 	patterns := []*regexp.Regexp{
// 		regexp.MustCompile(fmt.Sprintf(`(?i)%s['"]?\s*(?:=>|=|:|\s+)\s*'((?:\\'|[^'])+)'`, escapedKey)),
// 		regexp.MustCompile(fmt.Sprintf(`(?i)%s['"]?\s*(?:=>|=|:|\s+)\s*"((?:\\"|[^"])+)"`, escapedKey)),
// 		regexp.MustCompile(fmt.Sprintf(`(?i)%s['"]?\s*(?:=>|=|:|\s+)\s*((?:\\\s|\\"|\\'|[^\s'"])+)`, escapedKey)),
// 	}

// 	return keyScanner{
// 		key:      key,
// 		patterns: patterns,
// 	}
// }

// func (ks keyScanner) Each(f func(string), input string) {
// }

// func (ks keyScanner) unescape(val string) string {
// }
