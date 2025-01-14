package log

import (
	"fmt"
	"net/http"
)

// AddCmdLogger adds a command line logger to the LoggerService.
// The command line logger writes formatted log messages to stdout.
// It inherits timestamp, correlation ID, and icon settings from the LoggerService.
//
// Example:
//
//	service := log.New()
//	service.WithTimestamp().WithIcons()
//	service.AddCmdLogger()
//	service.Info("Hello from command line!")
//	// Output: [2024-03-20T10:00:00Z] ℹ info: Hello from command line!
func (l *LoggerService) AddCmdLogger() {
	Register(&CmdLogger{
		useTimestamp:      l.UseTimestamp,
		userCorrelationId: l.useCorrelationId,
		useIcons:          l.useIcons,
	})
}

// AddFileLogger adds a file logger to the LoggerService.
// The file logger writes formatted log messages to the specified file.
// It inherits timestamp, correlation ID, and icon settings from the LoggerService.
//
// Example:
//
//	service := log.New()
//	service.WithTimestamp()
//	service.AddFileLogger("app.log")
//	service.Info("Hello from file logger!")
//	// Content of app.log: [2024-03-20T10:00:00Z] info: Hello from file logger!
func (l *LoggerService) AddFileLogger(filename string) {
	Register(&FileLogger{
		userCorrelationId: l.useCorrelationId,
		useIcons:          l.useIcons,
		useTimestamp:      l.UseTimestamp,
		filename:          filename,
	})
}

// AddChannelLogger adds a channel-based logger to the LoggerService.
// The channel logger sends log messages through a channel, allowing for
// asynchronous processing of log messages via OnMessage subscribers.
// It inherits timestamp, correlation ID, and icon settings from the LoggerService.
//
// Example:
//
//	service := log.New()
//	service.AddChannelLogger()
//	service.OnMessage(func(msg LogMessage) {
//	    fmt.Printf("Received: %s\n", msg)
//	})
//	service.Info("Hello from channel!")
func (l *LoggerService) AddChannelLogger() {
	channelLogger := &ChannelLogger{
		useTimestamp:      l.UseTimestamp,
		userCorrelationId: l.useCorrelationId,
		useIcons:          l.useIcons,
	}
	Register(channelLogger)
}

// WithDebug sets the log level to Debug, enabling all log messages
// at Debug level and above (Debug, Info, Warning, Error).
//
// Example:
//
//	service := log.New()
//	service.WithDebug()
//	service.Debug("This will be logged")
//	service.Trace("This won't be logged")
func (l *LoggerService) WithDebug() *LoggerService {
	l.LogLevel = Debug
	return l
}

// WithTrace sets the log level to Trace, enabling all log messages
// at all levels (Trace, Debug, Info, Warning, Error).
//
// Example:
//
//	service := log.New()
//	service.WithTrace()
//	service.Debug("This will be logged")
//	service.Trace("This will also be logged")
func (l *LoggerService) WithTrace() *LoggerService {
	l.LogLevel = Trace
	return l
}

// WithWarning sets the log level to Warning, enabling only Warning
// and Error level messages.
//
// Example:
//
//	service := log.New()
//	service.WithWarning()
//	service.Info("This won't be logged")
//	service.Warn("This will be logged")
//	service.Error("This will be logged")
func (l *LoggerService) WithWarning() *LoggerService {
	l.LogLevel = Warning
	return l
}

// WithTimestamp enables timestamp prefixing for all log messages.
// Returns the LoggerService for method chaining.
//
// Example:
//
//	service := log.New()
//	service.WithTimestamp()
//	service.Info("Hello")
//	// Output: [2024-03-20T10:00:00Z] info: Hello
func (l *LoggerService) WithTimestamp() *LoggerService {
	for _, logger := range l.Loggers {
		logger.UseTimestamp(true)
	}

	l.UseTimestamp = true
	return l
}

// ToggleTimestamp toggles the timestamp display on/off for all loggers.
// Returns the LoggerService for method chaining.
//
// Example:
//
//	service := log.New()
//	service.WithTimestamp()  // Enable timestamps
//	service.Info("With timestamp")
//	service.ToggleTimestamp() // Disable timestamps
//	service.Info("Without timestamp")
func (l *LoggerService) ToggleTimestamp() *LoggerService {
	l.UseTimestamp = !l.UseTimestamp

	for _, logger := range l.Loggers {
		logger.UseTimestamp(l.UseTimestamp)
	}

	return l
}

