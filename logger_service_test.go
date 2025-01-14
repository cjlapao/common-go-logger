package log

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoggerService_Configuration(t *testing.T) {
	t.Run("log level configuration", func(t *testing.T) {
		service := &LoggerService{}

		service.WithDebug()
		assert.Equal(t, Debug, service.LogLevel)

		service.WithTrace()
		assert.Equal(t, Trace, service.LogLevel)

		service.WithWarning()
		assert.Equal(t, Warning, service.LogLevel)
	})

	t.Run("timestamp configuration", func(t *testing.T) {
		service := &LoggerService{}
		mockLogger := &MockLogger{}
		service.Loggers = append(service.Loggers, mockLogger)

		service.WithTimestamp()
		assert.True(t, service.UseTimestamp)
		assert.True(t, mockLogger.useTimestamp)

		service.ToggleTimestamp()
		assert.False(t, service.UseTimestamp)
		assert.False(t, mockLogger.useTimestamp)

		service.EnableTimestamp(true)
		assert.True(t, service.UseTimestamp)
		assert.True(t, mockLogger.useTimestamp)
	})

	t.Run("correlation ID configuration", func(t *testing.T) {
		service := &LoggerService{}
		mockLogger := &MockLogger{}
		service.Loggers = append(service.Loggers, mockLogger)

		service.WithCorrelationId()
		assert.True(t, service.useCorrelationId)
		assert.True(t, mockLogger.userCorrelationId)
	})

	t.Run("icons configuration", func(t *testing.T) {
		service := &LoggerService{}
		mockLogger := &MockLogger{}
		service.Loggers = append(service.Loggers, mockLogger)

		service.WithIcons()
		assert.True(t, service.useIcons)
		assert.True(t, mockLogger.useIcons)
	})
}

func TestLoggerService_LoggingMethods(t *testing.T) {
	tests := []struct {
		name      string
		logLevel  Level
		logFunc   func(*LoggerService)
		message   string
		shouldLog bool
	}{
		{
			name:      "Info with InfoLevel",
			logLevel:  Info,
			logFunc:   func(s *LoggerService) { s.Info("test info") },
			message:   "test info",
			shouldLog: true,
		},
		{
			name:      "Debug with DebugLevel",
			logLevel:  Debug,
			logFunc:   func(s *LoggerService) { s.Debug("test debug") },
			message:   "test debug",
			shouldLog: true,
		},
		{
			name:      "Debug with InfoLevel",
			logLevel:  Info,
			logFunc:   func(s *LoggerService) { s.Debug("test debug") },
			message:   "test debug",
			shouldLog: false,
		},
		{
			name:      "Error with any level",
			logLevel:  Info,
			logFunc:   func(s *LoggerService) { s.Error("test error") },
			message:   "test error",
			shouldLog: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := &MockLogger{}
			service := &LoggerService{
				LogLevel: tt.logLevel,
				Loggers:  []Logger{mockLogger},
			}

			tt.logFunc(service)

			if tt.shouldLog {
				assert.Contains(t, mockLogger.PrintedMessages[0].Message, tt.message)
			} else {
				assert.NotContains(t, mockLogger.PrintedMessages, tt.message)
			}
		})
	}
}

func TestLoggerService_ErrorHandling(t *testing.T) {
	mockLogger := &MockLogger{}
	service := &LoggerService{
		LogLevel: Error,
		Loggers:  []Logger{mockLogger},
	}

	t.Run("LogError", func(t *testing.T) {
		mockLogger.PrintedMessages = make([]MockedLogMessage, 0)
		err := errors.New("test error")
		service.LogError(err)
		assert.Contains(t, mockLogger.PrintedMessages[0].Message, "test error")
	})

	t.Run("Exception", func(t *testing.T) {
		mockLogger.PrintedMessages = make([]MockedLogMessage, 0)
		err := errors.New("test exception")
		service.Exception(err, "error occurred")
		assert.Contains(t, mockLogger.PrintedMessages[0].Message, "error occurred")
	})
}

