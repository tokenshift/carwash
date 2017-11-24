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
	files, err := filepath.Glob("./tests/*.input")
	checkFatal(err)

	for _, filename := range files {
		testname := strings.TrimSuffix(filename, filepath.Ext(filename))
		testSingleExample(t, testname)
	}
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

	if _, err := os.Stat(fmt.Sprintf("%s.vals", name)); err == nil {
		loadValsFile(fmt.Sprintf("%s.vals", name), scrubber)
	} else if !os.IsNotExist(err) {
		t.Error(err)
		return
	}

	inScanner := bufio.NewScanner(input)
	outScanner := bufio.NewScanner(expected)

	for inScanner.Scan() {
		if !outScanner.Scan() {
			t.Errorf("%s.expected ran out of lines before %s.input", name, name)
			return
		}

		scrubbed := scrubber.Scrub(inScanner.Text())
		if scrubbed != outScanner.Text() {
			t.Errorf("%s\nExpected: %s\nActual:   %s", name, outScanner.Text(), scrubbed)
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
