package log

import (
	"errors"
	"testing"
	"time"

	strcolor "github.com/cjlapao/common-go/strcolor"
	"github.com/stretchr/testify/assert"
)

func TestLogMessage_String(t *testing.T) {
	fixedTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		message  LogMessage
		expected string
	}{
		{
			name: "with icon",
			message: LogMessage{
				Level:     "info",
				Message:   "test message",
				Timestamp: fixedTime,
				Icon:      "📌",
			},
			expected: "[2024-01-01T12:00:00Z] 📌 info: test message",
		},
		{
			name: "without icon",
			message: LogMessage{
				Level:     "error",
				Message:   "error message",
				Timestamp: fixedTime,
			},
			expected: "[2024-01-01T12:00:00Z] error: error message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.message.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestChannelLogger_Init(t *testing.T) {
	logger := &ChannelLogger{}
	initialized := logger.Init().(*ChannelLogger)

	assert.False(t, initialized.useTimestamp)
	assert.False(t, initialized.userCorrelationId)
	assert.False(t, initialized.useIcons)
	assert.Empty(t, initialized.subscribers)
}

func TestChannelLogger_Configuration(t *testing.T) {
	logger := &ChannelLogger{}
	logger = logger.Init().(*ChannelLogger)

	t.Run("timestamp configuration", func(t *testing.T) {
		logger.UseTimestamp(true)
		assert.True(t, logger.IsTimestampEnabled())
		logger.UseTimestamp(false)
		assert.False(t, logger.IsTimestampEnabled())
	})

	t.Run("correlation ID configuration", func(t *testing.T) {
		logger.UseCorrelationId(true)
		assert.True(t, logger.userCorrelationId)
		logger.UseCorrelationId(false)
		assert.False(t, logger.userCorrelationId)
	})

	t.Run("icons configuration", func(t *testing.T) {
		logger.UseIcons(true)
		assert.True(t, logger.useIcons)
		logger.UseIcons(false)
		assert.False(t, logger.useIcons)
	})
}

func TestChannelLogger_Subscribe(t *testing.T) {
	logger := &ChannelLogger{}
	logger = logger.Init().(*ChannelLogger)

	t.Run("subscribe with empty ID", func(t *testing.T) {
		id, ch := logger.Subscribe("", func(msg LogMessage) bool { return true })
		assert.NotEmpty(t, id)
		assert.NotNil(t, ch)
		assert.Len(t, logger.subscribers, 1)
	})

	t.Run("subscribe with custom ID", func(t *testing.T) {
		customID := "test-id"
		id, ch := logger.Subscribe(customID, func(msg LogMessage) bool { return true })
		assert.Equal(t, "sub_"+customID, id)
		assert.NotNil(t, ch)
	})

	t.Run("subscribe with existing ID", func(t *testing.T) {
		customID := "duplicate-id"
		id1, ch1 := logger.Subscribe(customID, func(msg LogMessage) bool { return true })
		id2, ch2 := logger.Subscribe(customID, func(msg LogMessage) bool { return true })
		assert.Equal(t, id1, id2)
		assert.Equal(t, ch1, ch2)
	})
}

func TestChannelLogger_Unsubscribe(t *testing.T) {
	logger := &ChannelLogger{}
	logger = logger.Init().(*ChannelLogger)

	t.Run("unsubscribe existing subscription", func(t *testing.T) {
		id, ch := logger.Subscribe("test", func(msg LogMessage) bool { return true })
		success := logger.Unsubscribe(id)
		assert.True(t, success)

		// Verify channel is closed
		_, ok := <-ch
		assert.False(t, ok)
	})

	t.Run("unsubscribe non-existent subscription", func(t *testing.T) {
		success := logger.Unsubscribe("non-existent")
		assert.False(t, success)
	})
}

func TestChannelLogger_Channel(t *testing.T) {
	logger := &ChannelLogger{}
	logger = logger.Init().(*ChannelLogger)

	id, ch := logger.Channel()
	assert.NotEmpty(t, id)
	assert.NotNil(t, ch)
	assert.Len(t, logger.subscribers, 1)
}

func TestChannelLogger_MessageFiltering(t *testing.T) {
	logger := &ChannelLogger{}
	logger = logger.Init().(*ChannelLogger)

	// Subscribe to error messages only
	_, errorCh := logger.Subscribe("", func(msg LogMessage) bool {
		return msg.Level == "error"
	})

	// Subscribe to info messages only
	_, infoCh := logger.Subscribe("", func(msg LogMessage) bool {
		return msg.Level == "info"
	})

	go func() {
		logger.Info("info message")
		logger.Error("error message")
	}()

	// Check error channel
	select {
	case msg := <-errorCh:
		assert.Equal(t, "error", msg.Level)
		assert.Equal(t, "error message", msg.Message)
	case <-time.After(time.Second):
		t.Error("timeout waiting for error message")
	}

	// Check info channel
	select {
	case msg := <-infoCh:
		assert.Equal(t, "info", msg.Level)
		assert.Equal(t, "info message", msg.Message)
	case <-time.After(time.Second):
		t.Error("timeout waiting for info message")
	}
}

func TestChannelLogger_Close(t *testing.T) {
	logger := &ChannelLogger{}
	logger = logger.Init().(*ChannelLogger)

	// Create multiple subscriptions
	id1, ch1 := logger.Subscribe("", func(msg LogMessage) bool { return true })
	id2, ch2 := logger.Subscribe("", func(msg LogMessage) bool { return true })

	logger.Close()

	// Verify all channels are closed
	_, ok1 := <-ch1
	assert.False(t, ok1)
	_, ok2 := <-ch2
	assert.False(t, ok2)

	// Verify subscribers are cleared
	assert.Nil(t, logger.subscribers)

	// Verify unsubscribe after close
	assert.False(t, logger.Unsubscribe(id1))
	assert.False(t, logger.Unsubscribe(id2))
}

func TestChannelLogger_LoggingMethods(t *testing.T) {
	logger := &ChannelLogger{}
	logger = logger.Init().(*ChannelLogger)

	id, ch := logger.Channel()
	defer logger.Unsubscribe(id)

	tests := []struct {
		name     string
		logFunc  func()
		expected LogMessage
	}{
		{
			name:    "Info logging",
			logFunc: func() { logger.Info("test message") },
			expected: LogMessage{
				Level:   "info",
				Message: "test message",
				Icon:    IconInfo,
			},
		},
		{
			name:    "Error logging",
			logFunc: func() { logger.Error("error message") },
			expected: LogMessage{
				Level:   "error",
				Message: "error message",
				Icon:    IconRevolvingLight,
			},
		},
		{
			name:    "Warning logging",
			logFunc: func() { logger.Warn("warning message") },
			expected: LogMessage{
				Level:   "warn",
				Message: "warning message",
				Icon:    IconWarning,
			},
		},
		// Add more test cases for other logging methods
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.logFunc()
			select {
			case msg := <-ch:
				assert.Equal(t, tt.expected.Level, msg.Level)
				assert.Equal(t, tt.expected.Message, msg.Message)
				assert.Equal(t, tt.expected.Icon, msg.Icon)
			case <-time.After(time.Second):
				t.Fatal("timeout waiting for log message")
			}
		})
	}
}

