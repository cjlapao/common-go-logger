package log

// Package log provides mocking capabilities for testing logging functionality.
// The mock logger captures log messages and provides methods to verify logging behavior.

import (
	"fmt"
	"io"
	"os"

	"github.com/cjlapao/common-go/strcolor"
)

// MockedLogMessage represents a captured log message for testing purposes.
// It contains the essential components of a log message without the timestamp.
type MockedLogMessage struct {
	Message string // The formatted log message
	Level   string // The log level (info, error, warn, etc.)
	Icon    string // The icon used in the message (if any)
}

// MockLogger implements the Logger interface for testing purposes.
// It captures log messages and provides methods to verify logging behavior.
//
// Example usage:
//
//	// Create and initialize a mock logger
//	mockLogger := &MockLogger{}
//	mockLogger = mockLogger.Init().(*MockLogger)
//
//	// Use the mock logger
//	mockLogger.Info("test message")
//
//	// Verify the last message
//	if mockLogger.LastPrintedMessage.Level != "info" {
//	    t.Errorf("Expected level 'info', got %s", mockLogger.LastPrintedMessage.Level)
//	}
//
//	// Verify all messages
//	for _, msg := range mockLogger.PrintedMessages {
//	    // Check message properties
//	}
type MockLogger struct {
	LastPrintedMessage MockedLogMessage   // The most recent message logged
	PrintedMessages    []MockedLogMessage // All messages logged
	LastCallType       string             // The type of the last logging call made
	useTimestamp       bool               // Whether timestamps are enabled
	userCorrelationId  bool               // Whether correlation IDs are enabled
	useIcons           bool               // Whether icons are enabled
	writer             io.Writer          // The output writer (usually stdout for testing)
}

// Init initializes a new MockLogger with default settings.
// Returns the initialized logger as a Logger interface.
//
// Example:
//
//	mockLogger := &MockLogger{}
//	logger := mockLogger.Init()
//	logger.Info("test message")
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

// Clear resets the mock logger's message history.
// This is useful between tests to ensure a clean state.
//
// Example:
//
//	mockLogger := &MockLogger{}
//	mockLogger.Info("first test")
//	mockLogger.Clear() // Clear history
//	mockLogger.Info("second test")
//	// Only "second test" will be in PrintedMessages
func (l *MockLogger) Clear() {
	l.LastPrintedMessage = MockedLogMessage{}
	l.PrintedMessages = []MockedLogMessage{}
}

// IsTimestampEnabled returns whether timestamp logging is enabled.
//
// Example:
//
//	mockLogger := &MockLogger{}
//	mockLogger.UseTimestamp(true)
//	if mockLogger.IsTimestampEnabled() {
//	    // Timestamps are enabled
//	}
func (l *MockLogger) IsTimestampEnabled() bool {
	return l.useTimestamp
}

// UseTimestamp enables or disables timestamp logging.
//
// Example:
//
//	mockLogger := &MockLogger{}
//	mockLogger.UseTimestamp(true)  // Enable timestamps
//	mockLogger.Info("With timestamp")
//	mockLogger.UseTimestamp(false) // Disable timestamps
//	mockLogger.Info("Without timestamp")
func (l *MockLogger) UseTimestamp(value bool) {
	l.useTimestamp = value
}

// UseCorrelationId enables or disables correlation ID logging.
//
// Example:
//
//	mockLogger := &MockLogger{}
//	mockLogger.UseCorrelationId(true)
//	// Now correlation IDs will be included in log messages
//	// when CORRELATION_ID environment variable is set
func (l *MockLogger) UseCorrelationId(value bool) {
	l.userCorrelationId = value
}

// UseIcons enables or disables icon display in log messages.
//
// Example:
//
//	mockLogger := &MockLogger{}
//	mockLogger.UseIcons(true)
//	mockLogger.Info("With icon")    // Will include ℹ
//	mockLogger.UseIcons(false)
//	mockLogger.Info("Without icon") // No icon
func (l *MockLogger) UseIcons(value bool) {
	l.useIcons = value
}

// Log records a message with the specified level.
//
// Example:
//
//	mockLogger := &MockLogger{}
//	mockLogger.Log("Processing %d items", Info, 42)
//	// Verify:
//	if mockLogger.LastPrintedMessage.Level != "info" {
//	    t.Error("Wrong log level")
//	}
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

