// some utility functions in Golang
package goutils

import (
	"crypto/rand"
	"fmt"
	"github.com/mgutz/ansi"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

func GetTimeStamp() string {
	t := time.Now()
	timeString := fmt.Sprintf("%d/%02d/%02d %02d:%02d:%02d.%03d",
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
		t.Nanosecond()/100000)
	return fmt.Sprintf("%s", timeString)
}

func GetCaller(skip int) (file string, line int, function string) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(skip, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return frame.File, frame.Line, frame.Function
}

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

func DoItOrDie(err error, message string, v ...interface{}) {
	if err != nil {
		var loggerError = log.New(os.Stderr, "FATAL_ERROR: ", log.Lshortfile)
		filename, line, funcname := GetCaller(3)
		red := ansi.ColorFunc("red+b:yellow+h")
		logErr := loggerError.Output(2,
			red(fmt.Sprintf(
				"[%s], Function: %s, File: %s:%d, Error: [%v], Message: %s",
				GetTimeStamp(), funcname, filename, line, err, fmt.Sprintf(message, v...))))
		if logErr != nil {
			log.Fatalln("ERROR in DoItOrDie trying to output Err(message) to stderr console !")
		}
		log.Fatalf("# FATAL ERROR %s in function %s, [%v]", fmt.Sprintf(message, v...), funcname, err)
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
