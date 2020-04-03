package logs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/kataras/golog"
	"github.com/rimxgo/config"
)

var Logger *golog.Logger

func init() {
	Logger = golog.New()
	Logger.SetLevel(config.GetString("logLevel"))
	Logger.Handle(simpleOutput)
	// Logger.Handle(jsonOutput)
}

func Any(v interface{}) string {
	bt, _ := json.Marshal(v)
	return string(bt)
}

// https://golang.org/doc/go1.9#callersframes
func getCaller() (string, int) {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])
	wd, _ := os.Getwd()

	for {
		frame, more := frames.Next()
		file := frame.File

		if !strings.Contains(file, "github.com/kataras/golog") {
			if relFile, err := filepath.Rel(wd, file); err == nil {
				file = "./" + relFile
			}
			return file, frame.Line
		}

		if !more {
			break
		}
	}

	return "???", 0
}

func simpleOutput(l *golog.Log) bool {
	prefix := golog.GetTextForLevel(l.Level, true)

	filename, line := getCaller()
	message := fmt.Sprintf(
		"%s %s [%s:%d] %s",
		prefix, l.FormatTime(), filename, line, l.Message,
	)

	if l.NewLine {
		message += "\n"
	}

	fmt.Print(message)
	return true
}

func jsonOutput(l *golog.Log) bool {
	fn, line := getCaller()

	var (
		datetime = l.FormatTime()
		level    = golog.GetTextForLevel(l.Level, false)
		message  = l.Message
		source   = fmt.Sprintf("%s#%d", fn, line)
	)

	jsonStr := fmt.Sprintf(
		"{\"datetime\":\"%s\", \"level\":\"%s\", \"message\":\"%s\", \"source\":\"%s\"}",
		datetime, level, message, source,
	)

	fmt.Println(jsonStr)
	return true
}