// EnableTimestamp explicitly sets the timestamp display state for all loggers.
// Returns the LoggerService for method chaining.
//
// Example:
//
//	service := log.New()
//	service.EnableTimestamp(true)
//	service.Info("With timestamp")
//	service.EnableTimestamp(false)
//	service.Info("Without timestamp")
func (l *LoggerService) EnableTimestamp(value bool) *LoggerService {
	for _, logger := range l.Loggers {
		logger.UseTimestamp(value)
	}

	l.UseTimestamp = value

	return l
}

// WithCorrelationId enables correlation ID display in log messages.
// Correlation IDs help track related log messages across different parts of the system.
// Returns the LoggerService for method chaining.
//
// Example:
//
//	service := log.New()
//	service.WithCorrelationId()
//	os.Setenv("CORRELATION_ID", "req-123")
//	service.Info("Processing request")
//	// Output: [req-123] info: Processing request
func (l *LoggerService) WithCorrelationId() *LoggerService {
	l.useCorrelationId = true
	for _, logger := range l.Loggers {
		logger.UseCorrelationId(true)
	}
	return l
}

// WithIcons enables icon display in log messages.
// Icons provide visual indicators for different types of log messages.
// Returns the LoggerService for method chaining.
//
// Example:
//
//	service := log.New()
//	service.WithIcons()
//	service.Info("Information")    // Output: ℹ info: Information
//	service.Error("Problem")       // Output: 🚨 error: Problem
//	service.Success("Complete")    // Output: 👍 success: Complete
func (l *LoggerService) WithIcons() *LoggerService {
	l.useIcons = true
	for _, logger := range l.Loggers {
		logger.UseIcons(true)
	}
	return l
}

// Log logs a message with the specified level and format.
// This is a low-level logging function that allows direct control of the log level.
//
// Example:
//
//	service := log.New()
//	service.Log("Processing item %d", log.Info, 42)
//	// Output: info: Processing item 42
func (l *LoggerService) Log(format string, level Level, words ...interface{}) {
	for _, logger := range l.Loggers {
		logger.Log(format, level, words...)
	}
}

// LogIcon logs a message with a custom icon and specified level.
// This is a low-level logging function that allows custom icons.
//
// Example:
//
//	service := log.New()
//	service.LogIcon("🌟", "Special event %s", log.Info, "occurred")
//	// Output: 🌟 info: Special event occurred
func (l *LoggerService) LogIcon(icon LoggerIcon, format string, level Level, words ...interface{}) {
	for _, logger := range l.Loggers {
		logger.LogIcon(icon, format, level, words...)
	}
}

// LogHighlight logs a message with highlighted words using the specified color.
// The color is applied to the interpolated values, not the format string.
// This is useful for emphasizing important parts of log messages.
//
// Example:
//
//	service := log.New()
//	service.HighlightColor = strcolor.Green
//	service.LogHighlight("Status: %s, Count: %d", log.Info, "ACTIVE", 42)
//	// Output: info: Status: ACTIVE (in green), Count: 42 (in green)
//
//	// With different color:
//	service.HighlightColor = strcolor.Red
//	service.LogHighlight("Warning: %s", log.Warning, "Critical state")
//	// Output: warn: Warning: Critical state (in red)
func (l *LoggerService) LogHighlight(format string, level Level, words ...interface{}) {
	for _, logger := range l.Loggers {
		logger.LogHighlight(format, level, l.HighlightColor, words...)
	}
}

// Info logs an informational message.
// Messages are only logged if the service's log level is Info or higher.
//
// Example:
//
//	service := log.New()
//	service.Info("Server started on port %d", 8080)
//	// Output: info: Server started on port 8080
func (l *LoggerService) Info(format string, words ...interface{}) {
	if l.LogLevel >= Info {
		for _, logger := range l.Loggers {
			logger.Info(format, words...)
		}
	}
}

// Success logs a success message with a thumbs-up icon.
// Messages are only logged if the service's log level is Info or higher.
//
// Example:
//
//	service := log.New().WithIcons()
//	service.Success("Operation completed: %s", "backup")
//	// Output: 👍 success: Operation completed: backup
func (l *LoggerService) Success(format string, words ...interface{}) {
	if l.LogLevel >= Info {
		for _, logger := range l.Loggers {
			logger.Success(format, words...)
		}
	}
}