// LogIcon records a message with a custom icon and level.
//
// Example:
//
//	mockLogger := &MockLogger{}
//	mockLogger.LogIcon("🌟", "Special event %s", Info, "occurred")
//	// Verify:
//	if mockLogger.LastPrintedMessage.Icon != "🌟" {
//	    t.Error("Wrong icon")
//	}
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

// LogHighlight records a message with highlighted text.
//
// Example:
//
//	mockLogger := &MockLogger{}
//	mockLogger.LogHighlight("Status: %s", Info, strcolor.Green, "SUCCESS")
//	// Verify the message contains the highlighted text
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

// Info records an informational message.
//
// Example:
//
//	mockLogger := &MockLogger{}
//	mockLogger.Info("Server started on port %d", 8080)
//	// Verify:
//	if mockLogger.LastPrintedMessage.Level != "info" {
//	    t.Error("Wrong log level")
//	}
func (l *MockLogger) Info(format string, words ...interface{}) {
	l.printMessage(format, IconInfo, "info", false, false, words...)
}

// Success records a success message.
//
// Example:
//
//	mockLogger := &MockLogger{}
//	mockLogger.Success("Operation %s completed", "backup")
//	// Verify:
//	if mockLogger.LastPrintedMessage.Icon != string(IconThumbsUp) {
//	    t.Error("Wrong icon")
//	}
func (l *MockLogger) Success(format string, words ...interface{}) {
	l.printMessage(format, IconThumbsUp, "success", false, false, words...)
}

// TaskSuccess records a task completion message.
//
// Example:
//
//	mockLogger := &MockLogger{}
//	mockLogger.TaskSuccess("Task %s", true, "backup")
//	// Verify:
//	if !strings.Contains(mockLogger.LastPrintedMessage.Message, "backup") {
//	    t.Error("Message not recorded correctly")
//	}
func (l *MockLogger) TaskSuccess(format string, isComplete bool, words ...interface{}) {
	l.printMessage(format, "", "success", true, isComplete, words...)
}

// Warn records a warning message.
//
// Example:
//
//	mockLogger := &MockLogger{}
//	mockLogger.Warn("Disk usage at %d%%", 90)
//	// Verify:
//	if mockLogger.LastPrintedMessage.Level != "warn" {
//	    t.Error("Wrong log level")
//	}
func (l *MockLogger) Warn(format string, words ...interface{}) {
	l.printMessage(format, IconWarning, "warn", false, false, words...)
}

// TaskWarn records a task warning message.
//
// Example:
//
//	mockLogger := &MockLogger{}
//	mockLogger.TaskWarn("Task %s warning", "backup")
//	// Verify:
//	if mockLogger.LastPrintedMessage.Level != "warn" {
//	    t.Error("Wrong log level")
//	}
func (l *MockLogger) TaskWarn(format string, words ...interface{}) {
	l.printMessage(format, "", "warn", true, false, words...)
}

// Command records a command execution message.
//
// Example:
//
//	mockLogger := &MockLogger{}
//	mockLogger.Command("Executing: %s", "git pull")
//	// Verify:
//	if mockLogger.LastPrintedMessage.Icon != string(IconWrench) {
//	    t.Error("Wrong icon")
//	}
func (l *MockLogger) Command(format string, words ...interface{}) {
	l.printMessage(format, IconWrench, "command", false, false, words...)
}

// Disabled records a disabled feature message.
//
// Example:
//
//	mockLogger := &MockLogger{}
//	mockLogger.Disabled("Feature %s is disabled", "beta")
//	// Verify:
//	if mockLogger.LastPrintedMessage.Icon != string(IconBlackSquare) {
//	    t.Error("Wrong icon")
//	}
func (l *MockLogger) Disabled(format string, words ...interface{}) {
	l.printMessage(format, IconBlackSquare, "disabled", false, false, words...)
}

// Notice records a notice message.
//
// Example:
//
//	mockLogger := &MockLogger{}
//	mockLogger.Notice("Maintenance in %d minutes", 5)
//	// Verify:
//	if mockLogger.LastPrintedMessage.Icon != string(IconFlag) {
//	    t.Error("Wrong icon")
//	}
func (l *MockLogger) Notice(format string, words ...interface{}) {
	l.printMessage(format, IconFlag, "notice", false, false, words...)
}

// Debug records a debug message.
//
// Example:
//
//	mockLogger := &MockLogger{}
//	mockLogger.Debug("Variable x = %v", someVar)
//	// Verify:
//	if mockLogger.LastPrintedMessage.Level != "debug" {
//	    t.Error("Wrong log level")
//	}
func (l *MockLogger) Debug(format string, words ...interface{}) {
	l.printMessage(format, IconFire, "debug", false, false, words...)
}