func TestLoggerService_GetRequestPrefix(t *testing.T) {
	service := &LoggerService{}

	tests := []struct {
		name     string
		request  *http.Request
		logURL   bool
		expected string
	}{
		{
			name: "with request ID and URL logging",
			request: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				req.Header.Set("X-Request-Id", "123")
				return req
			}(),
			logURL:   true,
			expected: "[123] [GET] [/test] ",
		},
		{
			name: "with request ID without URL logging",
			request: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				req.Header.Set("X-Request-Id", "123")
				return req
			}(),
			logURL:   false,
			expected: "[123] ",
		},
		{
			name:     "without request ID with URL logging",
			request:  httptest.NewRequest("GET", "/test", nil),
			logURL:   true,
			expected: "[GET] [/test] ",
		},
		{
			name:     "without request ID without URL logging",
			request:  httptest.NewRequest("GET", "/test", nil),
			logURL:   false,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.GetRequestPrefix(tt.request, tt.logURL)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLoggerService_AddLoggers(t *testing.T) {
	service := New()

	t.Run("AddCmdLogger", func(t *testing.T) {
		initialCount := len(service.Loggers)
		service.AddCmdLogger()
		assert.Equal(t, initialCount, len(service.Loggers))
	})

	t.Run("AddFileLogger", func(t *testing.T) {
		initialCount := len(service.Loggers)
		service.AddFileLogger("test.log")
		assert.Equal(t, initialCount+1, len(service.Loggers))
	})
}

func TestLoggerService_FatalError(t *testing.T) {
	// Setup
	mockLogger := &MockLogger{}
	service := &LoggerService{
		LogLevel: Error,
		Loggers:  []Logger{mockLogger},
	}

	// Test cases
	tests := []struct {
		name        string
		err         error
		format      string
		args        []interface{}
		shouldPanic bool
	}{
		{
			name:        "should log error and panic when error is not nil",
			err:         fmt.Errorf("test error"),
			format:      "Fatal error occurred: %s",
			args:        []interface{}{"test"},
			shouldPanic: true,
		},
		{
			name:        "should only log when error is nil",
			err:         nil,
			format:      "Message without error: %s",
			args:        []interface{}{"test"},
			shouldPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != tt.shouldPanic {
					t.Errorf("FatalError() panic = %v, shouldPanic = %v", r != nil, tt.shouldPanic)
				}
				if tt.shouldPanic && r != tt.err {
					t.Errorf("FatalError() panic with %v, want %v", r, tt.err)
				}
			}()

			service.FatalError(tt.err, tt.format, tt.args...)
		})
	}
}

func TestLoggerService_Log(t *testing.T) {
	// Setup
	service := &LoggerService{
		Loggers: []Logger{},
	}

	// Create a mock logger to verify the call
	mockLogger := &MockLogger{}
	service.Loggers = append(service.Loggers, mockLogger)

	// Test case
	testFormat := "test message %s"
	testLevel := Info
	testWord := "hello"

	// Execute
	service.Log(testFormat, testLevel, testWord)

	// Verify that the mock logger received the correct parameters
	if mockLogger.LastPrintedMessage.Message != fmt.Sprintf(testFormat, testWord) {
		t.Errorf("Expected format %s, got %s", testFormat, mockLogger.LastPrintedMessage.Message)
	}
	if mockLogger.LastPrintedMessage.Level != testLevel.String() {
		t.Errorf("Expected level %v, got %v", testLevel, mockLogger.LastPrintedMessage.Level)
	}
}

func TestLoggerService_LogIcon(t *testing.T) {
	// Setup
	service := &LoggerService{
		Loggers: []Logger{},
	}

	// Create a mock logger to verify the call
	mockLogger := &MockLogger{}
	service.Loggers = append(service.Loggers, mockLogger)

	// Test case
	testIcon := IconInfo
	testFormat := "test message %s"
	testLevel := Info
	testWord := "hello"

	// Execute
	service.LogIcon(testIcon, testFormat, testLevel, testWord)

	// Verify that the mock logger received the correct parameters
	if mockLogger.LastPrintedMessage.Message != fmt.Sprintf(testFormat, testWord) {
		t.Errorf("Expected format %s, got %s", testFormat, mockLogger.LastPrintedMessage.Message)
	}
	if mockLogger.LastPrintedMessage.Level != testLevel.String() {
		t.Errorf("Expected level %v, got %v", testLevel, mockLogger.LastPrintedMessage.Level)
	}
	if mockLogger.LastPrintedMessage.Icon != string(testIcon) {
		t.Errorf("Expected icon %v, got %v", testIcon, mockLogger.LastPrintedMessage.Icon)
	}
}

