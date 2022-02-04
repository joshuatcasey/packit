package scribe

import (
	"fmt"
	"io"
	"strings"
)

// A Logger provides a standard logging interface for doing basic low level
// logging tasks as well as debug logging.
type Logger struct {
	writer io.Writer
	LeveledLogger
	Debug LeveledLogger
}

// NewLogger takes a writer and returns a Logger that writes to the given
// writer. The default writter sends all debug logging to io.Discard.
func NewLogger(writer io.Writer) Logger {
	return Logger{
		writer:        writer,
		LeveledLogger: NewLeveledLogger(writer),
		Debug:         NewLeveledLogger(io.Discard),
	}
}

// WithLevel takes in a log level string and configures the log level of the
// logger. To enable debug logging the log level must be set to "DEBUG".
func (l Logger) WithLevel(level string) Logger {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return Logger{
			writer:        l.writer,
			LeveledLogger: NewLeveledLogger(l.writer),
			Debug:         NewLeveledLogger(l.writer),
		}
	default:
		return Logger{
			writer:        l.writer,
			LeveledLogger: NewLeveledLogger(l.writer),
			Debug:         NewLeveledLogger(io.Discard),
		}
	}
}

// A LeveledLogger provides a standard interface for basic formatted logging.
type LeveledLogger struct {
	Title      FuncWriter
	Process    FuncWriter
	Subprocess FuncWriter
	Action     FuncWriter
	Detail     FuncWriter
	Subdetail  FuncWriter
}

// NewLeveledLogger takes a writer and returns a LeveledLogger that writes to the given
// writer.
func NewLeveledLogger(writer io.Writer) LeveledLogger {
	return LeveledLogger{
		Title:      NewFuncWriter(NewWriter(writer)),
		Process:    NewFuncWriter(NewWriter(writer, WithIndent(1))),
		Subprocess: NewFuncWriter(NewWriter(writer, WithIndent(2))),
		Action:     NewFuncWriter(NewWriter(writer, WithIndent(3))),
		Detail:     NewFuncWriter(NewWriter(writer, WithIndent(4))),
		Subdetail:  NewFuncWriter(NewWriter(writer, WithIndent(5))),
	}
}

// Break inserts a line break in the log output
func (l LeveledLogger) Break() {
	l.Title("")
}

type _fw struct{}

type FuncWriter func(format string, v ...interface{}) io.Writer

func NewFuncWriter(writer io.Writer) FuncWriter {
	return func(format string, v ...interface{}) io.Writer {
		skip := len(v) > 0 && v[0] == _fw{}
		if !skip {
			if !strings.HasSuffix(format, "\n") {
				format = format + "\n"
			}
			fmt.Fprintf(writer, format, v...)
		}

		return writer
	}
}

func (fw FuncWriter) Write(b []byte) (int, error) {
	return fw("", _fw{}).Write(b)
}