// Warn logs a warning message with a warning icon.
// Messages are only logged if the service's log level is Warning or higher.
//
// Example:
//
//	service := log.New().WithIcons()
//	service.Warn("Disk usage high: %d%%", 90)
//	// Output: ⚠ warn: Disk usage high: 90%
func (l *LoggerService) Warn(format string, words ...interface{}) {
	if l.LogLevel >= Warning {
		for _, logger := range l.Loggers {
			logger.Warn(format, words...)
		}
	}
}

// Command logs a command execution with a wrench icon.
// Messages are only logged if the service's log level is Info or higher.
//
// Example:
//
//	service := log.New().WithIcons()
//	service.Command("Executing: %s", "git pull")
//	// Output: 🔧 command: Executing: git pull
func (l *LoggerService) Command(format string, words ...interface{}) {
	if l.LogLevel >= Info {
		for _, logger := range l.Loggers {
			logger.Command(format, words...)
		}
	}
}

// Disabled logs a disabled feature message with a black square icon.
// Messages are only logged if the service's log level is Info or higher.
//
// Example:
//
//	service := log.New().WithIcons()
//	service.Disabled("Feature %s is disabled", "beta-testing")
//	// Output: ⬛ disabled: Feature beta-testing is disabled
func (l *LoggerService) Disabled(format string, words ...interface{}) {
	if l.LogLevel >= Info {
		for _, logger := range l.Loggers {
			logger.Disabled(format, words...)
		}
	}
}

// Notice logs a notice message with a flag icon.
// Messages are only logged if the service's log level is Info or higher.
//
// Example:
//
//	service := log.New().WithIcons()
//	service.Notice("Maintenance scheduled for %s", "tomorrow")
//	// Output: 🚩 notice: Maintenance scheduled for tomorrow
func (l *LoggerService) Notice(format string, words ...interface{}) {
	if l.LogLevel >= Info {
		for _, logger := range l.Loggers {
			logger.Notice(format, words...)
		}
	}
}

// Debug logs a debug message with a fire icon.
// Messages are only logged if the service's log level is Debug or higher.
//
// Example:
//
//	service := log.New().WithDebug().WithIcons()
//	service.Debug("Variable x = %d", 42)
//	// Output: 🔥 debug: Variable x = 42
func (l *LoggerService) Debug(format string, words ...interface{}) {
	if l.LogLevel >= Debug {
		for _, logger := range l.Loggers {
			logger.Debug(format, words...)
		}
	}
}

// Trace logs a trace message with a bulb icon.
// Messages are only logged if the service's log level is Trace.
// This is the most detailed level of logging, useful for debugging and development.
//
// Example:
//
//	service := log.New().WithTrace().WithIcons()
//	service.Trace("Entering function %s with args: %v", "processItem", []string{"a", "b"})
//	// Output: 💡 trace: Entering function processItem with args: [a b]
//
//	// With timestamp enabled:
//	service.WithTimestamp()
//	service.Trace("Variable state: %+v", myVar)
//	// Output: [2024-03-20T10:00:00Z] 💡 trace: Variable state: {Field:value}
func (l *LoggerService) Trace(format string, words ...interface{}) {
	if l.LogLevel >= Trace {
		for _, logger := range l.Loggers {
			logger.Debug(format, words...)
		}
	}
}

// Error logs an error message with a revolving light icon.
// Messages are only logged if the service's log level is Error or higher.
//
// Example:
//
//	service := log.New().WithIcons()
//	service.Error("Failed to connect: %s", "timeout")
//	// Output: 🚨 error: Failed to connect: timeout
func (l *LoggerService) Error(format string, words ...interface{}) {
	if l.LogLevel >= Error {
		for _, logger := range l.Loggers {
			logger.Error(format, words...)
		}
	}
}

// LogError logs an error object directly.
// Messages are only logged if the service's log level is Error or higher.
//
// Example:
//
//	service := log.New()
//	err := errors.New("connection failed")
//	service.LogError(err)
//	// Output: error: connection failed
func (l *LoggerService) LogError(message error) {
	if l.LogLevel >= Error {
		if message != nil {
			for _, logger := range l.Loggers {
				logger.Error(message.Error())
			}
		}
	}
}