func TestLoggerService_Success(t *testing.T) {
	tests := []struct {
		name      string
		logLevel  Level
		format    string
		args      []interface{}
		shouldLog bool
	}{
		{
			name:      "should log when level is Info",
			logLevel:  Info,
			format:    "success message %s",
			args:      []interface{}{"test"},
			shouldLog: true,
		},
		{
			name:      "should log when level is Debug",
			logLevel:  Debug,
			format:    "success message %s",
			args:      []interface{}{"test"},
			shouldLog: true,
		},
		{
			name:      "should not log when level is Error",
			logLevel:  Error,
			format:    "success message %s",
			args:      []interface{}{"test"},
			shouldLog: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockLogger := &MockLogger{}
			service := &LoggerService{
				LogLevel: tt.logLevel,
				Loggers:  []Logger{mockLogger},
			}

			// Execute
			service.Success(tt.format, tt.args...)

			// Verify
			expectedMsg := fmt.Sprintf(tt.format, tt.args...)
			if tt.shouldLog {
				assert.Equal(t, expectedMsg, mockLogger.LastPrintedMessage.Message)
				assert.Equal(t, "success", mockLogger.LastPrintedMessage.Level)
				assert.Equal(t, string(IconThumbsUp), mockLogger.LastPrintedMessage.Icon)
			} else {
				assert.Empty(t, mockLogger.LastPrintedMessage.Message)
			}
		})
	}
}

func TestLoggerService_Warn(t *testing.T) {
	tests := []struct {
		name      string
		logLevel  Level
		format    string
		args      []interface{}
		shouldLog bool
	}{
		{
			name:      "should log when level is Warning",
			logLevel:  Warning,
			format:    "warning message %s",
			args:      []interface{}{"test"},
			shouldLog: true,
		},
		{
			name:      "should log when level is Debug",
			logLevel:  Debug,
			format:    "warning message %s",
			args:      []interface{}{"test"},
			shouldLog: true,
		},
		{
			name:      "should not log when level is Error",
			logLevel:  Error,
			format:    "warning message %s",
			args:      []interface{}{"test"},
			shouldLog: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockLogger := &MockLogger{}
			service := &LoggerService{
				LogLevel: tt.logLevel,
				Loggers:  []Logger{mockLogger},
			}

			// Execute
			service.Warn(tt.format, tt.args...)

			// Verify
			expectedMsg := fmt.Sprintf(tt.format, tt.args...)
			if tt.shouldLog {
				assert.Equal(t, expectedMsg, mockLogger.LastPrintedMessage.Message)
				assert.Equal(t, "warn", mockLogger.LastPrintedMessage.Level)
				assert.Equal(t, string(IconWarning), mockLogger.LastPrintedMessage.Icon)
			} else {
				assert.Empty(t, mockLogger.LastPrintedMessage.Message)
			}
		})
	}
}

func TestLoggerService_Command(t *testing.T) {
	tests := []struct {
		name      string
		logLevel  Level
		format    string
		args      []interface{}
		shouldLog bool
	}{
		{
			name:      "should log when level is Info",
			logLevel:  Info,
			format:    "command message %s",
			args:      []interface{}{"test"},
			shouldLog: true,
		},
		{
			name:      "should log when level is Debug",
			logLevel:  Debug,
			format:    "command message %s",
			args:      []interface{}{"test"},
			shouldLog: true,
		},
		{
			name:      "should not log when level is Error",
			logLevel:  Error,
			format:    "command message %s",
			args:      []interface{}{"test"},
			shouldLog: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockLogger := &MockLogger{}
			service := &LoggerService{
				LogLevel: tt.logLevel,
				Loggers:  []Logger{mockLogger},
			}

			// Execute
			service.Command(tt.format, tt.args...)

			// Verify
			expectedMsg := fmt.Sprintf(tt.format, tt.args...)
			if tt.shouldLog {
				assert.Equal(t, expectedMsg, mockLogger.LastPrintedMessage.Message)
				assert.Equal(t, "command", mockLogger.LastPrintedMessage.Level)
				assert.Equal(t, string(IconWrench), mockLogger.LastPrintedMessage.Icon)
			} else {
				assert.Empty(t, mockLogger.LastPrintedMessage.Message)
			}
		})
	}
}