func TestChannelLogger_FatalError(t *testing.T) {
	logger := &ChannelLogger{}
	logger = logger.Init().(*ChannelLogger)

	id, ch := logger.Channel()
	defer logger.Unsubscribe(id)

	err := errors.New("test error")

	// Test panic recovery
	defer func() {
		if r := recover(); r == nil {
			t.Error("FatalError should panic")
		}
	}()

	logger.FatalError(err, "fatal error occurred")

	select {
	case msg := <-ch:
		assert.Equal(t, "error", msg.Level)
		assert.Equal(t, "fatal error occurred", msg.Message)
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for log message")
	}
}

func TestChannelLogger_Log(t *testing.T) {
	logger := &ChannelLogger{}
	logger = logger.Init().(*ChannelLogger)

	id, ch := logger.Channel()
	defer logger.Unsubscribe(id)

	tests := []struct {
		name     string
		level    Level
		format   string
		words    []interface{}
		expected LogMessage
	}{
		{
			name:   "error level (0)",
			level:  0,
			format: "error %s",
			words:  []interface{}{"message"},
			expected: LogMessage{
				Level:   "error",
				Message: "error message",
				Icon:    "",
			},
		},
		{
			name:   "warn level (1)",
			level:  1,
			format: "warn %s",
			words:  []interface{}{"message"},
			expected: LogMessage{
				Level:   "warn",
				Message: "warn message",
				Icon:    "",
			},
		},
		{
			name:   "info level (2)",
			level:  2,
			format: "info %s",
			words:  []interface{}{"message"},
			expected: LogMessage{
				Level:   "info",
				Message: "info message",
				Icon:    "",
			},
		},
		{
			name:   "debug level (3)",
			level:  3,
			format: "debug %s",
			words:  []interface{}{"message"},
			expected: LogMessage{
				Level:   "debug",
				Message: "debug message",
				Icon:    "",
			},
		},
		{
			name:   "trace level (4)",
			level:  4,
			format: "trace %s",
			words:  []interface{}{"message"},
			expected: LogMessage{
				Level:   "trace",
				Message: "trace message",
				Icon:    "",
			},
		},
		{
			name:   "multiple format arguments",
			level:  2,
			format: "%s: value=%d, active=%v",
			words:  []interface{}{"test", 42, true},
			expected: LogMessage{
				Level:   "info",
				Message: "test: value=42, active=true",
				Icon:    "",
			},
		},
		{
			name:   "no format arguments",
			level:  2,
			format: "simple message",
			words:  []interface{}{},
			expected: LogMessage{
				Level:   "info",
				Message: "simple message",
				Icon:    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger.Log(tt.format, tt.level, tt.words...)

			select {
			case msg := <-ch:
				assert.Equal(t, tt.expected.Level, msg.Level)
				assert.Equal(t, tt.expected.Message, msg.Message)
				assert.Equal(t, tt.expected.Icon, msg.Icon)
				assert.NotZero(t, msg.Timestamp)
			case <-time.After(time.Second):
				t.Fatal("timeout waiting for log message")
			}
		})
	}

	// Test with no subscribers
	logger.Close()
	logger.Log("this should not panic", 2)
}

