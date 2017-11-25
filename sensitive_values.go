package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

type SensitiveValues struct {
	values  []string
	exists  StringSet
	pattern *regexp.Regexp
}

func NewSensitiveValues() SensitiveValues {
	return SensitiveValues{
		values:  make([]string, 0, 1024),
		exists:  NewStringSet(),
		pattern: nil,
	}
}

func (sv *SensitiveValues) Add(val string) bool {
	val = strings.ToLower(strings.TrimSpace(val))
	if val == "" {
		return false
	} else if sv.exists.Add(val) {
		sv.values = append(sv.values, val)
		sort.Sort(ByLongest(sv.values))
		sv.recreatePattern()
		return true
	} else {
		return false
	}
}

func (sv SensitiveValues) Pattern() *regexp.Regexp {
	return sv.pattern
}

func (sv *SensitiveValues) recreatePattern() {
	components := make([]string, len(sv.values))
	for i, val := range sv.values {
		components[i] = regexp.QuoteMeta(val)
	}

	pattern := fmt.Sprintf("(?i)%s", strings.Join(components, "|"))

	sv.pattern = regexp.MustCompile(pattern)
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
