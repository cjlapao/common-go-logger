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

// AdoCmdLogger Command Line Logger implementation
type AdoCmdLogger struct {
	useTimestamp      bool
	userCorrelationId bool
	useIcons          bool
	writer            io.Writer
}

func (l AdoCmdLogger) Init() Logger {
	return &AdoCmdLogger{
		useTimestamp:      false,
		userCorrelationId: false,
		useIcons:          false,
		writer:            os.Stdout,
	}
}

func (l *AdoCmdLogger) IsTimestampEnabled() bool {
	return l.useTimestamp
}

func (l *AdoCmdLogger) UseTimestamp(value bool) {
	l.useTimestamp = value
}

func (l *AdoCmdLogger) UseCorrelationId(value bool) {
	l.userCorrelationId = value
}

func (l *AdoCmdLogger) UseIcons(value bool) {
	l.useIcons = value
}

// Log Log information message
func (l *AdoCmdLogger) Log(format string, level Level, words ...interface{}) {
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
func (l *AdoCmdLogger) LogIcon(icon LoggerIcon, format string, level Level, words ...interface{}) {
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
func (l *AdoCmdLogger) LogHighlight(format string, level Level, highlightColor strcolor.ColorCode, words ...interface{}) {
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
func (l *AdoCmdLogger) Info(format string, words ...interface{}) {
	l.printMessage(format, IconInfo, "info", false, false, words...)
}

// Success log message
func (l *AdoCmdLogger) Success(format string, words ...interface{}) {
	l.printMessage(format, IconThumbsUp, "success", false, false, words...)
}

// TaskSuccess log message
func (l *AdoCmdLogger) TaskSuccess(format string, isComplete bool, words ...interface{}) {
	l.printMessage(format, "", "success", true, isComplete, words...)
}

// Warn log message
func (l *AdoCmdLogger) Warn(format string, words ...interface{}) {
	l.printMessage(format, IconWarning, "warn", false, false, words...)
}

// TaskWarn log message
func (l *AdoCmdLogger) TaskWarn(format string, words ...interface{}) {
	l.printMessage(format, "", "warn", true, false, words...)
}

// Command log message
func (l *AdoCmdLogger) Command(format string, words ...interface{}) {
	l.printMessage(format, IconWrench, "command", false, false, words...)
}

// Disabled log message
func (l *AdoCmdLogger) Disabled(format string, words ...interface{}) {
	l.printMessage(format, IconBlackSquare, "disabled", false, false, words...)
}

// Notice log message
func (l *AdoCmdLogger) Notice(format string, words ...interface{}) {
	l.printMessage(format, IconFlag, "notice", false, false, words...)
}

// Debug log message
func (l *AdoCmdLogger) Debug(format string, words ...interface{}) {
	l.printMessage(format, IconFire, "debug", false, false, words...)
}

// Trace log message
func (l *AdoCmdLogger) Trace(format string, words ...interface{}) {
	l.printMessage(format, IconBulb, "trace", false, false, words...)
}

// Error log message
func (l *AdoCmdLogger) Error(format string, words ...interface{}) {
	l.printMessage(format, IconRevolvingLight, "error", false, false, words...)
}

// Error log message
func (l *AdoCmdLogger) Exception(err error, format string, words ...interface{}) {
	if format == "" {
		format = err.Error()
	} else {
		format = format + ", err " + err.Error()
	}
	l.printMessage(format, IconRevolvingLight, "error", false, false, words...)
}

// LogError log message
func (l *AdoCmdLogger) LogError(message error) {
	if message != nil {
		l.printMessage(message.Error(), IconRevolvingLight, "error", false, false)
	}
}

// TaskError log message
func (l *AdoCmdLogger) TaskError(format string, isComplete bool, words ...interface{}) {
	l.printMessage(format, "", "error", true, isComplete, l.useTimestamp)
}

// Fatal log message
func (l *AdoCmdLogger) Fatal(format string, words ...interface{}) {
	l.printMessage(format, IconRevolvingLight, "error", false, true, words...)
}

// FatalError log message
func (l *AdoCmdLogger) FatalError(e error, format string, words ...interface{}) {
	l.Error(format, words...)
	if e != nil {
		panic(e)
	}
}

// printMessage Prints a message in the system
func (l *AdoCmdLogger) printMessage(format string, icon LoggerIcon, level string, isTask bool, isComplete bool, words ...interface{}) {
	agentID := os.Getenv("AGENT_ID")
	isPipeline := false
	if len(agentID) != 0 {
		isPipeline = true
	}

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

	if !isPipeline {
		format = format + "\u001b[0m" + "\n"
	} else {
		if (level == "warn" || level == "error") && isTask {
			format = format + "\n"
		} else {
			format = format + "\033[0m" + "\n"
		}
	}

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
						if isPipeline {
							word += "\033[" + fmt.Sprint(SuccessColor) + "m"
						} else {
							word += "\u001b[" + fmt.Sprint(SuccessColor) + "m"
						}
					case "warn":
						if isPipeline {
							if !isTask {
								word += "\033[" + fmt.Sprint(WarningColor) + "m"
							}
						} else {
							word += "\u001b[" + fmt.Sprint(WarningColor) + "m"
						}
					case "error":
						if isPipeline {
							if !isTask {
								word += "\033[" + fmt.Sprint(ErrorColor) + "m"
							}
						} else {
							word += "\u001b[" + fmt.Sprint(ErrorColor) + "m"
						}
					case "debug":
						if isPipeline {
							word += "\033[" + fmt.Sprint(DebugColor) + "m"
						} else {
							word += "\u001b[" + fmt.Sprint(DebugColor) + "m"
						}
					case "trace":
						if isPipeline {
							word += "\033[" + fmt.Sprint(TraceColor) + "m"
						} else {
							word += "\u001b[" + fmt.Sprint(TraceColor) + "m"
						}
					case "info":
						if isPipeline {
							word += "\033[" + fmt.Sprint(InfoColor) + "m"
						} else {
							word += "\u001b[" + fmt.Sprint(InfoColor) + "m"
						}
					case "notice":
						if isPipeline {
							word += "\033[" + fmt.Sprint(NoticeColor) + "m"
						} else {
							word += "\u001b[" + fmt.Sprint(NoticeColor) + "m"
						}
					case "command":
						if isPipeline {
							word += "\033[" + fmt.Sprint(CommandColor) + "m"
						} else {
							word += "\u001b[" + fmt.Sprint(CommandColor) + "m"
						}
					case "disabled":
						if isPipeline {
							word += "\033[" + fmt.Sprint(DisabledColor) + "m"
						} else {
							word += "\u001b[" + fmt.Sprint(DisabledColor) + "m"
						}
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
		if isPipeline {
			format = "\033[" + fmt.Sprint(SuccessColor) + "m" + format
			format = "##[section]" + format
			fmt.Fprintf(l.writer, format, formattedWords...)
			if isTask && isComplete {
				fmt.Fprintf(l.writer, "\033["+fmt.Sprint(SuccessColor)+"m"+"##vso[task.complete result=Succeeded;]\n")
			}
		} else {
			successWriter(l.writer, format, formattedWords...)
		}

		if isComplete {
			if isPipeline && isTask {
				fmt.Fprintf(l.writer, "\033["+fmt.Sprint(SuccessColor)+"m"+"##[section] Completed\n")
			} else {
				successWriter(l.writer, "Completed")
			}
			os.Exit(0)
		}
	case "warn":
		if isPipeline {
			if isTask {
				format = "##vso[task.LogIssue type=warning;]" + format
				fmt.Fprintf(l.writer, format, formattedWords...)
			} else {
				format = "\033[" + fmt.Sprint(WarningColor) + "m" + format
				fmt.Fprintf(l.writer, format, formattedWords...)
			}
		} else {
			warningWriter(l.writer, format, formattedWords...)
		}
	case "error":
		if isPipeline {
			if isTask {
				format = "##vso[task.LogIssue type=error;]" + format
				fmt.Fprintf(l.writer, format, formattedWords...)
			} else {
				format = "\033[" + fmt.Sprint(ErrorColor) + "m" + format
				fmt.Fprintf(l.writer, format, formattedWords...)
			}
		} else {
			errorWriter(l.writer, format, formattedWords...)
		}

		if isComplete {
			if isPipeline && isTask {
				format = "\033[" + fmt.Sprint(ErrorColor) + "m" + format
				fmt.Fprintf(l.writer, format, formattedWords...)
				fmt.Fprintf(l.writer, "##vso[task.complete result=Failed;]\n")
				os.Exit(0)
			} else {
				errorWriter(l.writer, "Failed\n")
				os.Exit(1)
			}
		}
	case "debug":
		if isPipeline {
			format = "\033[" + fmt.Sprint(DebugColor) + "m" + format
			fmt.Fprintf(l.writer, format, formattedWords...)
		} else {
			debugWriter(l.writer, format, formattedWords...)
		}
	case "trace":
		if isPipeline {
			format = "\033[" + fmt.Sprint(TraceColor) + "m" + format
			fmt.Fprintf(l.writer, format, formattedWords...)
		} else {
			traceWriter(l.writer, format, formattedWords...)
		}
	case "info":
		if isPipeline {
			format = "\033[" + fmt.Sprint(InfoColor) + "m" + format
			fmt.Fprintf(l.writer, format, formattedWords...)
		} else {
			infoWriter(l.writer, format, formattedWords...)
		}
	case "notice":
		if isPipeline {
			format = "\033[" + fmt.Sprint(NoticeColor) + "m" + format
			fmt.Fprintf(l.writer, format, formattedWords...)
		} else {
			noticeWriter(l.writer, format, formattedWords...)
		}
	case "command":
		if isPipeline {
			format = "\033[" + fmt.Sprint(CommandColor) + "m" + format
			format = "##[command]" + format
			fmt.Fprintf(l.writer, format, formattedWords...)
		} else {
			commandWriter(l.writer, format, formattedWords...)
		}
	case "disabled":
		if isPipeline {
			format = "\033[" + fmt.Sprint(DisabledColor) + "m" + format
			fmt.Fprintf(l.writer, format, formattedWords...)
		} else {
			disableWriter(l.writer, format, formattedWords...)
		}
	}
}