func TestChannelLogger_LogIcon(t *testing.T) {
	logger := &ChannelLogger{}
	logger = logger.Init().(*ChannelLogger)

	id, ch := logger.Channel()
	defer logger.Unsubscribe(id)

	tests := []struct {
		name     string
		level    Level
		icon     LoggerIcon
		format   string
		words    []interface{}
		expected LogMessage
	}{
		{
			name:   "error level with icon",
			level:  0,
			icon:   "🚫",
			format: "error %s",
			words:  []interface{}{"message"},
			expected: LogMessage{
				Level:   "error",
				Message: "error message",
				Icon:    "🚫",
			},
		},
		{
			name:   "warn level with icon",
			level:  1,
			icon:   "⚠️",
			format: "warn %s",
			words:  []interface{}{"message"},
			expected: LogMessage{
				Level:   "warn",
				Message: "warn message",
				Icon:    "⚠️",
			},
		},
		{
			name:   "info level with icon",
			level:  2,
			icon:   "ℹ️",
			format: "info %s",
			words:  []interface{}{"message"},
			expected: LogMessage{
				Level:   "info",
				Message: "info message",
				Icon:    "ℹ️",
			},
		},
		{
			name:   "debug level with icon",
			level:  3,
			icon:   "🔍",
			format: "debug %s",
			words:  []interface{}{"message"},
			expected: LogMessage{
				Level:   "debug",
				Message: "debug message",
				Icon:    "🔍",
			},
		},
		{
			name:   "trace level with icon",
			level:  4,
			icon:   "🔎",
			format: "trace %s",
			words:  []interface{}{"message"},
			expected: LogMessage{
				Level:   "trace",
				Message: "trace message",
				Icon:    "🔎",
			},
		},
		{
			name:   "multiple format arguments with icon",
			level:  2,
			icon:   "📝",
			format: "%s: value=%d, active=%v",
			words:  []interface{}{"test", 42, true},
			expected: LogMessage{
				Level:   "info",
				Message: "test: value=42, active=true",
				Icon:    "📝",
			},
		},
		{
			name:   "empty icon",
			level:  2,
			icon:   "",
			format: "message with no icon",
			words:  []interface{}{},
			expected: LogMessage{
				Level:   "info",
				Message: "message with no icon",
				Icon:    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test with icons disabled
			logger.UseIcons(false)
			logger.LogIcon(tt.icon, tt.format, tt.level, tt.words...)

			select {
			case msg := <-ch:
				assert.Equal(t, tt.expected.Level, msg.Level)
				assert.Equal(t, tt.expected.Message, msg.Message)
				assert.Equal(t, tt.expected.Icon, msg.Icon)
				assert.NotZero(t, msg.Timestamp)
			case <-time.After(time.Second):
				t.Fatal("timeout waiting for log message")
			}

			// Test with icons enabled
			logger.UseIcons(true)
			logger.LogIcon(tt.icon, tt.format, tt.level, tt.words...)

			select {
			case msg := <-ch:
				assert.Equal(t, tt.expected.Level, msg.Level)
				if tt.icon != "" {
					assert.Equal(t, string(tt.icon)+" "+tt.expected.Message, msg.Message)
				} else {
					assert.Equal(t, tt.expected.Message, msg.Message)
				}
				assert.Equal(t, tt.expected.Icon, msg.Icon)
				assert.NotZero(t, msg.Timestamp)
			case <-time.After(time.Second):
				t.Fatal("timeout waiting for log message")
			}
		})
	}

	// Test with no subscribers
	logger.Close()
	logger.LogIcon("📌", "this should not panic", 2)
}

