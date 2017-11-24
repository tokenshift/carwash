package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"
)

const DefaultBufferSize = "64"

var (
	app = kingpin.New("carwash", "Sanitize sensitive values from an input stream.")

	debug = app.Flag("debug", "Turn on debug output (for carwash developers).").Short('d').Bool()

	keysFile = app.Flag("keys-file", "Path to a file containing keys to consider sensitive, one on each line.").Short('k').ExistingFile()
	valsFile = app.Flag("vals-file", "Path to a file containing values to consider sensitive, one on each line.").Short('v').ExistingFile()

	noEnv    = app.Flag("no-env", "Disable scanning of environment variables for sensitive values (enabled by default).").Short('E').Bool()
	noPredef = app.Flag("no-predef", "Don't add a predefined set of keys to consider sensitive (like PASSWORD, TOKEN, KEY...).").Short('D').Bool()

	bufferSize = app.Flag("buffer", "Set the number of lines to buffer in order to handle multiline sensitive values (e.g. private keys).").Short('b').Default(DefaultBufferSize).Int()
	mask       = app.Flag("mask", "Set the mask text to use to obscure sensitive values.").Short('m').Default("********").String()
)

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))

	debugLog("Initializing with buffer size", *bufferSize)
	scrubber := NewScrubber(*mask)

	if *keysFile == "" {
		debugLog("No keys file specified; skipping.")
	} else {
		debugLog("Loading keys file:", *keysFile)
		loadKeysFile(*keysFile, scrubber)
	}

	if *valsFile == "" {
		debugLog("No values file specified; skipping.")
	} else {
		debugLog("Loading values file:", *valsFile)
		loadValsFile(*valsFile, scrubber)
	}

	if *noEnv {
		debugLog("--no-env was specified, skipping")
	} else {
		debugLog("Scanning environment variables for sensitive values")
		// TODO
	}

	if *noPredef {
		debugLog("--no-predef was specified, skipping adding predefined keys")
	} else {
		debugLog("Adding predefined keys")
		addPredefinedKeys(scrubber)
	}

	debugLog("Beginning scanning of STDIN")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Println(scrubber.Scrub(scanner.Text()))
	}
	checkFatal(scanner.Err())
}

func debugLog(args ...interface{}) {
	if *debug {
		fmt.Fprintln(os.Stderr, args...)
	}
}

func checkFatal(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func loadKeysFile(filename string, scrubber Scrubber) {
	f, err := os.Open(filename)
	checkFatal(err)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		scrubber.AddSensitiveKey(scanner.Text())
	}

	checkFatal(scanner.Err())
}

func loadValsFile(filename string, scrubber Scrubber) {
	f, err := os.Open(filename)
	checkFatal(err)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		scrubber.AddSensitiveValue(scanner.Text())
	}

	checkFatal(scanner.Err())
}

func addPredefinedKeys(scrubber Scrubber) {
	// General Keys
	scrubber.AddSensitiveKey("KEY")
	scrubber.AddSensitiveKey("PASSWORD")
	scrubber.AddSensitiveKey("TOKEN")
	scrubber.AddSensitiveKey("SECRET")

	// AWS Keys
	scrubber.AddSensitiveKey("ACCESS_KEY_ID")
	scrubber.AddSensitiveKey("AccessKeyId")
	scrubber.AddSensitiveKey("SECRET_ACCESS_KEY")
	scrubber.AddSensitiveKey("SecretAccessKey")
	scrubber.AddSensitiveKey("SESSION_TOKEN")
	scrubber.AddSensitiveKey("SessionToken")
}
