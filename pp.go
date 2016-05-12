package main

import (
	"os"
	"bufio"
	"strings"
	"strconv"
	"encoding/json"
	"C"
)

func main() {
}


//export ProcessInRuby
func ProcessInRuby(intOpts *C.char) *C.char {

	type RubyOptions struct {
		File             string `json:"file"`
		Out              string `json:"out"`
		IgnoreReferenceN bool `json:"ignore_reference_n"`
		MinDepth         int `json:"min_depth"`
		MinNonRefCount   int `json:"min_non_ref_count"`
	}

	optString := C.GoString(intOpts)
	println("GO RECEIVED: " + optString)

	ro := RubyOptions{}
	json.Unmarshal([]byte(optString), &ro)
	options := Options{ro.MinDepth, ro.MinNonRefCount, ro.IgnoreReferenceN}

	if _, err := os.Stat(ro.File); os.IsNotExist(err) {
		// path/to/whatever does not exist
		panic(ro.File + " DOES NOT EXIST")
	}

	inFile, _ := os.Open(ro.File)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		text := scanner.Text()
		s := strings.Split(text, "\t");

		if (len(s) != 6) {
			panic("there were " + strconv.Itoa(len(s)) + " chunks instead of 6")
		}

		if (isSNP(s, options)) {
			str := s[0] + "\t" + s[1] + "\t" + s[2] + "\t" + s[3] + "\t" + s[4] + "\t" + s[5] + "\n"
			writeLine(ro.Out, str)
		}

	}
	return C.CString(ro.Out)
}

func isSNP(p []string, options Options) bool {
	if (p[4] == "*") {
		return false
	}
	if (options.ignoreReferenceN) {
		if (p[2] == "N" || p[2] == "n") {
			return false
		}
	}

	i, err := strconv.Atoi(p[3])
	if err != nil {
		os.Exit(2)
	}

	if (i >= options.minDepth && nonRefCount(p[4]) >= options.minNonRefCount) {
		return true
	}
	return false

}

type Options struct {
	minDepth         int
	minNonRefCount   int
	ignoreReferenceN bool
}

func nonRefCount(str string) int {
	return strings.Count(str, "A") + strings.Count(str, "T") + strings.Count(str, "G") + strings.Count(str, "C") + strings.Count(str, "a") + strings.Count(str, "t") + strings.Count(str, "g") + strings.Count(str, "c")
}

func writeLine(filename string, text string) {
	f, err := os.OpenFile(filename, os.O_APPEND | os.O_WRONLY | os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(text); err != nil {
		panic(err)
	}
}