func TestChannelLogger_LogHighlight(t *testing.T) {
	logger := &ChannelLogger{}
	logger = logger.Init().(*ChannelLogger)

	id, ch := logger.Channel()
	defer logger.Unsubscribe(id)

	tests := []struct {
		name          string
		level         Level
		format        string
		highlightText []interface{}
		expected      LogMessage
	}{
		{
			name:          "error level with highlight",
			level:         0,
			format:        "error: %s occurred",
			highlightText: []interface{}{"critical failure"},
			expected: LogMessage{
				Level:   "error",
				Message: "error: \x1b[31mcritical failure\x1b[0m occurred",
				Icon:    "",
			},
		},
		{
			name:          "info level with multiple highlights",
			level:         2,
			format:        "values: %s, %s",
			highlightText: []interface{}{"abc", "123"},
			expected: LogMessage{
				Level:   "info",
				Message: "values: \x1b[31mabc\x1b[0m, \x1b[31m123\x1b[0m",
				Icon:    "",
			},
		},
		{
			name:          "warning level with multiple highlights",
			level:         1,
			format:        "values: %s, %s",
			highlightText: []interface{}{"abc", "123"},
			expected: LogMessage{
				Level:   "warn",
				Message: "values: \x1b[31mabc\x1b[0m, \x1b[31m123\x1b[0m",
				Icon:    "",
			},
		},
		{
			name:          "debug level with number",
			level:         3,
			format:        "count: %v",
			highlightText: []interface{}{42},
			expected: LogMessage{
				Level:   "debug",
				Message: "count: \x1b[31m42\x1b[0m",
				Icon:    "",
			},
		},
		{
			name:          "trace level with number",
			level:         4,
			format:        "count: %v",
			highlightText: []interface{}{42},
			expected: LogMessage{
				Level:   "trace",
				Message: "count: \x1b[31m42\x1b[0m",
				Icon:    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger.LogHighlight(tt.format, tt.level, strcolor.Red, tt.highlightText...)

			select {
			case msg := <-ch:
				assert.Equal(t, tt.expected.Level, msg.Level)
				assert.Equal(t, tt.expected.Message, msg.Message)
				assert.Equal(t, tt.expected.Icon, msg.Icon)
				assert.NotZero(t, msg.Timestamp)
			case <-time.After(time.Second):
				t.Fatal("timeout waiting for log message")
			}
		})
	}

	// Test with no subscribers
	logger.Close()
	logger.LogHighlight("this should not panic", 2, strcolor.Red, "test")
}

