package log

import (
	"fmt"
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

func NewMockLogger() *LoggerService {
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

	Register(&CmdLogger{})
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

func GetMockLogger() (*MockLogger, error) {
	for _, logger := range globalLogger.Loggers {
		if logger, ok := logger.(*MockLogger); ok {
			return logger, nil
		}
	}

	return nil, fmt.Errorf("MockLogger not found")
}
