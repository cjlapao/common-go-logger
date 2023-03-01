package log

import (
	"github.com/cjlapao/common-go/strcolor"
)

// Logger Interface
type Logger interface {
	UseTimestamp(value bool)
	UseCorrelationId(value bool)
	UseIcons(value bool)

	Init() Logger
	Log(format string, level Level, words ...interface{})
	LogIcon(icon LoggerIcon, format string, level Level, words ...interface{})
	LogHighlight(format string, level Level, highlightColor strcolor.ColorCode, words ...interface{})
	Info(format string, words ...interface{})
	Success(format string, words ...interface{})
	TaskSuccess(format string, isComplete bool, words ...interface{})
	Warn(format string, words ...interface{})
	TaskWarn(format string, words ...interface{})
	Command(format string, words ...interface{})
	Disabled(format string, words ...interface{})
	Notice(format string, words ...interface{})
	Debug(format string, words ...interface{})
	Trace(format string, words ...interface{})
	Error(format string, words ...interface{})
	Exception(err error, format string, words ...interface{})
	LogError(message error)
	TaskError(format string, isComplete bool, words ...interface{})
	Fatal(format string, words ...interface{})
	FatalError(e error, format string, words ...interface{})
}