func TestChannelLogger_Success(t *testing.T) {
	// Create a new logger
	logger := &ChannelLogger{}
	logger = logger.Init().(*ChannelLogger)

	// Subscribe to the logger
	id, ch := logger.Subscribe("test", func(msg LogMessage) bool {
		return true // Accept all messages
	})
	defer logger.Unsubscribe(id)

	// Test cases
	tests := []struct {
		name     string
		format   string
		args     []interface{}
		expected LogMessage
	}{
		{
			name:   "simple success message",
			format: "Operation completed",
			args:   nil,
			expected: LogMessage{
				Level:   "success",
				Message: "Operation completed",
				Icon:    IconThumbsUp,
			},
		},
		{
			name:   "success message with formatting",
			format: "Created %d items",
			args:   []interface{}{42},
			expected: LogMessage{
				Level:   "success",
				Message: "Created 42 items",
				Icon:    IconThumbsUp,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call Success method
			logger.Success(tt.format, tt.args...)

			// Receive message from channel
			select {
			case msg := <-ch:
				// Verify level and icon
				if msg.Level != tt.expected.Level {
					t.Errorf("expected level %s, got %s", tt.expected.Level, msg.Level)
				}
				if msg.Icon != tt.expected.Icon {
					t.Errorf("expected icon %s, got %s", tt.expected.Icon, msg.Icon)
				}
				if msg.Message != tt.expected.Message {
					t.Errorf("expected message %s, got %s", tt.expected.Message, msg.Message)
				}
			case <-time.After(time.Second):
				t.Error("timeout waiting for message")
			}
		})
	}
}

func TestChannelLogger_Command(t *testing.T) {
	// Create a new logger
	logger := &ChannelLogger{}
	logger = logger.Init().(*ChannelLogger)

	// Subscribe to the logger
	id, ch := logger.Subscribe("test", func(msg LogMessage) bool {
		return true // Accept all messages
	})
	defer logger.Unsubscribe(id)

	// Test cases
	tests := []struct {
		name     string
		format   string
		args     []interface{}
		expected LogMessage
	}{
		{
			name:   "simple command message",
			format: "git pull",
			args:   nil,
			expected: LogMessage{
				Level:   "command",
				Message: "git pull",
				Icon:    IconWrench,
			},
		},
		{
			name:   "command message with formatting",
			format: "docker run -p %d:%d nginx",
			args:   []interface{}{8080, 80},
			expected: LogMessage{
				Level:   "command",
				Message: "docker run -p 8080:80 nginx",
				Icon:    IconWrench,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call Command method
			logger.Command(tt.format, tt.args...)

			// Receive message from channel
			select {
			case msg := <-ch:
				// Verify level and icon
				if msg.Level != tt.expected.Level {
					t.Errorf("expected level %s, got %s", tt.expected.Level, msg.Level)
				}
				if msg.Icon != tt.expected.Icon {
					t.Errorf("expected icon %s, got %s", tt.expected.Icon, msg.Icon)
				}
				if msg.Message != tt.expected.Message {
					t.Errorf("expected message %s, got %s", tt.expected.Message, msg.Message)
				}
			case <-time.After(time.Second):
				t.Error("timeout waiting for message")
			}
		})
	}
}

