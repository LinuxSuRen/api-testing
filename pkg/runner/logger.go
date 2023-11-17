/*
MIT License

Copyright (c) 2023 API Testing Authors.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package runner

import (
	"fmt"
	"io"
)

// LevelWriter represents a writer with level
type LevelWriter interface {
	Info(format string, a ...any)
	Debug(format string, a ...any)
	Trace(format string, a ...any)
}

// FormatPrinter represents a formart printer with level
type FormatPrinter interface {
	Fprintf(w io.Writer, level, format string, a ...any) (n int, err error)
}

type defaultLevelWriter struct {
	level LogLevel
	io.Writer
	FormatPrinter
}

// NewDefaultLevelWriter creates a default LevelWriter instance
func NewDefaultLevelWriter(level string, writer io.Writer) LevelWriter {
	result := &defaultLevelWriter{
		Writer: writer,
	}
	switch level {
	case "trace":
		result.level = LogLevelTrace
	case "debug":
		result.level = LogLevelDebug
	case "info":
		result.level = LogLevelInfo
	}
	return result
}

type LogLevel int

const (
	LogLevelInfo  LogLevel = 3
	LogLevelDebug LogLevel = 5
	LogLevelTrace LogLevel = 7
)

// Fprintf implements interface FormatPrinter
func (w *defaultLevelWriter) Fprintf(writer io.Writer, level LogLevel, format string, a ...any) (n int, err error) {
	if level <= w.level {
		return fmt.Fprintf(writer, format, a...)
	}
	return
}

// Info writes the info level message
func (w *defaultLevelWriter) Info(format string, a ...any) {
	w.Fprintf(w.Writer, LogLevelInfo, format, a...)
}

// Debug writes the debug level message
func (w *defaultLevelWriter) Debug(format string, a ...any) {
	w.Fprintf(w.Writer, LogLevelDebug, format, a...)
}

func (w *defaultLevelWriter) Trace(format string, a ...any) {
	w.Fprintf(w.Writer, LogLevelTrace, format, a...)
}
