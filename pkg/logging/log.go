/*
Copyright 2024 API Testing Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package logging

import (
	"io"
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogLevel defines a log level for api-testing logs.
type LogLevel string

// Log level const.
const (
	// LogLevelDebug defines the "debug" logging level.
	LogLevelDebug LogLevel = "debug"

	// LogLevelInfo defines the "Info" logging level.
	LogLevelInfo LogLevel = "info"

	// LogLevelWarn defines the "Warn" logging level.
	LogLevelWarn LogLevel = "warn"

	// LogLevelError defines the "Error" logging level.
	LogLevelError LogLevel = "error"
)

// apiTestingLog defines a make up part that supports a configured logging level.
type apiTestingLogComponent string

const (
	// LogComponentApiTestingDefault defines the "default"-wide logging component. When specified,
	// all other logging components are ignored.
	LogComponentApiTestingDefault apiTestingLogComponent = "default"

	// LogComponentApiTestingTesting represents the logging component for testing.
	LogComponentApiTestingTesting apiTestingLogComponent = "testing"
)

// ApiTestingLogging defines logging for api-testing.
type ApiTestingLogging struct {
	// Level is the logging level. If unspecified, defaults to "info".
	Level map[apiTestingLogComponent]LogLevel
}

// Logger represents a logger.
type Logger struct {
	// Embedded Logger interface
	logr.Logger
	logging       *ApiTestingLogging
	sugaredLogger *zap.SugaredLogger
}

func NewLogger(logging *ApiTestingLogging) Logger {
	logger := initZapLogger(os.Stdout, logging, logging.Level[LogComponentApiTestingDefault])

	return Logger{
		Logger:        zapr.NewLogger(logger),
		logging:       logging,
		sugaredLogger: logger.Sugar(),
	}
}

// FileLogger returns a file logger.
// file is the path of the log file.
// name is the name of the logger.
// level is the log level of the logger.
// The returned logger can write logs to the specified file.
func FileLogger(file string, name string, level LogLevel) Logger {
	writer, err := os.OpenFile(file, os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}

	logging := DefaultApiTestingLogging()
	logger := initZapLogger(writer, logging, level)

	return Logger{
		Logger:        zapr.NewLogger(logger).WithName(name),
		logging:       logging,
		sugaredLogger: logger.Sugar(),
	}
}

func DefaultLogger(level LogLevel) Logger {
	logging := DefaultApiTestingLogging()
	logger := initZapLogger(os.Stdout, logging, level)

	return Logger{
		Logger:        zapr.NewLogger(logger),
		logging:       logging,
		sugaredLogger: logger.Sugar(),
	}
}

// WithName returns a new Logger instance with the specified name element added
// to the Logger's name.  Successive calls with WithName append additional
// suffixes to the Logger's name.  It's strongly recommended that name segments
// contain only letters, digits, and hyphens (see the package documentation for
// more information).
func (l Logger) WithName(name string) Logger {
	logLevel := l.logging.Level[apiTestingLogComponent(name)]
	logger := initZapLogger(os.Stdout, l.logging, logLevel)

	return Logger{
		Logger:        zapr.NewLogger(logger).WithName(name),
		logging:       l.logging,
		sugaredLogger: logger.Sugar(),
	}
}

// WithValues returns a new Logger instance with additional key/value pairs.
// See Info for documentation on how key/value pairs work.
func (l Logger) WithValues(keysAndValues ...interface{}) Logger {

	l.Logger = l.Logger.WithValues(keysAndValues...)
	return l
}

// A Sugar wraps the base Logger functionality in a slower, but less
// verbose, API. Any Logger can be converted to a SugaredLogger with its Sugar
// method.
//
// Unlike the Logger, the SugaredLogger doesn't insist on structured logging.
// For each log level, it exposes four methods:
//
//   - methods named after the log level for log.Print-style logging
//   - methods ending in "w" for loosely-typed structured logging
//   - methods ending in "f" for log.Printf-style logging
//   - methods ending in "ln" for log.Println-style logging
//
// Used:
//
//	Info(...any)           Print-style logging
//	Infow(...any)          Structured logging (read as "info with")
//	Infof(string, ...any)  Printf-style logging
//	Infoln(...any)         Println-style logging
func (l Logger) Sugar() *zap.SugaredLogger {

	return l.sugaredLogger
}

func initZapLogger(w io.Writer, logging *ApiTestingLogging, level LogLevel) *zap.Logger {

	parseLevel, _ := zapcore.ParseLevel(string(logging.DefaultApiTestingLoggingLevel(level)))
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()), zapcore.AddSync(w), zap.NewAtomicLevelAt(parseLevel))

	return zap.New(core, zap.AddCaller())
}

// DefaultApiTestingLogging returns a new apiTestingLogging with default configuration parameters.
func DefaultApiTestingLogging() *ApiTestingLogging {
	return &ApiTestingLogging{
		Level: map[apiTestingLogComponent]LogLevel{
			LogComponentApiTestingDefault: LogLevelInfo,
		},
	}
}

// DefaultApiTestingLoggingLevel returns a new ApiTestingLogging with default configuration parameters.
// When LogComponentApiTestingDefault specified, all other logging components are ignored.
func (logging *ApiTestingLogging) DefaultApiTestingLoggingLevel(level LogLevel) LogLevel {

	if level != "" {
		return level
	}

	if logging.Level[LogComponentApiTestingDefault] != "" {
		return logging.Level[LogComponentApiTestingDefault]
	}

	return LogLevelInfo
}