func TestLoggerService_Disabled(t *testing.T) {
	tests := []struct {
		name      string
		logLevel  Level
		format    string
		args      []interface{}
		shouldLog bool
	}{
		{
			name:      "should log when level is Info",
			logLevel:  Info,
			format:    "disabled message %s",
			args:      []interface{}{"test"},
			shouldLog: true,
		},
		{
			name:      "should log when level is Debug",
			logLevel:  Debug,
			format:    "disabled message %s",
			args:      []interface{}{"test"},
			shouldLog: true,
		},
		{
			name:      "should not log when level is Error",
			logLevel:  Error,
			format:    "disabled message %s",
			args:      []interface{}{"test"},
			shouldLog: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockLogger := &MockLogger{}
			service := &LoggerService{
				LogLevel: tt.logLevel,
				Loggers:  []Logger{mockLogger},
			}

			// Execute
			service.Disabled(tt.format, tt.args...)

			// Verify
			expectedMsg := fmt.Sprintf(tt.format, tt.args...)
			if tt.shouldLog {
				assert.Equal(t, expectedMsg, mockLogger.LastPrintedMessage.Message)
				assert.Equal(t, "disabled", mockLogger.LastPrintedMessage.Level)
				assert.Equal(t, string(IconBlackSquare), mockLogger.LastPrintedMessage.Icon)
			} else {
				assert.Empty(t, mockLogger.LastPrintedMessage.Message)
			}
		})
	}
}

func TestLoggerService_Notice(t *testing.T) {
	tests := []struct {
		name      string
		logLevel  Level
		format    string
		args      []interface{}
		shouldLog bool
	}{
		{
			name:      "should log when level is Info",
			logLevel:  Info,
			format:    "notice message %s",
			args:      []interface{}{"test"},
			shouldLog: true,
		},
		{
			name:      "should log when level is Debug",
			logLevel:  Debug,
			format:    "notice message %s",
			args:      []interface{}{"test"},
			shouldLog: true,
		},
		{
			name:      "should not log when level is Error",
			logLevel:  Error,
			format:    "notice message %s",
			args:      []interface{}{"test"},
			shouldLog: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockLogger := &MockLogger{}
			service := &LoggerService{
				LogLevel: tt.logLevel,
				Loggers:  []Logger{mockLogger},
			}

			// Execute
			service.Notice(tt.format, tt.args...)

			// Verify
			expectedMsg := fmt.Sprintf(tt.format, tt.args...)
			if tt.shouldLog {
				assert.Equal(t, expectedMsg, mockLogger.LastPrintedMessage.Message)
				assert.Equal(t, "notice", mockLogger.LastPrintedMessage.Level)
				assert.Equal(t, string(IconFlag), mockLogger.LastPrintedMessage.Icon)
			} else {
				assert.Empty(t, mockLogger.LastPrintedMessage.Message)
			}
		})
	}
}

func TestLoggerService_Trace(t *testing.T) {
	tests := []struct {
		name      string
		logLevel  Level
		format    string
		args      []interface{}
		shouldLog bool
	}{
		{
			name:      "should log when level is Trace",
			logLevel:  Trace,
			format:    "trace message %s",
			args:      []interface{}{"test"},
			shouldLog: true,
		},
		{
			name:      "should not log when level is Debug",
			logLevel:  Debug,
			format:    "trace message %s",
			args:      []interface{}{"test"},
			shouldLog: false,
		},
		{
			name:      "should not log when level is Info",
			logLevel:  Info,
			format:    "trace message %s",
			args:      []interface{}{"test"},
			shouldLog: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockLogger := &MockLogger{}
			service := &LoggerService{
				LogLevel: tt.logLevel,
				Loggers:  []Logger{mockLogger},
			}

			// Execute
			service.Trace(tt.format, tt.args...)

			// Verify
			expectedMsg := fmt.Sprintf(tt.format, tt.args...)
			if tt.shouldLog {
				assert.Equal(t, expectedMsg, mockLogger.LastPrintedMessage.Message)
				assert.Equal(t, "debug", mockLogger.LastPrintedMessage.Level) // Note: Trace uses Debug internally
				assert.Equal(t, string(IconFire), mockLogger.LastPrintedMessage.Icon)
			} else {
				assert.Empty(t, mockLogger.LastPrintedMessage.Message)
			}
		})
	}
}

