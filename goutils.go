// some utility functions in Golang
package goutils

import (
	"crypto/rand"
	"fmt"
	"github.com/lao-tseu-is-alive/golog"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
)

func PrintCallStack() {
	// Ask runtime.Callers for up to 10 pcs, including runtime.Callers itself.
	pc := make([]uintptr, 10)
	n := runtime.Callers(2, pc) // skip runtime.Callers and PrintCallStack levels
	if n == 0 {
		// No pcs available. Stop now.
		// This can happen if the first argument to runtime.Callers is large.
		return
	}
	pc = pc[:n] // pass only valid pcs to runtime.CallersFrames
	frames := runtime.CallersFrames(pc)

	level := 0
	// Loop to get frames.
	for {
		frame, more := frames.Next()
		if strings.Contains(frame.File, "runtime/") {
			break
		}
		fmt.Printf("-level:%v | %s | %s:%d\n", level, frame.Function, frame.File, frame.Line)
		level += 1
		if !more {
			break
		}
	}
}

func GetKeyValue(s, sep string) (string, string) {
	res := strings.SplitN(s, sep, 2)
	return res[0], res[1]
}

func GetEnvOrDefault(key, defVal string) string {
	val, exist := os.LookupEnv(key)
	if !exist {
		return defVal
	}
	return val
}

// Generates a cryptographically secure random 16 bytes UUID using crypto/rand package
// returns a string of lenght 36 like this : bcaf4890-8b63-423b-258f-7a11004a8bf0
// On Linux and FreeBSD, it uses getrandom(2) if available, /dev/urandom otherwise.
// https://golang.org/pkg/crypto/rand/
// For a more classic RFC 4122 v4 GUID you can use https://github.com/satori/go.uuid
// more info at https://blog.kowalczyk.info/article/JyRZ/generating-good-unique-ids-in-go.html
func GetUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

//  loops over a slice, and selects and returns only the elements
//  that meet a criteria captured by a function value
//  https://stackoverflow.com/questions/37562873/most-idiomatic-way-to-select-elements-from-an-array-in-golang
func filter(ss []string, test func(string) bool) (ret []string) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

// get back the text content of a filename if FileEncoding is used decode the content with this encoding (default is
func GetFileTextContent(filename string, FileEncoding string) string {
	golog.Un(golog.Trace("getFileTextContent(%s, %s)", filename, FileEncoding))
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	// Read all in raw form.
	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	content := string(b)

	if FileEncoding == "" {
		//no special encoding asked so send content as is
		return content
	} else {
		var decoder *encoding.Decoder
		switch strings.ToLower(FileEncoding) {
		case "win1252", "cp1252", "windows-1252", "windows1252":
			// Decode CP1252 to unicode https://en.wikipedia.org/wiki/Windows-1252
			decoder = charmap.Windows1252.NewDecoder()
		default:
			decoder = charmap.ISO8859_1.NewDecoder()
		}
		reader := decoder.Reader(strings.NewReader(content))
		b, err = ioutil.ReadAll(reader)
		if err != nil {
			panic(err)
		}
		return string(b)
	}
}
