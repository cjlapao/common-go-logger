package log

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	strcolor "github.com/cjlapao/common-go/strcolor"

	"github.com/fatih/color"
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
func (l *CmdLogger) LogIcon(icon LoggerIcon, format string, level Level, words ...interface{}) {
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
func (l *CmdLogger) LogHighlight(format string, level Level, highlightColor strcolor.ColorCode, words ...interface{}) {
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
func (l *CmdLogger) Info(format string, words ...interface{}) {
	l.printMessage(format, IconInfo, "info", false, false, words...)
}

// Success log message
func (l *CmdLogger) Success(format string, words ...interface{}) {
	l.printMessage(format, IconThumbsUp, "success", false, false, words...)
}

// TaskSuccess log message
func (l *CmdLogger) TaskSuccess(format string, isComplete bool, words ...interface{}) {
	l.printMessage(format, "", "success", true, isComplete, words...)
}

// Warn log message
func (l *CmdLogger) Warn(format string, words ...interface{}) {
	l.printMessage(format, IconWarning, "warn", false, false, words...)
}

// TaskWarn log message
func (l *CmdLogger) TaskWarn(format string, words ...interface{}) {
	l.printMessage(format, "", "warn", true, false, words...)
}

// Command log message
func (l *CmdLogger) Command(format string, words ...interface{}) {
	l.printMessage(format, IconWrench, "command", false, false, words...)
}

// Disabled log message
func (l *CmdLogger) Disabled(format string, words ...interface{}) {
	l.printMessage(format, IconBlackSquare, "disabled", false, false, words...)
}

// Notice log message
func (l *CmdLogger) Notice(format string, words ...interface{}) {
	l.printMessage(format, IconFlag, "notice", false, false, words...)
}

// Debug log message
func (l *CmdLogger) Debug(format string, words ...interface{}) {
	l.printMessage(format, IconFire, "debug", false, false, words...)
}

// Trace log message
func (l *CmdLogger) Trace(format string, words ...interface{}) {
	l.printMessage(format, IconBulb, "trace", false, false, words...)
}

// Error log message
func (l *CmdLogger) Error(format string, words ...interface{}) {
	l.printMessage(format, IconRevolvingLight, "error", false, false, words...)
}

// Error log message
func (l *CmdLogger) Exception(err error, format string, words ...interface{}) {
	if format == "" {
		format = err.Error()
	} else {
		format = format + ", err " + err.Error()
	}
	l.printMessage(format, IconRevolvingLight, "error", false, false, words...)
}

// LogError log message
func (l *CmdLogger) LogError(message error) {
	if message != nil {
		l.printMessage(message.Error(), IconRevolvingLight, "error", false, false)
	}
}

// TaskError log message
func (l *CmdLogger) TaskError(format string, isComplete bool, words ...interface{}) {
	l.printMessage(format, "", "error", true, isComplete, l.useTimestamp)
}

// Fatal log message
func (l *CmdLogger) Fatal(format string, words ...interface{}) {
	l.printMessage(format, IconRevolvingLight, "error", false, true, words...)
}

// FatalError log message
func (l *CmdLogger) FatalError(e error, format string, words ...interface{}) {
	l.Error(format, words...)
	if e != nil {
		panic(e)
	}
}

// printMessage Prints a message in the system
func (l *CmdLogger) printMessage(format string, icon LoggerIcon, level string, isTask bool, isComplete bool, words ...interface{}) {
	if l.useIcons && icon != "" {
		format = fmt.Sprintf("%s %s", icon, format)
	}

	if l.userCorrelationId {
		correlationId := os.Getenv("CORRELATION_ID")
		if correlationId != "" {
			format = "[" + correlationId + "] " + format
		}
	}

	if l.useTimestamp {
		format = fmt.Sprintf("%s %s", time.Now().Format(time.RFC3339), format)
	}

	format = format + "\u001b[0m" + "\n"

	successWriter := color.New(SuccessColor).FprintfFunc()
	warningWriter := color.New(WarningColor).FprintfFunc()
	errorWriter := color.New(ErrorColor).FprintfFunc()
	debugWriter := color.New(DebugColor).FprintfFunc()
	traceWriter := color.New(TraceColor).FprintfFunc()
	infoWriter := color.New(InfoColor).FprintfFunc()
	noticeWriter := color.New(NoticeColor).FprintfFunc()
	commandWriter := color.New(CommandColor).FprintfFunc()
	disableWriter := color.New(DisabledColor).FprintfFunc()

	formattedWords := make([]interface{}, len(words))
	if len(words) > 0 {
		for i := range words {
			word := ""
			switch t := words[i].(type) {
			case string:
				word = t
			default:
				word = fmt.Sprintf("%v", words[i])
			}

			if word != "" {
				word = strings.ReplaceAll(word, "\n\n", "\n")
				if word[0] == 27 {
					switch strings.ToLower(level) {
					case "success":
						word += "\u001b[" + fmt.Sprint(SuccessColor) + "m"
					case "warn":
						word += "\u001b[" + fmt.Sprint(WarningColor) + "m"
					case "error":
						word += "\u001b[" + fmt.Sprint(ErrorColor) + "m"
					case "debug":
						word += "\u001b[" + fmt.Sprint(DebugColor) + "m"
					case "trace":
						word += "\u001b[" + fmt.Sprint(TraceColor) + "m"
					case "info":
						word += "\u001b[" + fmt.Sprint(InfoColor) + "m"
					case "notice":
						word += "\u001b[" + fmt.Sprint(NoticeColor) + "m"
					case "command":
						word += "\u001b[" + fmt.Sprint(CommandColor) + "m"
					case "disabled":
						word += "\u001b[" + fmt.Sprint(DisabledColor) + "m"
					}
					formattedWords[i] = word
				} else {
					formattedWords[i] = word
				}
			}
		}
	}

	switch strings.ToLower(level) {
	case "success":
		successWriter(l.writer, format, formattedWords...)

		if isComplete {
			successWriter(l.writer, "Completed")
			os.Exit(0)
		}
	case "warn":
		warningWriter(l.writer, format, formattedWords...)
	case "error":
		errorWriter(l.writer, format, formattedWords...)

		if isComplete {
			errorWriter(l.writer, "Failed\n")
			os.Exit(1)
		}
	case "debug":
		debugWriter(l.writer, format, formattedWords...)
	case "trace":
		traceWriter(l.writer, format, formattedWords...)
	case "info":
		infoWriter(l.writer, format, formattedWords...)
	case "notice":
		noticeWriter(l.writer, format, formattedWords...)
	case "command":
		commandWriter(l.writer, format, formattedWords...)
	case "disabled":
		disableWriter(l.writer, format, formattedWords...)
	}
}
