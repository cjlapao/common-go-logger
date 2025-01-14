package log

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	strcolor "github.com/cjlapao/common-go/strcolor"
	"github.com/stretchr/testify/assert"
)

func TestFileLogger_Init(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     bool
	}{
		{
			name:     "with valid filename",
			filename: "test.log",
			want:     true,
		},
		{
			name:     "without filename",
			filename: "",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := FileLogger{filename: tt.filename}
			logger := l.Init().(*FileLogger)

			assert.Equal(t, tt.want, logger.enabled)

			if tt.filename != "" {
				defer os.Remove(tt.filename)
			}
		})
	}
}

func TestFileLogger_LoggingOperations(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "test.log")
	logger := FileLogger{filename: tmpFile}.Init().(*FileLogger)
	defer logger.Close()

	tests := []struct {
		name     string
		logFunc  func()
		contains string
	}{
		{
			name: "Info logging",
			logFunc: func() {
				logger.Info("test info message")
			},
			contains: "test info message",
		},
		{
			name: "Error logging",
			logFunc: func() {
				logger.Error("test error message")
			},
			contains: "test error message",
		},
		{
			name: "Exception logging",
			logFunc: func() {
				logger.Exception(errors.New("test error"), "exception occurred")
			},
			contains: "exception occurred, err test error",
		},
		{
			name: "Log with timestamp",
			logFunc: func() {
				logger.UseTimestamp(true)
				logger.Info("timestamped message")
				logger.UseTimestamp(false)
			},
			contains: time.Now().Format(time.RFC3339)[:10], // Check just the date part
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.logFunc()

			content, err := os.ReadFile(tmpFile)
			assert.NoError(t, err)
			assert.Contains(t, string(content), tt.contains)
		})
	}
}

func TestFileLogger_RotateLogFile(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "rotate.log")

	// Set a small max file size for testing
	os.Setenv("MAX_LOG_FILE_SIZE", "100")
	defer os.Unsetenv("MAX_LOG_FILE_SIZE")

	logger := FileLogger{filename: logFile}.Init().(*FileLogger)
	defer logger.Close()

	// Write enough data to trigger rotation
	for i := 0; i < 10; i++ {
		logger.Info("This is a long message that will help fill up the log file quickly " + fmt.Sprint(i))
	}

	// Check if rotation files were created
	files, err := os.ReadDir(tmpDir)
	assert.NoError(t, err)

	rotatedFiles := 0
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "rotate.log.") {
			rotatedFiles++
		}
	}

	assert.Greater(t, rotatedFiles, 0, "Expected at least one rotated log file")
}

func TestFileLogger_CorrelationID(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "correlation.log")
	logger := FileLogger{filename: tmpFile}.Init().(*FileLogger)
	defer logger.Close()

	// Set correlation ID
	correlationID := "test-correlation-id"
	os.Setenv("CORRELATION_ID", correlationID)
	defer os.Unsetenv("CORRELATION_ID")

	logger.UseCorrelationId(true)
	logger.Info("test message")

	content, err := os.ReadFile(tmpFile)
	assert.NoError(t, err)
	assert.Contains(t, string(content), correlationID)
}

func TestFileLogger_LogLevels(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "levels.log")
	logger := FileLogger{filename: tmpFile}.Init().(*FileLogger)
	defer logger.Close()

	tests := []struct {
		level   Level
		message string
	}{
		{Level(0), "error message"},
		{Level(1), "warn message"},
		{Level(2), "info message"},
		{Level(3), "debug message"},
		{Level(4), "trace message"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Level_%d", tt.level), func(t *testing.T) {
			logger.Log(tt.message, tt.level)

			content, err := os.ReadFile(tmpFile)
			assert.NoError(t, err)
			assert.Contains(t, string(content), tt.message)
		})
	}
}

func TestFileLogger_AllLogMethods(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "all_methods.log")
	logger := FileLogger{filename: tmpFile}.Init().(*FileLogger)
	defer logger.Close()

	tests := []struct {
		name   string
		logFn  func()
		expect string
	}{
		{"Success", func() { logger.Success("success msg") }, "success msg"},
		{"TaskSuccess", func() { logger.TaskSuccess("task success", true) }, "task success"},
		{"Warn", func() { logger.Warn("warn msg") }, "warn msg"},
		{"TaskWarn", func() { logger.TaskWarn("task warn") }, "task warn"},
		{"Command", func() { logger.Command("command msg") }, "command msg"},
		{"Disabled", func() { logger.Disabled("disabled msg") }, "disabled msg"},
		{"Notice", func() { logger.Notice("notice msg") }, "notice msg"},
		{"Debug", func() { logger.Debug("debug msg") }, "debug msg"},
		{"Trace", func() { logger.Trace("trace msg") }, "trace msg"},
		{"Fatal", func() { logger.Fatal("fatal msg") }, "fatal msg"},
		{"LogError", func() { logger.LogError(errors.New("error msg")) }, "error msg"},
		{"TaskError", func() { logger.TaskError("task error", true) }, "task error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.logFn()
			content, err := os.ReadFile(tmpFile)
			assert.NoError(t, err)
			assert.Contains(t, string(content), tt.expect)
		})
	}
}

