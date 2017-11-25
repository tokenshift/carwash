package main

import (
	"bufio"
	"bytes"
	"fmt"
	. "testing"
)

func BenchmarkManyReplacements(b *B) {
	var testData bytes.Buffer

	for i := 0; i < 1000; i += 1 {
		for _, testname := range exampleNames() {
			if contents, ok := readFile(testname, ".input"); ok {
				fmt.Fprintln(&testData, contents)
			}
		}
	}

	scrubber := NewScrubber("********")

	b.ResetTimer()

	for i := 0; i < b.N; i += 1 {
		scanner := bufio.NewScanner(&testData)
		for scanner.Scan() {
			scrubber.Scrub(scanner.Text())
		}
	}
}