// Trace records a trace message.
//
// Example:
//
//	mockLogger := &MockLogger{}
//	mockLogger.Trace("Entering function %s", "processItem")
//	// Verify:
//	if mockLogger.LastPrintedMessage.Level != "trace" {
//	    t.Error("Wrong log level")
//	}
func (l *MockLogger) Trace(format string, words ...interface{}) {
	l.printMessage(format, IconBulb, "trace", false, false, words...)
}

// Error records an error message.
//
// Example:
//
//	mockLogger := &MockLogger{}
//	mockLogger.Error("Failed to connect: %s", "timeout")
//	// Verify:
//	if mockLogger.LastPrintedMessage.Level != "error" {
//	    t.Error("Wrong log level")
//	}
func (l *MockLogger) Error(format string, words ...interface{}) {
	l.printMessage(format, IconRevolvingLight, "error", false, false, words...)
}

// Exception records an error with additional context.
//
// Example:
//
//	mockLogger := &MockLogger{}
//	err := errors.New("not found")
//	mockLogger.Exception(err, "Failed to load %s", "config.json")
//	// Verify:
//	if !strings.Contains(mockLogger.LastPrintedMessage.Message, "not found") {
//	    t.Error("Error message not included")
//	}
func (l *MockLogger) Exception(err error, format string, words ...interface{}) {
	if format == "" {
		format = err.Error()
	} else {
		format = format + ", err " + err.Error()
	}
	l.printMessage(format, IconRevolvingLight, "error", false, false, words...)
}

// LogError records an error object directly.
//
// Example:
//
//	mockLogger := &MockLogger{}
//	err := errors.New("connection failed")
//	mockLogger.LogError(err)
//	// Verify:
//	if mockLogger.LastPrintedMessage.Message != "connection failed" {
//	    t.Error("Wrong error message")
//	}
func (l *MockLogger) LogError(message error) {
	if message != nil {
		l.printMessage(message.Error(), IconRevolvingLight, "error", false, false)
	}
}

// TaskError records a task error message.
//
// Example:
//
//	mockLogger := &MockLogger{}
//	mockLogger.TaskError("Task %s failed", true, "backup")
//	// Verify:
//	if mockLogger.LastPrintedMessage.Level != "error" {
//	    t.Error("Wrong log level")
//	}
func (l *MockLogger) TaskError(format string, isComplete bool, words ...interface{}) {
	l.printMessage(format, "", "error", true, isComplete, l.useTimestamp)
}

// Fatal records a fatal error message.
//
// Example:
//
//	mockLogger := &MockLogger{}
//	mockLogger.Fatal("System failure: %s", "out of memory")
//	// Verify:
//	if mockLogger.LastPrintedMessage.Level != "error" {
//	    t.Error("Wrong log level")
//	}
func (l *MockLogger) Fatal(format string, words ...interface{}) {
	l.printMessage(format, IconRevolvingLight, "error", false, true, words...)
}

// FatalError records an error message and panics if the error is not nil.
//
// Example:
//
//	mockLogger := &MockLogger{}
//	defer func() {
//	    if r := recover(); r == nil {
//	        t.Error("Expected panic")
//	    }
//	}()
//	err := errors.New("critical failure")
//	mockLogger.FatalError(err, "System crashed")
func (l *MockLogger) FatalError(e error, format string, words ...interface{}) {
	l.Error(format, words...)
	if e != nil {
		panic(e)
	}
}

// printMessage captures a log message for testing purposes.
// This internal method is used by all logging methods to record messages.
//
// Parameters:
//   - format: The message format string
//   - icon: The icon to be displayed (if icons are enabled)
//   - level: The log level (info, error, warn, etc.)
//   - isTask: Whether this is a task-related message
//   - isComplete: Whether this is a completion message
//   - words: Format string arguments
//
// Example usage (internal):
//
//	l.printMessage("Processing %s", IconInfo, "info", false, false, "data")
func (l *MockLogger) printMessage(format string, icon LoggerIcon, level string, isTask bool, isComplete bool, words ...interface{}) {
	l.LastPrintedMessage = MockedLogMessage{Message: fmt.Sprintf(format, words...), Level: level, Icon: string(icon)}
	l.PrintedMessages = append(l.PrintedMessages, l.LastPrintedMessage)
}
