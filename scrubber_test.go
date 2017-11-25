package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	. "testing"
)

func TestAllExamples(t *T) {
	for _, testname := range exampleNames() {
		testSingleExample(t, testname)
	}
}

func exampleNames() []string {
	files, err := filepath.Glob("./tests/*.input")
	checkFatal(err)

	names := make([]string, 0, len(files))
	for _, filename := range files {
		testname := strings.TrimSuffix(filename, filepath.Ext(filename))
		names = append(names, testname)
	}

	return names
}

func testSingleExample(t *T, name string) {
	input, err := os.Open(fmt.Sprintf("%s.input", name))
	if err != nil {
		t.Error(err)
		return
	}

	defer input.Close()

	expected, err := os.Open(fmt.Sprintf("%s.expected", name))
	if err != nil {
		t.Error(err)
		return
	}

	defer expected.Close()

	scrubber := NewScrubber("********")
	addPredefinedKeys(scrubber)

	if _, err := os.Stat(fmt.Sprintf("%s.vals", name)); err == nil {
		loadValsFile(fmt.Sprintf("%s.vals", name), scrubber)
	} else if !os.IsNotExist(err) {
		t.Error(err)
		return
	}

	inScanner := bufio.NewScanner(input)
	outScanner := bufio.NewScanner(expected)
	line := 0

	for inScanner.Scan() {
		if !outScanner.Scan() {
			t.Errorf("%s.expected ran out of lines before %s.input", name, name)
			return
		}

		line += 1

		scrubbed := scrubber.Scrub(inScanner.Text())
		if scrubbed != outScanner.Text() {
			t.Errorf("%s:%d\nExpected: %s\nActual:   %s",
				name, line,
				outScanner.Text(),
				scrubbed)
		}
	}

	if inScanner.Err() != nil {
		t.Error(inScanner.Err())
	}

	if outScanner.Err() != nil {
		t.Error(outScanner.Err())
	}
}

func readFile(name, extension string) (string, bool) {
	filename := fmt.Sprintf("%s%s", name, extension)
	if f, err := ioutil.ReadFile(filename); err != nil {
		return "", false
	} else {
		return string(f), true
	}
}
