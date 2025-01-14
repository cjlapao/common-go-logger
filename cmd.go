package log

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	strcolor "github.com/cjlapao/common-go/strcolor"
)

// CmdLogger Command Line Logger implementation
type CmdLogger struct {
	useTimestamp      bool
	userCorrelationId bool
	useIcons          bool
	writer            io.Writer
}

func (l CmdLogger) Init() Logger {
	return &CmdLogger{
		useTimestamp:      false,
		userCorrelationId: false,
		useIcons:          false,
		writer:            os.Stdout,
	}
}

func (l *CmdLogger) IsTimestampEnabled() bool {
	return l.useTimestamp
}

func (l *CmdLogger) UseTimestamp(value bool) {
	l.useTimestamp = value
}

func (l *CmdLogger) UseCorrelationId(value bool) {
	l.userCorrelationId = value
}

func (l *CmdLogger) UseIcons(value bool) {
	l.useIcons = value
}

// Log Log information message
func (l *CmdLogger) Log(format string, level Level, words ...interface{}) {
	switch level {
	case 0:
		l.printMessage(format, "", "error", words...)
	case 1:
		l.printMessage(format, "", "warn", words...)
	case 2:
		l.printMessage(format, "", "info", words...)
	case 3:
		l.printMessage(format, "", "debug", words...)
	case 4:
		l.printMessage(format, "", "trace", words...)
	}
}

// Log Log information message
func (l *CmdLogger) LogIcon(icon LoggerIcon, format string, level Level, words ...interface{}) {
	switch level {
	case 0:
		l.printMessage(format, icon, "error", words...)
	case 1:
		l.printMessage(format, icon, "warn", words...)
	case 2:
		l.printMessage(format, icon, "info", words...)
	case 3:
		l.printMessage(format, icon, "debug", words...)
	case 4:
		l.printMessage(format, icon, "trace", words...)
	}
}

// LogHighlight Log information message
func (l *CmdLogger) LogHighlight(format string, level Level, highlightColor strcolor.ColorCode, words ...interface{}) {
	if len(words) > 0 {
		for i := range words {
			words[i] = GetColorString(ColorCode(highlightColor), fmt.Sprintf("%v", words[i]))
		}
	}

	switch level {
	case 0:
		l.printMessage(format, "", "error", words...)
	case 1:
		l.printMessage(format, "", "warn", words...)
	case 2:
		l.printMessage(format, "", "info", words...)
	case 3:
		l.printMessage(format, "", "debug", words...)
	case 4:
		l.printMessage(format, "", "trace", words...)
	}
}

// Info log information message
func (l *CmdLogger) Info(format string, words ...interface{}) {
	l.printMessage(format, IconInfo, "info", words...)
}

// Success log message
func (l *CmdLogger) Success(format string, words ...interface{}) {
	l.printMessage(format, IconThumbsUp, "success", words...)
}

// Warn log message
func (l *CmdLogger) Warn(format string, words ...interface{}) {
	l.printMessage(format, IconWarning, "warn", words...)
}

// Command log message
func (l *CmdLogger) Command(format string, words ...interface{}) {
	l.printMessage(format, IconWrench, "command", words...)
}

// Disabled log message
func (l *CmdLogger) Disabled(format string, words ...interface{}) {
	l.printMessage(format, IconBlackSquare, "disabled", words...)
}

// Notice log message
func (l *CmdLogger) Notice(format string, words ...interface{}) {
	l.printMessage(format, IconFlag, "notice", words...)
}

// Debug log message
func (l *CmdLogger) Debug(format string, words ...interface{}) {
	l.printMessage(format, IconFire, "debug", words...)
}

// Trace log message
func (l *CmdLogger) Trace(format string, words ...interface{}) {
	l.printMessage(format, IconBulb, "trace", words...)
}

// Error log message
func (l *CmdLogger) Error(format string, words ...interface{}) {
	l.printMessage(format, IconRevolvingLight, "error", words...)
}

// Error log message
func (l *CmdLogger) Exception(err error, format string, words ...interface{}) {
	if format == "" {
		format = err.Error()
	} else {
		format = format + ", err " + err.Error()
	}
	l.printMessage(format, IconRevolvingLight, "error", words...)
}

// LogError log message
func (l *CmdLogger) LogError(message error) {
	if message != nil {
		l.printMessage(message.Error(), IconRevolvingLight, "error")
	}
}

// Fatal log message
func (l *CmdLogger) Fatal(format string, words ...interface{}) {
	l.printMessage(format, IconRevolvingLight, "error", words...)
}

// FatalError log message
func (l *CmdLogger) FatalError(e error, format string, words ...interface{}) {
	l.Error(format, words...)
	if e != nil {
		panic(e)
	}
}

// printMessage Prints a message in the system
func (l *CmdLogger) printMessage(format string, icon LoggerIcon, level string, words ...interface{}) {
	// First format the arguments according to the format string
	message := fmt.Sprintf(format, words...)

	if l.useIcons && icon != "" {
		message = fmt.Sprintf("%s %s", icon, message)
	}

	if l.userCorrelationId {
		correlationId := os.Getenv("CORRELATION_ID")
		if correlationId != "" {
			message = "[" + correlationId + "] " + message
		}
	}

	if l.useTimestamp {
		message = fmt.Sprintf("%s %s", time.Now().Format(time.RFC3339), message)
	}

	message = message + "\u001b[0m" + "\n"

	// Use the appropriate color writer for each log level
	switch strings.ToLower(level) {
	case "success":
		successWriter(l.writer, message)
	case "warn":
		warningWriter(l.writer, message)
	case "error":
		errorWriter(l.writer, message)
	case "debug":
		debugWriter(l.writer, message)
	case "trace":
		traceWriter(l.writer, message)
	case "info":
		infoWriter(l.writer, message)
	case "notice":
		noticeWriter(l.writer, message)
	case "command":
		commandWriter(l.writer, message)
	case "disabled":
		disableWriter(l.writer, message)
	}
}

func successWriter(w io.Writer, message string) {
	fmt.Fprintf(w, "\u001b[32m%s", message)
}

func warningWriter(w io.Writer, message string) {
	fmt.Fprintf(w, "\u001b[33m%s", message)
}

func errorWriter(w io.Writer, message string) {
	fmt.Fprintf(w, "\u001b[31m%s", message)
}

func debugWriter(w io.Writer, message string) {
	fmt.Fprintf(w, "\u001b[36m%s", message)
}

func traceWriter(w io.Writer, message string) {
	fmt.Fprintf(w, "\u001b[37m%s", message)
}

func infoWriter(w io.Writer, message string) {
	fmt.Fprintf(w, "\u001b[0m%s", message)
}

func noticeWriter(w io.Writer, message string) {
	fmt.Fprintf(w, "\u001b[34m%s", message)
}

func commandWriter(w io.Writer, message string) {
	fmt.Fprintf(w, "\u001b[35m%s", message)
}

func disableWriter(w io.Writer, message string) {
	fmt.Fprintf(w, "\u001b[90m%s", message)
}