func TestChannelLogger_Disabled(t *testing.T) {
	// Create a new logger
	logger := &ChannelLogger{}
	logger = logger.Init().(*ChannelLogger)

	// Subscribe to the logger
	id, ch := logger.Subscribe("test", func(msg LogMessage) bool {
		return true // Accept all messages
	})
	defer logger.Unsubscribe(id)

	// Test cases
	tests := []struct {
		name     string
		format   string
		args     []interface{}
		expected LogMessage
	}{
		{
			name:   "simple disabled message",
			format: "Feature X is disabled",
			args:   nil,
			expected: LogMessage{
				Level:   "disabled",
				Message: "Feature X is disabled",
				Icon:    IconBlackSquare,
			},
		},
		{
			name:   "disabled message with formatting",
			format: "Feature %s is disabled in version %s",
			args:   []interface{}{"OAuth", "2.0"},
			expected: LogMessage{
				Level:   "disabled",
				Message: "Feature OAuth is disabled in version 2.0",
				Icon:    IconBlackSquare,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call Disabled method
			logger.Disabled(tt.format, tt.args...)

			// Receive message from channel
			select {
			case msg := <-ch:
				// Verify level and icon
				if msg.Level != tt.expected.Level {
					t.Errorf("expected level %s, got %s", tt.expected.Level, msg.Level)
				}
				if msg.Icon != tt.expected.Icon {
					t.Errorf("expected icon %s, got %s", tt.expected.Icon, msg.Icon)
				}
				if msg.Message != tt.expected.Message {
					t.Errorf("expected message %s, got %s", tt.expected.Message, msg.Message)
				}
			case <-time.After(time.Second):
				t.Error("timeout waiting for message")
			}
		})
	}
}

func TestChannelLogger_Notice(t *testing.T) {
	// Create a new logger
	logger := &ChannelLogger{}
	logger = logger.Init().(*ChannelLogger)

	// Subscribe to the logger
	id, ch := logger.Subscribe("test", func(msg LogMessage) bool {
		return true // Accept all messages
	})
	defer logger.Unsubscribe(id)

	// Test cases
	tests := []struct {
		name     string
		format   string
		args     []interface{}
		expected LogMessage
	}{
		{
			name:   "simple notice message",
			format: "System maintenance scheduled",
			args:   nil,
			expected: LogMessage{
				Level:   "notice",
				Message: "System maintenance scheduled",
				Icon:    IconFlag,
			},
		},
		{
			name:   "notice message with formatting",
			format: "Database backup starting in %d minutes on %s",
			args:   []interface{}{5, "primary server"},
			expected: LogMessage{
				Level:   "notice",
				Message: "Database backup starting in 5 minutes on primary server",
				Icon:    IconFlag,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call Notice method
			logger.Notice(tt.format, tt.args...)

			// Receive message from channel
			select {
			case msg := <-ch:
				// Verify level and icon
				if msg.Level != tt.expected.Level {
					t.Errorf("expected level %s, got %s", tt.expected.Level, msg.Level)
				}
				if msg.Icon != tt.expected.Icon {
					t.Errorf("expected icon %s, got %s", tt.expected.Icon, msg.Icon)
				}
				if msg.Message != tt.expected.Message {
					t.Errorf("expected message %s, got %s", tt.expected.Message, msg.Message)
				}
			case <-time.After(time.Second):
				t.Error("timeout waiting for message")
			}
		})
	}
}

func TestChannelLogger_Debug(t *testing.T) {
	// Create a new logger
	logger := &ChannelLogger{}
	logger = logger.Init().(*ChannelLogger)

	// Subscribe to the logger
	id, ch := logger.Subscribe("test", func(msg LogMessage) bool {
		return true // Accept all messages
	})
	defer logger.Unsubscribe(id)

	// Test cases
	tests := []struct {
		name     string
		format   string
		args     []interface{}
		expected LogMessage
	}{
		{
			name:   "simple debug message",
			format: "Connection pool status",
			args:   nil,
			expected: LogMessage{
				Level:   "debug",
				Message: "Connection pool status",
				Icon:    IconFire,
			},
		},
		{
			name:   "debug message with formatting",
			format: "Active connections: %d, Queue size: %d",
			args:   []interface{}{42, 7},
			expected: LogMessage{
				Level:   "debug",
				Message: "Active connections: 42, Queue size: 7",
				Icon:    IconFire,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call Debug method
			logger.Debug(tt.format, tt.args...)

			// Receive message from channel
			select {
			case msg := <-ch:
				// Verify level and icon
				if msg.Level != tt.expected.Level {
					t.Errorf("expected level %s, got %s", tt.expected.Level, msg.Level)
				}
				if msg.Icon != tt.expected.Icon {
					t.Errorf("expected icon %s, got %s", tt.expected.Icon, msg.Icon)
				}
				if msg.Message != tt.expected.Message {
					t.Errorf("expected message %s, got %s", tt.expected.Message, msg.Message)
				}
			case <-time.After(time.Second):
				t.Error("timeout waiting for message")
			}
		})
	}
}

