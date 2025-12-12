package logs

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"testing"
)

var logger = log.New(os.Stderr, "", log.LstdFlags)

// Println mirrors TeaGo logs.Println.
func Println(args ...any) {
	_ = logger.Output(2, sprint(args...))
}

// Printf mirrors TeaGo logs.Printf.
func Printf(format string, args ...any) {
	_ = logger.Output(2, fmt.Sprintf(format, args...))
}

// Error mirrors TeaGo logs.Error.
func Error(args ...any) {
	_ = logger.Output(2, sprint(args...))
}

// Errorf mirrors TeaGo logs.Errorf.
func Errorf(format string, args ...any) {
	_ = logger.Output(2, fmt.Sprintf(format, args...))
}

// PrintAsJSON pretty prints an object for debugging/tests.
func PrintAsJSON(v any, extra ...any) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		Println(err)
		return
	}
	if len(extra) > 0 {
		if t, ok := extra[0].(*testing.T); ok {
			t.Log(string(data))
			return
		}
	}
	Println(string(data))
}

func sprint(args ...any) string {
	return fmt.Sprintln(args...)
}

// SetWriter allows redirecting logger output.
func SetWriter(w io.Writer) {
	if w == nil {
		return
	}
	logger.SetOutput(w)
}