// Exception logs an error with additional context information.
// Messages are only logged if the service's log level is Error or higher.
//
// Example:
//
//	service := log.New()
//	err := errors.New("not found")
//	service.Exception(err, "Failed to load config from %s", "config.json")
//	// Output: error: Failed to load config from config.json, err not found
func (l *LoggerService) Exception(err error, format string, words ...interface{}) {
	if l.LogLevel >= Error {
		for _, logger := range l.Loggers {
			logger.Exception(err, format, words...)
		}
	}
}

// Fatal logs a fatal error message with a revolving light icon.
// Messages are only logged if the service's log level is Error or higher.
//
// Example:
//
//	service := log.New().WithIcons()
//	service.Fatal("System failure: %s", "out of memory")
//	// Output: 🚨 error: System failure: out of memory
func (l *LoggerService) Fatal(format string, words ...interface{}) {
	if l.LogLevel >= Error {
		for _, logger := range l.Loggers {
			logger.Fatal(format, words...)
		}
	}
}

// FatalError logs an error message and then panics if the error is not nil.
// This should be used for unrecoverable errors that require immediate shutdown.
//
// Example:
//
//	service := log.New()
//	err := errors.New("critical failure")
//	// This will log the error and then panic:
//	service.FatalError(err, "System crashed: %s", "unrecoverable state")
func (l *LoggerService) FatalError(e error, format string, words ...interface{}) {
	for _, logger := range l.Loggers {
		logger.Error(format, words...)
	}

	if e != nil {
		panic(e)
	}
}

// GetRequestPrefix generates a prefix for HTTP request logging.
// It includes the request ID if present in X-Request-Id header and optionally
// includes the HTTP method and path. This is useful for consistent request logging
// across your application.
//
// Parameters:
//   - r: The HTTP request to generate the prefix for
//   - logUrl: If true, includes the HTTP method and path in the prefix
//
// Example:
//
//	service := log.New()
//
//	// Create a request with ID
//	req, _ := http.NewRequest("GET", "/api/users", nil)
//	req.Header.Set("X-Request-Id", "req-123")
//
//	// Get prefix with URL
//	prefix := service.GetRequestPrefix(req, true)
//	service.Info(prefix + "Processing request")
//	// Output: [req-123] [GET] [/api/users] Processing request
//
//	// Get prefix without URL
//	prefix = service.GetRequestPrefix(req, false)
//	service.Info(prefix + "Request received")
//	// Output: [req-123] Request received
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

// OnMessage registers a callback function to receive log messages from the channel logger.
// The callback will be executed asynchronously for each log message.
// Returns a subscription ID that can be used to unsubscribe later.
// If no channel logger is configured, returns an empty string.
//
// Example:
//
//	service := log.New()
//	service.AddChannelLogger()
//
//	// Subscribe with custom ID
//	subID := service.OnMessage("my-handler", func(msg LogMessage) {
//	    fmt.Printf("Received [%s]: %s\n", msg.Level, msg.Message)
//	})
//
//	// Log some messages
//	service.Info("Test message")
//	service.Error("Something went wrong")
//
//	// Later, unsubscribe:
//	service.RemoveMessageHandler(subID)
func (l *LoggerService) OnMessage(id string, callback func(LogMessage)) string {
	// Find the channel logger instance
	var channelLogger *ChannelLogger
	for _, logger := range l.Loggers {
		if cl, ok := logger.(*ChannelLogger); ok {
			channelLogger = cl
			break
		}
	}

	if channelLogger == nil {
		return ""
	}

	// Subscribe with a filter that accepts all messages
	subID, ch := channelLogger.Subscribe(id, func(LogMessage) bool { return true })

	// Start goroutine to process messages
	go func() {
		for msg := range ch {
			callback(msg)
		}
	}()

	return subID
}

// RemoveMessageHandler unsubscribes a message handler using its subscription ID.
// Returns true if the handler was successfully removed, false if the handler wasn't found
// or if no channel logger is configured.
//
// Example:
//
//	service := log.New()
//	service.AddChannelLogger()
//
//	// Subscribe to messages
//	subID := service.OnMessage("my-handler", func(msg LogMessage) {
//	    fmt.Printf("Got message: %s\n", msg)
//	})
//
//	// Remove the subscription
//	success := service.RemoveMessageHandler(subID)
//	if !success {
//	    fmt.Println("Failed to remove message handler")
//	}
func (l *LoggerService) RemoveMessageHandler(subscriptionID string) bool {
	for _, logger := range l.Loggers {
		if cl, ok := logger.(*ChannelLogger); ok {
			return cl.Unsubscribe(subscriptionID)
		}
	}
	return false
}