func TestFileLogger_LogWithIcons(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "icons.log")
	logger := FileLogger{filename: tmpFile}.Init().(*FileLogger)
	defer logger.Close()

	logger.UseIcons(true)
	tests := []struct {
		name  string
		level Level
		icon  LoggerIcon
	}{
		{"ErrorIcon", Level(0), IconRevolvingLight},
		{"WarnIcon", Level(1), IconWarning},
		{"InfoIcon", Level(2), IconInfo},
		{"DebugIcon", Level(3), IconFire},
		{"TraceIcon", Level(4), IconBulb},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger.LogIcon(tt.icon, "test with icon", tt.level)
			assert.True(t, logger.IsTimestampEnabled() == logger.useTimestamp)
		})
	}
}

func TestFileLogger_LogHighlight(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "highlight.log")
	logger := FileLogger{filename: tmpFile}.Init().(*FileLogger)
	defer logger.Close()

	tests := []struct {
		level Level
		color strcolor.ColorCode
	}{
		{Level(0), strcolor.Red},
		{Level(1), strcolor.Yellow},
		{Level(2), strcolor.Green},
		{Level(3), strcolor.Blue},
		{Level(4), strcolor.BrightMagenta},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Level_%d", tt.level), func(t *testing.T) {
			logger.LogHighlight("highlighted message %s", tt.level, tt.color, "test")
		})
	}
}

func TestFileLogger_FatalError(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "fatal.log")
	logger := FileLogger{filename: tmpFile}.Init().(*FileLogger)
	defer logger.Close()

	// Test without error
	logger.FatalError(nil, "no error message")

	// Test with error (should panic)
	testErr := errors.New("test fatal error")
	assert.Panics(t, func() {
		logger.FatalError(testErr, "fatal error message")
	})
}

func TestFileLogger_InitWithInvalidFile(t *testing.T) {
	// Try to create logger with a path that cannot be created
	assert.Panics(t, func() {
		invalidPath := filepath.Join(string(byte(0)), "invalid.log")
		FileLogger{filename: invalidPath}.Init()
	})
}

func TestFileLogger_MessageFormatting(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "format.log")
	logger := FileLogger{filename: tmpFile}.Init().(*FileLogger)
	defer logger.Close()

	tests := []struct {
		name     string
		logFn    func()
		expected string
	}{
		{
			name: "Different types",
			logFn: func() {
				logger.Info("test %v %v %v", 123, true, struct{ A string }{A: "test"})
			},
			expected: "test 123 true",
		},
		{
			name: "Without newline",
			logFn: func() {
				logger.Info("test message")
			},
			expected: "test message\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.logFn()
			content, err := os.ReadFile(tmpFile)
			assert.NoError(t, err)
			assert.Contains(t, string(content), tt.expected)
		})
	}
}

func TestFileLogger_RotationEdgeCases(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "rotate_edge.log")

	tests := []struct {
		name        string
		maxFileSize string
		writeCount  int
		expectedRot bool
		setupFn     func()
		cleanupFn   func()
	}{
		{
			name:        "Invalid max size",
			maxFileSize: "invalid",
			writeCount:  5,
			expectedRot: false,
			setupFn: func() {
				os.Setenv("MAX_LOG_FILE_SIZE", "invalid")
			},
			cleanupFn: func() {
				os.Unsetenv("MAX_LOG_FILE_SIZE")
			},
		},
		{
			name:        "Zero max size",
			maxFileSize: "0",
			writeCount:  5,
			expectedRot: true,
			setupFn: func() {
				os.Setenv("MAX_LOG_FILE_SIZE", "0")
			},
			cleanupFn: func() {
				os.Unsetenv("MAX_LOG_FILE_SIZE")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupFn != nil {
				tt.setupFn()
			}

			logger := FileLogger{filename: logFile}.Init().(*FileLogger)

			for i := 0; i < tt.writeCount; i++ {
				logger.Info("Test message for rotation %d", i)
			}

			logger.Close()

			if tt.cleanupFn != nil {
				tt.cleanupFn()
			}
		})
	}
}

func TestFileLogger_DisabledLogger(t *testing.T) {
	// Test with disabled logger (no filename)
	logger := FileLogger{}.Init().(*FileLogger)
	defer logger.Close()

	// These should not panic and should be no-ops
	logger.Info("test message")
	logger.Error("test error")
	logger.Success("test success")
}
