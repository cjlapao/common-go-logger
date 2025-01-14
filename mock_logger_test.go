package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockLogger_Log(t *testing.T) {
	tests := []struct {
		name     string
		level    Level
		format   string
		args     []interface{}
		expected MockedLogMessage
	}{
		{
			name:   "error level",
			level:  Error,
			format: "error message %s",
			args:   []interface{}{"test"},
			expected: MockedLogMessage{
				Level:   "error",
				Message: "error message test",
				Icon:    "",
			},
		},
		{
			name:   "warn level",
			level:  Warning,
			format: "warn message %s",
			args:   []interface{}{"test"},
			expected: MockedLogMessage{
				Level:   "warn",
				Message: "warn message test",
				Icon:    "",
			},
		},
		{
			name:   "info level",
			level:  Info,
			format: "info message %s",
			args:   []interface{}{"test"},
			expected: MockedLogMessage{
				Level:   "info",
				Message: "info message test",
				Icon:    "",
			},
		},
		{
			name:   "debug level",
			level:  Debug,
			format: "debug message %s",
			args:   []interface{}{"test"},
			expected: MockedLogMessage{
				Level:   "debug",
				Message: "debug message test",
				Icon:    "",
			},
		},
		{
			name:   "trace level",
			level:  Trace,
			format: "trace message %s",
			args:   []interface{}{"test"},
			expected: MockedLogMessage{
				Level:   "trace",
				Message: "trace message test",
				Icon:    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockLogger := &MockLogger{}
			mockLogger = mockLogger.Init().(*MockLogger)

			// Execute
			mockLogger.Log(tt.format, tt.level, tt.args...)

			// Verify
			assert.Equal(t, tt.expected.Level, mockLogger.LastPrintedMessage.Level)
			assert.Equal(t, tt.expected.Message, mockLogger.LastPrintedMessage.Message)
			assert.Equal(t, tt.expected.Icon, mockLogger.LastPrintedMessage.Icon)

			// Verify message was added to history
			assert.Len(t, mockLogger.PrintedMessages, 1)
			assert.Equal(t, tt.expected, mockLogger.PrintedMessages[0])
		})
	}

	t.Run("multiple messages", func(t *testing.T) {
		// Setup
		mockLogger := &MockLogger{}
		mockLogger = mockLogger.Init().(*MockLogger)

		// Execute multiple logs
		mockLogger.Log("first %s", Info, "message")
		mockLogger.Log("second %s", Warning, "message")
		mockLogger.Log("third %s", Error, "message")

		// Verify message history
		assert.Len(t, mockLogger.PrintedMessages, 3)
		assert.Equal(t, "info", mockLogger.PrintedMessages[0].Level)
		assert.Equal(t, "warn", mockLogger.PrintedMessages[1].Level)
		assert.Equal(t, "error", mockLogger.PrintedMessages[2].Level)
	})

	t.Run("clear messages", func(t *testing.T) {
		// Setup
		mockLogger := &MockLogger{}
		mockLogger = mockLogger.Init().(*MockLogger)

		// Log some messages
		mockLogger.Log("test %s", Info, "message")
		assert.Len(t, mockLogger.PrintedMessages, 1)

		// Clear messages
		mockLogger.Clear()

		// Verify messages were cleared
		assert.Len(t, mockLogger.PrintedMessages, 0)
		assert.Empty(t, mockLogger.LastPrintedMessage.Message)
	})
}

func TestMockLogger_LogIcon(t *testing.T) {
	tests := []struct {
		name     string
		level    Level
		icon     LoggerIcon
		format   string
		args     []interface{}
		expected MockedLogMessage
	}{
		{
			name:   "error level with custom icon",
			level:  Error,
			icon:   "🔴",
			format: "error message %s",
			args:   []interface{}{"test"},
			expected: MockedLogMessage{
				Level:   "error",
				Message: "error message test",
				Icon:    "🔴",
			},
		},
		{
			name:   "info level with custom icon",
			level:  Info,
			icon:   "💡",
			format: "info message %s",
			args:   []interface{}{"test"},
			expected: MockedLogMessage{
				Level:   "info",
				Message: "info message test",
				Icon:    "💡",
			},
		},
		{
			name:   "warn level with standard icon",
			level:  Warning,
			icon:   IconWarning,
			format: "warn message %s",
			args:   []interface{}{"test"},
			expected: MockedLogMessage{
				Level:   "warn",
				Message: "warn message test",
				Icon:    string(IconWarning),
			},
		},
		{
			name:   "debug level with empty icon",
			level:  Debug,
			icon:   "",
			format: "debug message %s",
			args:   []interface{}{"test"},
			expected: MockedLogMessage{
				Level:   "debug",
				Message: "debug message test",
				Icon:    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockLogger := &MockLogger{}
			mockLogger = mockLogger.Init().(*MockLogger)

			// Execute
			mockLogger.LogIcon(tt.icon, tt.format, tt.level, tt.args...)

			// Verify
			assert.Equal(t, tt.expected.Level, mockLogger.LastPrintedMessage.Level)
			assert.Equal(t, tt.expected.Message, mockLogger.LastPrintedMessage.Message)
			assert.Equal(t, tt.expected.Icon, mockLogger.LastPrintedMessage.Icon)

			// Verify message was added to history
			assert.Len(t, mockLogger.PrintedMessages, 1)
			assert.Equal(t, tt.expected, mockLogger.PrintedMessages[0])
		})
	}

	t.Run("multiple messages with icons", func(t *testing.T) {
		// Setup
		mockLogger := &MockLogger{}
		mockLogger = mockLogger.Init().(*MockLogger)

		// Execute multiple logs
		mockLogger.LogIcon("📘", "first %s", Info, "message")
		mockLogger.LogIcon("⚠️", "second %s", Warning, "message")
		mockLogger.LogIcon("❌", "third %s", Error, "message")

		// Verify message history
		assert.Len(t, mockLogger.PrintedMessages, 3)
		assert.Equal(t, "info", mockLogger.PrintedMessages[0].Level)
		assert.Equal(t, "📘", mockLogger.PrintedMessages[0].Icon)
		assert.Equal(t, "warn", mockLogger.PrintedMessages[1].Level)
		assert.Equal(t, "⚠️", mockLogger.PrintedMessages[1].Icon)
		assert.Equal(t, "error", mockLogger.PrintedMessages[2].Level)
		assert.Equal(t, "❌", mockLogger.PrintedMessages[2].Icon)
	})

	t.Run("with icons enabled", func(t *testing.T) {
		// Setup
		mockLogger := &MockLogger{}
		mockLogger = mockLogger.Init().(*MockLogger)
		mockLogger.UseIcons(true)

		// Execute
		mockLogger.LogIcon(IconInfo, "test %s", Info, "message")

		// Verify
		assert.Equal(t, string(IconInfo), mockLogger.LastPrintedMessage.Icon)
		assert.Equal(t, "info", mockLogger.LastPrintedMessage.Level)
		assert.Equal(t, "test message", mockLogger.LastPrintedMessage.Message)
	})
}
