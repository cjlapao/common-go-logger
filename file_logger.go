package log

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	strcolor "github.com/cjlapao/common-go/strcolor"
)

// FileLogger Command Line Logger implementation
type FileLogger struct {
	useTimestamp      bool
	userCorrelationId bool
	useIcons          bool
	filename          string
	enabled           bool
	writer            io.Writer
}

func (l FileLogger) Init() Logger {
	logger := &FileLogger{
		useTimestamp:      false,
		userCorrelationId: false,
		useIcons:          false,
		filename:          l.filename,
	}
	if l.filename != "" {
		file, err := os.OpenFile(l.filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
		if err != nil {
			panic(err)
		}
		logger.writer = file
		logger.enabled = true
	} else {
		logger.writer = os.Stdout
		logger.enabled = false
	}
	return logger
}

func (l *FileLogger) IsTimestampEnabled() bool {
	return l.useTimestamp
}

func (l *FileLogger) UseTimestamp(value bool) {
	l.useTimestamp = value
}

func (l *FileLogger) UseCorrelationId(value bool) {
	l.userCorrelationId = value
}

func (l *FileLogger) UseIcons(value bool) {
	l.useIcons = value
}

// Log Log information message
func (l *FileLogger) Log(format string, level Level, words ...interface{}) {
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
func (l *FileLogger) LogIcon(icon LoggerIcon, format string, level Level, words ...interface{}) {
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
func (l *FileLogger) LogHighlight(format string, level Level, highlightColor strcolor.ColorCode, words ...interface{}) {
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
func (l *FileLogger) Info(format string, words ...interface{}) {
	l.printMessage(format, IconInfo, "info", false, false, words...)
}

// Success log message
func (l *FileLogger) Success(format string, words ...interface{}) {
	l.printMessage(format, IconThumbsUp, "success", false, false, words...)
}

// TaskSuccess log message
func (l *FileLogger) TaskSuccess(format string, isComplete bool, words ...interface{}) {
	l.printMessage(format, "", "success", true, isComplete, words...)
}

// Warn log message
func (l *FileLogger) Warn(format string, words ...interface{}) {
	l.printMessage(format, IconWarning, "warn", false, false, words...)
}

// TaskWarn log message
func (l *FileLogger) TaskWarn(format string, words ...interface{}) {
	l.printMessage(format, "", "warn", true, false, words...)
}

// Command log message
func (l *FileLogger) Command(format string, words ...interface{}) {
	l.printMessage(format, IconWrench, "command", false, false, words...)
}

// Disabled log message
func (l *FileLogger) Disabled(format string, words ...interface{}) {
	l.printMessage(format, IconBlackSquare, "disabled", false, false, words...)
}

// Notice log message
func (l *FileLogger) Notice(format string, words ...interface{}) {
	l.printMessage(format, IconFlag, "notice", false, false, words...)
}

// Debug log message
func (l *FileLogger) Debug(format string, words ...interface{}) {
	l.printMessage(format, IconFire, "debug", false, false, words...)
}

// Trace log message
func (l *FileLogger) Trace(format string, words ...interface{}) {
	l.printMessage(format, IconBulb, "trace", false, false, words...)
}

// Error log message
func (l *FileLogger) Error(format string, words ...interface{}) {
	l.printMessage(format, IconRevolvingLight, "error", false, false, words...)
}

// Error log message
func (l *FileLogger) Exception(err error, format string, words ...interface{}) {
	if format == "" {
		format = err.Error()
	} else {
		format = format + ", err " + err.Error()
	}
	l.printMessage(format, IconRevolvingLight, "error", false, false, words...)
}

// LogError log message
func (l *FileLogger) LogError(message error) {
	if message != nil {
		l.printMessage(message.Error(), IconRevolvingLight, "error", false, false)
	}
}

// TaskError log message
func (l *FileLogger) TaskError(format string, isComplete bool, words ...interface{}) {
	l.printMessage(format, "", "error", true, isComplete, l.useTimestamp)
}

// Fatal log message
func (l *FileLogger) Fatal(format string, words ...interface{}) {
	l.printMessage(format, IconRevolvingLight, "error", false, true, words...)
}

// FatalError log message
func (l *FileLogger) FatalError(e error, format string, words ...interface{}) {
	l.Error(format, words...)
	if e != nil {
		panic(e)
	}
}

// printMessage Prints a message in the system
func (l *FileLogger) printMessage(format string, icon LoggerIcon, level string, isTask bool, isComplete bool, words ...interface{}) {
	if !l.enabled {
		return
	}

	if !strings.HasSuffix(format, "\n") {
		format = format + "\n"
	}

	if l.userCorrelationId {
		correlationId := os.Getenv("CORRELATION_ID")
		if correlationId != "" {
			format = "[" + correlationId + "] " + "[" + strings.ToUpper(level) + "]" + format
		}
	}

	if l.useTimestamp {
		format = fmt.Sprintf("%s %s", time.Now().Format(time.RFC3339), format)
	}

	formattedWords := make([]interface{}, len(words))
	if len(words) > 0 {
		for i := range words {
			formattedWords[i] = fmt.Sprintf("%v", words[i])
		}
	}

	l.rotateLogFile()
	l.writer.Write([]byte(fmt.Sprintf(format, formattedWords...)))
}

func (l *FileLogger) Close() {
	if l.enabled {
		file, ok := l.writer.(*os.File)
		if ok {
			file.Close()
		}
	}
}

func (l *FileLogger) rotateLogFile() {
	if l.enabled {
		file, ok := l.writer.(*os.File)
		if ok {
			fileInfo, err := file.Stat()
			if err != nil {
				return
			}
			// Get the maximum log file size from the environment variable
			maxSizeStr := os.Getenv("MAX_LOG_FILE_SIZE")
			maxSize := int64(1024 * 1024 * 5) // Default to 5MB if not set
			if maxSizeStr != "" {
				if parsedSize, err := strconv.ParseInt(maxSizeStr, 10, 64); err == nil {
					maxSize = parsedSize
				}
			}

			// File is smaller than 5MB keep it
			if fileInfo.Size() < maxSize {
				return
			}

			// Delete the last file if it exists
			lastFile := fmt.Sprintf("%s.%02d", l.filename, 9)
			if _, err := os.Stat(lastFile); err == nil {
				os.Remove(lastFile)
			}

			for i := 9; i >= 1; i-- {
				oldPath := fmt.Sprintf("%s.%02d", l.filename, i)
				newPath := fmt.Sprintf("%s.%02d", l.filename, i+1)
				if _, err := os.Stat(oldPath); err == nil {
					if err := os.Rename(oldPath, newPath); err != nil {
						return
					}
				}
			}
			if err := os.Rename(l.filename, fmt.Sprintf("%s.01", l.filename)); err != nil {
				return
			}
			file.Close()
			file, err := os.OpenFile(l.filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
			if err != nil {
				panic(err)
			}
			l.writer = file
		}
	}
}
