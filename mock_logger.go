package log

import (
	"fmt"
	"io"
	"os"

	"github.com/cjlapao/common-go/strcolor"
)

type MockedLogMessage struct {
	Message string
	Level   string
	Icon    string
}

type MockLogger struct {
	LastPrintedMessage MockedLogMessage
	PrintedMessages    []MockedLogMessage
	LastCallType       string
	useTimestamp       bool
	userCorrelationId  bool
	useIcons           bool
	writer             io.Writer
}

func (l MockLogger) Init() Logger {
	return &MockLogger{
		useTimestamp:       false,
		userCorrelationId:  false,
		useIcons:           false,
		writer:             os.Stdout,
		LastPrintedMessage: MockedLogMessage{},
		PrintedMessages:    []MockedLogMessage{},
	}
}

func (l *MockLogger) Clear() {
	l.LastPrintedMessage = MockedLogMessage{}
	l.PrintedMessages = []MockedLogMessage{}
}

func (l *MockLogger) IsTimestampEnabled() bool {
	return l.useTimestamp
}

func (l *MockLogger) UseTimestamp(value bool) {
	l.useTimestamp = value
}

func (l *MockLogger) UseCorrelationId(value bool) {
	l.userCorrelationId = value
}

func (l *MockLogger) UseIcons(value bool) {
	l.useIcons = value
}

func (l *MockLogger) Log(format string, level Level, words ...interface{}) {
	switch level {
	case 0:
		l.printMessage(format, "", "error", false, false, words...)
	case 1:
		l.printMessage(format, "", "warn", false, false, words...)
	case 2:
		l.printMessage(format, "", "info", false, false, words...)
	case 3:
		l.printMessage(format, "", "debug", false, false, words...)
	case 4:
		l.printMessage(format, "", "trace", false, false, words...)
	}
}

// Log Log information message
func (l *MockLogger) LogIcon(icon LoggerIcon, format string, level Level, words ...interface{}) {
	switch level {
	case 0:
		l.printMessage(format, icon, "error", false, false, words...)
	case 1:
		l.printMessage(format, icon, "warn", false, false, words...)
	case 2:
		l.printMessage(format, icon, "info", false, false, words...)
	case 3:
		l.printMessage(format, icon, "debug", false, false, words...)
	case 4:
		l.printMessage(format, icon, "trace", false, false, words...)
	}
}

// LogHighlight Log information message
func (l *MockLogger) LogHighlight(format string, level Level, highlightColor strcolor.ColorCode, words ...interface{}) {
	if len(words) > 0 {
		for i := range words {
			words[i] = strcolor.GetColorString(strcolor.ColorCode(highlightColor), fmt.Sprintf("%s", words[i]))
		}
	}

	switch level {
	case 0:
		l.printMessage(format, "", "error", false, false, words...)
	case 1:
		l.printMessage(format, "", "warn", false, false, words...)
	case 2:
		l.printMessage(format, "", "info", false, false, words...)
	case 3:
		l.printMessage(format, "", "debug", false, false, words...)
	case 4:
		l.printMessage(format, "", "trace", false, false, words...)
	}
}

// Info log information message
func (l *MockLogger) Info(format string, words ...interface{}) {
	l.printMessage(format, IconInfo, "info", false, false, words...)
}

// Success log message
func (l *MockLogger) Success(format string, words ...interface{}) {
	l.printMessage(format, IconThumbsUp, "success", false, false, words...)
}

// TaskSuccess log message
func (l *MockLogger) TaskSuccess(format string, isComplete bool, words ...interface{}) {
	l.printMessage(format, "", "success", true, isComplete, words...)
}

// Warn log message
func (l *MockLogger) Warn(format string, words ...interface{}) {
	l.printMessage(format, IconWarning, "warn", false, false, words...)
}

// TaskWarn log message
func (l *MockLogger) TaskWarn(format string, words ...interface{}) {
	l.printMessage(format, "", "warn", true, false, words...)
}

// Command log message
func (l *MockLogger) Command(format string, words ...interface{}) {
	l.printMessage(format, IconWrench, "command", false, false, words...)
}

// Disabled log message
func (l *MockLogger) Disabled(format string, words ...interface{}) {
	l.printMessage(format, IconBlackSquare, "disabled", false, false, words...)
}

// Notice log message
func (l *MockLogger) Notice(format string, words ...interface{}) {
	l.printMessage(format, IconFlag, "notice", false, false, words...)
}

// Debug log message
func (l *MockLogger) Debug(format string, words ...interface{}) {
	l.printMessage(format, IconFire, "debug", false, false, words...)
}

// Trace log message
func (l *MockLogger) Trace(format string, words ...interface{}) {
	l.printMessage(format, IconBulb, "trace", false, false, words...)
}

// Error log message
func (l *MockLogger) Error(format string, words ...interface{}) {
	l.printMessage(format, IconRevolvingLight, "error", false, false, words...)
}

// Error log message
func (l *MockLogger) Exception(err error, format string, words ...interface{}) {
	if format == "" {
		format = err.Error()
	} else {
		format = format + ", err " + err.Error()
	}
	l.printMessage(format, IconRevolvingLight, "error", false, false, words...)
}

// LogError log message
func (l *MockLogger) LogError(message error) {
	if message != nil {
		l.printMessage(message.Error(), IconRevolvingLight, "error", false, false)
	}
}

// TaskError log message
func (l *MockLogger) TaskError(format string, isComplete bool, words ...interface{}) {
	l.printMessage(format, "", "error", true, isComplete, l.useTimestamp)
}

// Fatal log message
func (l *MockLogger) Fatal(format string, words ...interface{}) {
	l.printMessage(format, IconRevolvingLight, "error", false, true, words...)
}

// FatalError log message
func (l *MockLogger) FatalError(e error, format string, words ...interface{}) {
	l.Error(format, words...)
	if e != nil {
		panic(e)
	}
}

// printMessage Prints a message in the system
func (l *MockLogger) printMessage(format string, icon LoggerIcon, level string, isTask bool, isComplete bool, words ...interface{}) {
	l.LastPrintedMessage = MockedLogMessage{Message: fmt.Sprintf(format, words...), Level: level, Icon: string(icon)}
	l.PrintedMessages = append(l.PrintedMessages, l.LastPrintedMessage)
}
