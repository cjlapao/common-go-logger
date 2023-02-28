package log

import (
	"fmt"
	"os"
	"strings"

	"github.com/cjlapao/common-go-logger/constants"
	"github.com/cjlapao/common-go-logger/entities"
	"github.com/cjlapao/common-go-logger/icons"
	"github.com/cjlapao/common-go-logger/interfaces"
	strcolor "github.com/cjlapao/common-go/strcolor"
)

var globalLogger *Logger

// Logger Default structure
type Logger struct {
	Loggers        []interfaces.Logger
	LogLevel       entities.Level
	HighlightColor strcolor.ColorCode
	UseTimestamp   bool
}

// Get Creates a new Logger instance
func Get() *Logger {
	if globalLogger == nil {
		result := Logger{
			LogLevel:       entities.Info,
			HighlightColor: strcolor.BrightYellow,
		}
		result.Loggers = []interfaces.Logger{}
		result.AddCmdLogger()

		_logLevel := os.Getenv(constants.LOG_LEVEL)
		if _logLevel == "debug" {
			result.LogLevel = entities.Debug
		}

		if _logLevel == "trace" {
			result.LogLevel = entities.Trace
		}

		globalLogger = &result
		return &result
	}

	return globalLogger
}

func Register[T interfaces.Logger](value T) {
	l := Get()
	found := false
	newType := fmt.Sprintf("%T", value)
	for _, logger := range l.Loggers {
		xType := fmt.Sprintf("%T", logger)
		if strings.EqualFold(newType, xType) {
			found = true
			break
		}
	}

	if !found {
		logger := value.Init()
		l.Loggers = append(l.Loggers, logger)
	}
}

// AddCmdLogger Add a command line logger to the system
func (l *Logger) AddCmdLogger() {
	Register(&CmdLogger{})
}

// func (l *Logger) AddCmdLoggerWithTimestamp() {
// 	found := false
// 	for _, logger := range l.Loggers {
// 		xType := fmt.Sprintf("%T", logger)
// 		if xType == "CmdLogger" {
// 			found = true
// 			logger.UseTimestamp(true)
// 			break
// 		}
// 	}

// 	if !found {
// 		logger := new(CmdLogger)
// 		logger.UseTimestamp(true)
// 		l.Loggers = append(l.Loggers, logger)
// 	}
// }

func (l *Logger) WithDebug() *Logger {
	l.LogLevel = entities.Debug
	return l
}

func (l *Logger) WithTrace() *Logger {
	l.LogLevel = entities.Trace
	return l
}

func (l *Logger) WithWarning() *Logger {
	l.LogLevel = entities.Warning
	return l
}

func (l *Logger) WithTimestamp() *Logger {
	for _, logger := range l.Loggers {
		logger.UseTimestamp(true)
	}
	return l
}

func (l *Logger) WithCorrelationId() *Logger {
	for _, logger := range l.Loggers {
		logger.UseCorrelationId(true)
	}
	return l
}

func (l *Logger) WithIcons() *Logger {
	for _, logger := range l.Loggers {
		logger.UseIcons(true)
	}
	return l
}

// Log Log information message
func (l *Logger) Log(format string, level entities.Level, words ...interface{}) {
	for _, logger := range l.Loggers {
		logger.Log(format, level, words...)
	}
}

// LogIcon Log message with custom icon
func (l *Logger) LogIcon(format string, icon icons.LoggerIcon, level entities.Level, words ...interface{}) {
	for _, logger := range l.Loggers {
		logger.Log(format, level, words...)
	}
}

// LogHighlight Log information message
func (l *Logger) LogHighlight(format string, level entities.Level, words ...interface{}) {
	for _, logger := range l.Loggers {
		logger.LogHighlight(format, level, l.HighlightColor, words...)
	}
}

// Info log information message
func (l *Logger) Info(format string, words ...interface{}) {
	if l.LogLevel >= entities.Info {
		for _, logger := range l.Loggers {
			logger.Info(format, words...)
		}
	}
}

// Success log message
func (l *Logger) Success(format string, words ...interface{}) {
	if l.LogLevel >= entities.Info {
		for _, logger := range l.Loggers {
			logger.Success(format, words...)
		}
	}
}

// TaskSuccess log message
func (l *Logger) TaskSuccess(format string, isComplete bool, words ...interface{}) {
	if l.LogLevel >= entities.Info {
		for _, logger := range l.Loggers {
			logger.TaskSuccess(format, isComplete, words...)
		}
	}
}

// Warn log message
func (l *Logger) Warn(format string, words ...interface{}) {
	if l.LogLevel >= entities.Warning {
		for _, logger := range l.Loggers {
			logger.Warn(format, words...)
		}
	}
}

// TaskWarn log message
func (l *Logger) TaskWarn(format string, words ...interface{}) {
	if l.LogLevel >= entities.Warning {
		for _, logger := range l.Loggers {
			logger.TaskWarn(format, words...)
		}
	}
}

// Command log message
func (l *Logger) Command(format string, words ...interface{}) {
	if l.LogLevel >= entities.Info {
		for _, logger := range l.Loggers {
			logger.Command(format, words...)
		}
	}
}

// Disabled log message
func (l *Logger) Disabled(format string, words ...interface{}) {
	if l.LogLevel >= entities.Info {
		for _, logger := range l.Loggers {
			logger.Disabled(format, words...)
		}
	}
}

// Notice log message
func (l *Logger) Notice(format string, words ...interface{}) {
	if l.LogLevel >= entities.Info {
		for _, logger := range l.Loggers {
			logger.Notice(format, words...)
		}
	}
}

// Debug log message
func (l *Logger) Debug(format string, words ...interface{}) {
	if l.LogLevel >= entities.Debug {
		for _, logger := range l.Loggers {
			logger.Debug(format, words...)
		}
	}
}

// Trace log message
func (l *Logger) Trace(format string, words ...interface{}) {
	if l.LogLevel >= entities.Trace {
		for _, logger := range l.Loggers {
			logger.Debug(format, words...)
		}
	}
}

// Error log message
func (l *Logger) Error(format string, words ...interface{}) {
	if l.LogLevel >= entities.Error {
		for _, logger := range l.Loggers {
			logger.Error(format, words...)
		}
	}
}

// LogError log message
func (l *Logger) LogError(message error) {
	if l.LogLevel >= entities.Error {
		if message != nil {
			for _, logger := range l.Loggers {
				logger.Error(message.Error())
			}
		}
	}
}

// Exception log message
func (l *Logger) Exception(err error, format string, words ...interface{}) {
	if l.LogLevel >= entities.Error {
		for _, logger := range l.Loggers {
			logger.Exception(err, format, words...)
		}
	}
}

// TaskError log message
func (l *Logger) TaskError(format string, isComplete bool, words ...interface{}) {
	if l.LogLevel >= entities.Error {
		for _, logger := range l.Loggers {
			logger.TaskError(format, isComplete, words...)
		}
	}
}

// Fatal log message
func (l *Logger) Fatal(format string, words ...interface{}) {
	if l.LogLevel >= entities.Error {
		for _, logger := range l.Loggers {
			logger.Fatal(format, words...)
		}
	}
}

// FatalError log message
func (l *Logger) FatalError(e error, format string, words ...interface{}) {
	for _, logger := range l.Loggers {
		logger.Error(format, words...)
	}

	if e != nil {
		panic(e)
	}
}