func TestChannelLogger_Trace(t *testing.T) {
	// Create a new logger
	logger := &ChannelLogger{}
	logger = logger.Init().(*ChannelLogger)

	// Subscribe to the logger
	id, ch := logger.Subscribe("test", func(msg LogMessage) bool {
		return true // Accept all messages
	})
	defer logger.Unsubscribe(id)

	// Test cases
	tests := []struct {
		name     string
		format   string
		args     []interface{}
		expected LogMessage
	}{
		{
			name:   "simple trace message",
			format: "Function entry point",
			args:   nil,
			expected: LogMessage{
				Level:   "trace",
				Message: "Function entry point",
				Icon:    IconBulb,
			},
		},
		{
			name:   "trace message with formatting",
			format: "Method %s called with params: %v",
			args:   []interface{}{"ProcessData", []string{"a", "b", "c"}},
			expected: LogMessage{
				Level:   "trace",
				Message: "Method ProcessData called with params: [a b c]",
				Icon:    IconBulb,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call Trace method
			logger.Trace(tt.format, tt.args...)

			// Receive message from channel
			select {
			case msg := <-ch:
				// Verify level and icon
				if msg.Level != tt.expected.Level {
					t.Errorf("expected level %s, got %s", tt.expected.Level, msg.Level)
				}
				if msg.Icon != tt.expected.Icon {
					t.Errorf("expected icon %s, got %s", tt.expected.Icon, msg.Icon)
				}
				if msg.Message != tt.expected.Message {
					t.Errorf("expected message %s, got %s", tt.expected.Message, msg.Message)
				}
			case <-time.After(time.Second):
				t.Error("timeout waiting for message")
			}
		})
	}
}

func TestChannelLogger_Exception(t *testing.T) {
	// Create a new logger
	logger := &ChannelLogger{}
	logger = logger.Init().(*ChannelLogger)

	// Subscribe to the logger
	id, ch := logger.Subscribe("test", func(msg LogMessage) bool {
		return true // Accept all messages
	})
	defer logger.Unsubscribe(id)

	// Create test errors
	testErr1 := errors.New("database connection failed")
	testErr2 := errors.New("invalid configuration")

	// Test cases
	tests := []struct {
		name     string
		err      error
		format   string
		args     []interface{}
		expected LogMessage
	}{
		{
			name:   "exception with empty format",
			err:    testErr1,
			format: "",
			args:   nil,
			expected: LogMessage{
				Level:   "error",
				Message: "database connection failed",
				Icon:    IconRevolvingLight,
			},
		},
		{
			name:   "exception with format and no args",
			err:    testErr1,
			format: "Failed to initialize database",
			args:   nil,
			expected: LogMessage{
				Level:   "error",
				Message: "Failed to initialize database, err database connection failed",
				Icon:    IconRevolvingLight,
			},
		},
		{
			name:   "exception with format and args",
			err:    testErr2,
			format: "Configuration error in %s module",
			args:   []interface{}{"authentication"},
			expected: LogMessage{
				Level:   "error",
				Message: "Configuration error in authentication module, err invalid configuration",
				Icon:    IconRevolvingLight,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call Exception method
			logger.Exception(tt.err, tt.format, tt.args...)

			// Receive message from channel
			select {
			case msg := <-ch:
				// Verify level and icon
				if msg.Level != tt.expected.Level {
					t.Errorf("expected level %s, got %s", tt.expected.Level, msg.Level)
				}
				if msg.Icon != tt.expected.Icon {
					t.Errorf("expected icon %s, got %s", tt.expected.Icon, msg.Icon)
				}
				if msg.Message != tt.expected.Message {
					t.Errorf("expected message %s, got %s", tt.expected.Message, msg.Message)
				}
			case <-time.After(time.Second):
				t.Error("timeout waiting for message")
			}
		})
	}
}