func TestLoggerService_Fatal(t *testing.T) {
	tests := []struct {
		name      string
		logLevel  Level
		format    string
		args      []interface{}
		shouldLog bool
	}{
		{
			name:      "should log when level is Error",
			logLevel:  Error,
			format:    "fatal message %s",
			args:      []interface{}{"test"},
			shouldLog: true,
		},
		{
			name:      "should log when level is Warning",
			logLevel:  Warning,
			format:    "fatal message %s",
			args:      []interface{}{"test"},
			shouldLog: true,
		},
		{
			name:      "should not log when level is Info",
			logLevel:  Info,
			format:    "fatal message %s",
			args:      []interface{}{"test"},
			shouldLog: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockLogger := &MockLogger{}
			service := &LoggerService{
				LogLevel: tt.logLevel,
				Loggers:  []Logger{mockLogger},
			}

			// Execute
			service.Fatal(tt.format, tt.args...)

			// Verify
			expectedMsg := fmt.Sprintf(tt.format, tt.args...)
			if tt.shouldLog {
				assert.Equal(t, expectedMsg, mockLogger.LastPrintedMessage.Message)
				assert.Equal(t, "error", mockLogger.LastPrintedMessage.Level)
				assert.Equal(t, string(IconRevolvingLight), mockLogger.LastPrintedMessage.Icon)
			} else {
				assert.Empty(t, mockLogger.LastPrintedMessage.Message)
			}
		})
	}
}

func TestLoggerService_OnMessage(t *testing.T) {
	// Setup
	service := New()
	channelLogger := &ChannelLogger{}
	channelLogger = channelLogger.Init().(*ChannelLogger)
	service.Loggers = append(service.Loggers, channelLogger)

	// Create channels to receive messages from each subscriber
	messages1 := make(chan LogMessage, 100)
	messages2 := make(chan LogMessage, 100)

	// Create WaitGroup for synchronization
	var wg sync.WaitGroup

	// First subscriber
	service.OnMessage("sub1", func(msg LogMessage) {
		messages1 <- msg
		wg.Done()
	})

	// Second subscriber
	service.OnMessage("sub2", func(msg LogMessage) {
		messages2 <- msg
		wg.Done()
	})

	// Test single message delivery
	t.Run("single message delivery", func(t *testing.T) {
		wg.Add(2) // One for each subscriber

		// Send test message
		service.Info("test message")

		// Wait for both subscribers with timeout
		if !waitTimeout(&wg, 5*time.Second) {
			t.Fatal("timeout waiting for messages")
		}

		// Verify both subscribers received the message
		msg1 := <-messages1
		assert.Equal(t, "info", msg1.Level)
		assert.Equal(t, "test message", msg1.Message)
		assert.Equal(t, string(IconInfo), string(msg1.Icon))

		msg2 := <-messages2
		assert.Equal(t, "info", msg2.Level)
		assert.Equal(t, "test message", msg2.Message)
		assert.Equal(t, string(IconInfo), string(msg2.Icon))
	})

	// Test multiple message delivery
	t.Run("multiple message delivery", func(t *testing.T) {
		messageCount := 3
		wg.Add(messageCount * 2) // Messages * Subscribers

		// Send multiple messages
		for i := 0; i < messageCount; i++ {
			service.Info(fmt.Sprintf("message %d", i))
		}

		// Wait for all messages with timeout
		if !waitTimeout(&wg, 5*time.Second) {
			t.Fatal("timeout waiting for messages")
		}

		// Verify all messages were received by both subscribers
		for i := 0; i < messageCount; i++ {
			msg1 := <-messages1
			assert.Contains(t, msg1.Message, "message")
			assert.Equal(t, "info", msg1.Level)

			msg2 := <-messages2
			assert.Contains(t, msg2.Message, "message")
			assert.Equal(t, "info", msg2.Level)
		}
	})
}

// Helper function to wait with timeout
func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return true
	case <-time.After(timeout):
		return false
	}
}

func TestMain(m *testing.M) {
	// Run all tests
	code := m.Run()

	// Cleanup after all tests
	cleanup()

	// Exit with the test result code
	os.Exit(code)
}

func cleanup() {
	// List of test files that might be created
	testFiles := []string{
		"test.log",
		"test.log.01",
		"test.log.02",
		"test.log.03",
		"test.log.04",
		"test.log.05",
		"test.log.06",
		"test.log.07",
		"test.log.08",
		"test.log.09",
	}

	// Remove each test file if it exists
	for _, file := range testFiles {
		if _, err := os.Stat(file); err == nil {
			err := os.Remove(file)
			if err != nil {
				fmt.Printf("Error removing test file %s: %v\n", file, err)
			}
		}
	}
}
