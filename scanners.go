package main

import (
	"fmt"
	"regexp"
	"strings"
)

type Scanner interface {
	Each(func(string), string)
}

type keyScanner struct {
	key     string
	pattern *regexp.Regexp
}

// Looks for the following patterns of key-value pairs:
// * KEY=value
// * KEY value
// * KEY: value
// * KEY => value
// Where the key or value may be single or double quoted.
func NewKeyScanner(key string) Scanner {
	lowerKey := strings.ToLower(key)
	escapedKey := regexp.QuoteMeta(lowerKey)

	pattern := regexp.MustCompile(fmt.Sprintf(`(?i)%s['"]?\s*(?:=>|=|:|\s+)\s*("(?:\\"|[^"])+"|'(?:\\'|[^'])+'|(?:\\\s|\\"|\\'|[^\s'"])+)`, escapedKey))

	return keyScanner{
		key:     lowerKey,
		pattern: pattern,
	}
}

func (ks keyScanner) Each(f func(string), input string) {
	matches := ks.pattern.FindAllStringSubmatch(input, -1)
	for _, match := range matches {
		value := stripQuotes(match[1])
		f(value)

		unescaped := unescape(value)
		if unescaped != value {
			f(unescaped)
		}
	}
}

func stripQuotes(val string) string {
	if strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"") {
		return val[1 : len(val)-1]
	} else if strings.HasPrefix(val, "'") && strings.HasSuffix(val, "'") {
		return val[1 : len(val)-1]
	} else {
		return val
	}
}

var rxEscapeChar = regexp.MustCompile(`\\(.)`)

func unescape(val string) string {
	return rxEscapeChar.ReplaceAllStringFunc(val, func(match string) string {
		switch match {
		case "\\0":
			return string([]byte{0})
		case "\\a":
			return "\a"
		case "\\b":
			return "\b"
		case "\\f":
			return "\f"
		case "\\n":
			return "\n"
		case "\\r":
			return "\r"
		case "\\t":
			return "\t"
		case "\\v":
			return "\v"
		default:
			return match[1:]
		}
	})
}