func TestChannelLogger_LogError(t *testing.T) {
	// Create a new logger
	logger := &ChannelLogger{}
	logger = logger.Init().(*ChannelLogger)

	// Subscribe to the logger
	id, ch := logger.Subscribe("test", func(msg LogMessage) bool {
		return true // Accept all messages
	})
	defer logger.Unsubscribe(id)

	// Test cases
	tests := []struct {
		name     string
		err      error
		expected *LogMessage // Using pointer to handle nil case
	}{
		{
			name: "standard error message",
			err:  errors.New("file not found"),
			expected: &LogMessage{
				Level:   "error",
				Message: "file not found",
				Icon:    IconRevolvingLight,
			},
		},
		{
			name:     "nil error message",
			err:      nil,
			expected: nil, // No message should be sent
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call LogError method
			logger.LogError(tt.err)

			if tt.expected == nil {
				// For nil error case, verify no message is sent
				select {
				case msg := <-ch:
					t.Errorf("expected no message for nil error, got: %v", msg)
				case <-time.After(100 * time.Millisecond):
					// This is the expected case for nil error
				}
				return
			}

			// For non-nil error case, verify the message
			select {
			case msg := <-ch:
				if msg.Level != tt.expected.Level {
					t.Errorf("expected level %s, got %s", tt.expected.Level, msg.Level)
				}
				if msg.Icon != tt.expected.Icon {
					t.Errorf("expected icon %s, got %s", tt.expected.Icon, msg.Icon)
				}
				if msg.Message != tt.expected.Message {
					t.Errorf("expected message %s, got %s", tt.expected.Message, msg.Message)
				}
			case <-time.After(time.Second):
				t.Error("timeout waiting for message")
			}
		})
	}
}

func TestChannelLogger_Fatal(t *testing.T) {
	// Create a new logger
	logger := &ChannelLogger{}
	logger = logger.Init().(*ChannelLogger)

	// Subscribe to the logger
	id, ch := logger.Subscribe("test", func(msg LogMessage) bool {
		return true // Accept all messages
	})
	defer logger.Unsubscribe(id)

	// Test cases
	tests := []struct {
		name     string
		format   string
		args     []interface{}
		expected LogMessage
	}{
		{
			name:   "simple fatal message",
			format: "Application crashed",
			args:   nil,
			expected: LogMessage{
				Level:   "error",
				Message: "Application crashed",
				Icon:    IconRevolvingLight,
			},
		},
		{
			name:   "fatal message with formatting",
			format: "Fatal error in module %s: memory allocation failed at address 0x%x",
			args:   []interface{}{"UserAuth", 0xDEADBEEF},
			expected: LogMessage{
				Level:   "error",
				Message: "Fatal error in module UserAuth: memory allocation failed at address 0xdeadbeef",
				Icon:    IconRevolvingLight,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call Fatal method
			logger.Fatal(tt.format, tt.args...)

			// Receive message from channel
			select {
			case msg := <-ch:
				// Verify level and icon
				if msg.Level != tt.expected.Level {
					t.Errorf("expected level %s, got %s", tt.expected.Level, msg.Level)
				}
				if msg.Icon != tt.expected.Icon {
					t.Errorf("expected icon %s, got %s", tt.expected.Icon, msg.Icon)
				}
				if msg.Message != tt.expected.Message {
					t.Errorf("expected message %s, got %s", tt.expected.Message, msg.Message)
				}
			case <-time.After(time.Second):
				t.Error("timeout waiting for message")
			}
		})
	}
}
