package log

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	strcolor "github.com/cjlapao/common-go/strcolor"
)

var globalLogger *LoggerService

// Logger Default structure
type LoggerService struct {
	Loggers        []Logger
	LogLevel       Level
	HighlightColor strcolor.ColorCode
	UseTimestamp   bool
}

// Get Creates a new Logger instance
func Get() *LoggerService {
	if globalLogger == nil {
		return New()
	}

	return globalLogger
}

func New() *LoggerService {
	globalLogger = &LoggerService{
		LogLevel:       Info,
		HighlightColor: strcolor.BrightYellow,
		Loggers:        []Logger{},
	}

	_logLevel := os.Getenv(LOG_LEVEL)
	if _logLevel == "debug" {
		globalLogger.LogLevel = Debug
	}

	if _logLevel == "trace" {
		globalLogger.LogLevel = Trace
	}

	globalLogger.AddCmdLogger()
	return globalLogger
}

func Register[T Logger](value T) {
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
func (l *LoggerService) AddCmdLogger() {
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

func (l *LoggerService) WithDebug() *LoggerService {
	l.LogLevel = Debug
	return l
}

func (l *LoggerService) WithTrace() *LoggerService {
	l.LogLevel = Trace
	return l
}

func (l *LoggerService) WithWarning() *LoggerService {
	l.LogLevel = Warning
	return l
}

func (l *LoggerService) WithTimestamp() *LoggerService {
	for _, logger := range l.Loggers {
		logger.UseTimestamp(true)
	}
	return l
}

func (l *LoggerService) ToggleTimestamp() *LoggerService {
	for _, logger := range l.Loggers {
		logger.UseTimestamp(!logger.IsTimestampEnabled())
	}
	return l
}

func (l *LoggerService) EnableTimestamp(value bool) *LoggerService {
	for _, logger := range l.Loggers {
		logger.UseTimestamp(value)
	}

	return l
}

func (l *LoggerService) WithCorrelationId() *LoggerService {
	for _, logger := range l.Loggers {
		logger.UseCorrelationId(true)
	}
	return l
}

func (l *LoggerService) WithIcons() *LoggerService {
	for _, logger := range l.Loggers {
		logger.UseIcons(true)
	}
	return l
}

// Log Log information message
func (l *LoggerService) Log(format string, level Level, words ...interface{}) {
	for _, logger := range l.Loggers {
		logger.Log(format, level, words...)
	}
}

// LogIcon Log message with custom icon
func (l *LoggerService) LogIcon(icon LoggerIcon, format string, level Level, words ...interface{}) {
	for _, logger := range l.Loggers {
		logger.LogIcon(icon, format, level, words...)
	}
}

// LogHighlight Log information message
func (l *LoggerService) LogHighlight(format string, level Level, words ...interface{}) {
	for _, logger := range l.Loggers {
		logger.LogHighlight(format, level, l.HighlightColor, words...)
	}
}

// Info log information message
func (l *LoggerService) Info(format string, words ...interface{}) {
	if l.LogLevel >= Info {
		for _, logger := range l.Loggers {
			logger.Info(format, words...)
		}
	}
}

// Success log message
func (l *LoggerService) Success(format string, words ...interface{}) {
	if l.LogLevel >= Info {
		for _, logger := range l.Loggers {
			logger.Success(format, words...)
		}
	}
}

// TaskSuccess log message
func (l *LoggerService) TaskSuccess(format string, isComplete bool, words ...interface{}) {
	if l.LogLevel >= Info {
		for _, logger := range l.Loggers {
			logger.TaskSuccess(format, isComplete, words...)
		}
	}
}

// Warn log message
func (l *LoggerService) Warn(format string, words ...interface{}) {
	if l.LogLevel >= Warning {
		for _, logger := range l.Loggers {
			logger.Warn(format, words...)
		}
	}
}

// TaskWarn log message
func (l *LoggerService) TaskWarn(format string, words ...interface{}) {
	if l.LogLevel >= Warning {
		for _, logger := range l.Loggers {
			logger.TaskWarn(format, words...)
		}
	}
}

// Command log message
func (l *LoggerService) Command(format string, words ...interface{}) {
	if l.LogLevel >= Info {
		for _, logger := range l.Loggers {
			logger.Command(format, words...)
		}
	}
}

// Disabled log message
func (l *LoggerService) Disabled(format string, words ...interface{}) {
	if l.LogLevel >= Info {
		for _, logger := range l.Loggers {
			logger.Disabled(format, words...)
		}
	}
}

// Notice log message
func (l *LoggerService) Notice(format string, words ...interface{}) {
	if l.LogLevel >= Info {
		for _, logger := range l.Loggers {
			logger.Notice(format, words...)
		}
	}
}

// Debug log message
func (l *LoggerService) Debug(format string, words ...interface{}) {
	if l.LogLevel >= Debug {
		for _, logger := range l.Loggers {
			logger.Debug(format, words...)
		}
	}
}

// Trace log message
func (l *LoggerService) Trace(format string, words ...interface{}) {
	if l.LogLevel >= Trace {
		for _, logger := range l.Loggers {
			logger.Debug(format, words...)
		}
	}
}

// Error log message
func (l *LoggerService) Error(format string, words ...interface{}) {
	if l.LogLevel >= Error {
		for _, logger := range l.Loggers {
			logger.Error(format, words...)
		}
	}
}

// LogError log message
func (l *LoggerService) LogError(message error) {
	if l.LogLevel >= Error {
		if message != nil {
			for _, logger := range l.Loggers {
				logger.Error(message.Error())
			}
		}
	}
}

// Exception log message
func (l *LoggerService) Exception(err error, format string, words ...interface{}) {
	if l.LogLevel >= Error {
		for _, logger := range l.Loggers {
			logger.Exception(err, format, words...)
		}
	}
}

// TaskError log message
func (l *LoggerService) TaskError(format string, isComplete bool, words ...interface{}) {
	if l.LogLevel >= Error {
		for _, logger := range l.Loggers {
			logger.TaskError(format, isComplete, words...)
		}
	}
}

// Fatal log message
func (l *LoggerService) Fatal(format string, words ...interface{}) {
	if l.LogLevel >= Error {
		for _, logger := range l.Loggers {
			logger.Fatal(format, words...)
		}
	}
}

// FatalError log message
func (l *LoggerService) FatalError(e error, format string, words ...interface{}) {
	for _, logger := range l.Loggers {
		logger.Error(format, words...)
	}

	if e != nil {
		panic(e)
	}
}

func (l *LoggerService) GetRequestPrefix(r *http.Request, logUrl bool) string {
	msg := ""
	if r.Header.Get("X-Request-Id") != "" {
		msg += fmt.Sprintf("[%s] ", r.Header.Get("X-Request-Id"))
	}

	if logUrl {
		msg += fmt.Sprintf("[%s] [%s] ", r.Method, r.URL.Path)
	}

	return msg
}
