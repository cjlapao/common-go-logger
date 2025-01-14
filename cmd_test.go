package log

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	strcolor "github.com/cjlapao/common-go/strcolor"
	"github.com/stretchr/testify/assert"
)

func TestCmdLogger_UseIcons(t *testing.T) {
	type args struct {
		value bool
	}
	tests := []struct {
		name string
		l    *CmdLogger
		args args
	}{
		{
			name: "icon enabled",
			l:    &CmdLogger{},
			args: args{
				value: true,
			},
		},
		{
			name: "icon enabled",
			l:    &CmdLogger{},
			args: args{
				value: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.UseIcons(tt.args.value)
			assert.Equal(t, tt.l.useIcons, tt.args.value)
		})
	}
}

func TestCmdLogger_UseTimestamp(t *testing.T) {
	type args struct {
		value bool
	}
	tests := []struct {
		name string
		l    *CmdLogger
		args args
	}{
		{
			name: "timestamp enabled",
			l:    &CmdLogger{},
			args: args{
				value: true,
			},
		},
		{
			name: "timestamp enabled",
			l:    &CmdLogger{},
			args: args{
				value: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.UseTimestamp(tt.args.value)
			assert.Equal(t, tt.l.useTimestamp, tt.args.value)
		})
	}
}

func TestCmdLogger_UseCorrelationId(t *testing.T) {
	type args struct {
		value bool
	}
	tests := []struct {
		name string
		l    *CmdLogger
		args args
	}{
		{
			name: "correlationId enabled",
			l:    &CmdLogger{},
			args: args{
				value: true,
			},
		},
		{
			name: "correlationId enabled",
			l:    &CmdLogger{},
			args: args{
				value: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.UseCorrelationId(tt.args.value)
			assert.Equal(t, tt.l.userCorrelationId, tt.args.value)
		})
	}
}

func TestCmdLogger_Log(t *testing.T) {
	type args struct {
		format string
		level  Level
		words  []interface{}
	}
	tests := []struct {
		name        string
		args        args
		expect      string
		useIcons    bool
		useTime     bool
		useCorrelId bool
		correlId    string
	}{
		{
			name: "info",
			args: args{
				format: "%s",
				level:  Info,
				words: []interface{}{
					"I am",
				},
			},
			expect: "\x1b[0mI am\x1b[0m\n",
		},
		{
			name: "info, no words",
			args: args{
				format: "just me",
				level:  Info,
			},
			expect: "\x1b[0mjust me\x1b[0m\n",
		},
		{
			name: "info, empty words",
			args: args{
				format: "just me",
				level:  Info,
				words:  []interface{}{},
			},
			expect: "\x1b[0mjust me\x1b[0m\n",
		},
		{
			name: "error",
			args: args{
				format: "%s",
				level:  Error,
				words: []interface{}{
					"some error",
				},
			},
			expect: "\x1b[31msome error\x1b[0m\n",
		},
		{
			name: "warn",
			args: args{
				format: "warning: %s",
				level:  Warning,
				words: []interface{}{
					"some warning",
				},
			},
			expect: "\x1b[33mwarning: some warning\x1b[0m\n",
		},
		{
			name: "info with icons",
			args: args{
				format: "%s",
				level:  Info,
				words: []interface{}{
					"with icons",
				},
			},
			useIcons: true,

			expect: "\x1b[0mwith icons\x1b[0m\n",
		},
		{
			name: "error with icons",
			args: args{
				format: "%s",
				level:  Error,
				words: []interface{}{
					"with icons",
				},
			},
			useIcons: true,
			expect:   "\x1b[31mwith icons\x1b[0m\n",
		},
		{
			name: "warn with icons",
			args: args{
				format: "%s",
				level:  Warning,
				words: []interface{}{
					"with icons",
				},
			},
			useIcons: true,
			expect:   "\x1b[33mwith icons\x1b[0m\n",
		},
		{
			name: "success with icons",
			args: args{
				format: "%s",
				level:  Info,
				words: []interface{}{
					"with icons",
				},
			},
			useIcons: true,
			expect:   "\x1b[0mwith icons\x1b[0m\n",
		},
		{
			name: "with timestamp",
			args: args{
				format: "%s",
				level:  Info,
				words: []interface{}{
					"with time",
				},
			},
			useTime: true,
			expect:  `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z with time\x1b\[0m\n`,
		},
		{
			name: "with correlation id",
			args: args{
				format: "%s",
				level:  Info,
				words: []interface{}{
					"with correlation",
				},
			},
			useCorrelId: true,
			correlId:    "test-123",
			expect:      "\x1b[0m[test-123] with correlation\x1b[0m\n",
		},
		{
			name: "with all features enabled",
			args: args{
				format: "%s",
				level:  Info,
				words: []interface{}{
					"all features",
				},
			},
			useIcons:    true,
			useTime:     true,
			useCorrelId: true,
			correlId:    "test-123",
			expect:      `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z \[test-123\] all features\x1b\[0m\n`,
		},
		{
			name: "debug level with all features",
			args: args{
				format: "%s",
				level:  Debug,
				words: []interface{}{
					"debug message",
				},
			},
			useIcons:    true,
			useTime:     true,
			useCorrelId: true,
			correlId:    "test-123",
			expect:      `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z \[test-123\] debug message\x1b\[0m\n`,
		},
		{
			name: "trace level with all features",
			args: args{
				format: "%s",
				level:  Trace,
				words: []interface{}{
					"trace message",
				},
			},
			useIcons:    true,
			useTime:     true,
			useCorrelId: true,
			correlId:    "test-123",
			expect:      `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z \[test-123\] trace message\x1b\[0m\n`,
		},
		{
			name: "warning level with all features",
			args: args{
				format: "%s",
				level:  Warning,
				words: []interface{}{
					"warning message",
				},
			},
			useIcons:    true,
			useTime:     true,
			useCorrelId: true,
			correlId:    "test-123",
			expect:      `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z \[test-123\] warning message\x1b\[0m\n`,
		},
		{
			name: "error level with all features",
			args: args{
				format: "%s",
				level:  Error,
				words: []interface{}{
					"error message",
				},
			},
			useIcons:    true,
			useTime:     true,
			useCorrelId: true,
			correlId:    "test-123",
			expect:      `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z \[test-123\] error message\x1b\[0m\n`,
		},
		{
			name: "success level with all features",
			args: args{
				format: "%s",
				level:  Info,
				words: []interface{}{
					"success message",
				},
			},
			useIcons:    true,
			useTime:     true,
			useCorrelId: true,
			correlId:    "test-123",
			expect:      `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z \[test-123\] success message\x1b\[0m\n`,
		},
		{
			name: "multiple format arguments",
			args: args{
				format: "%s %d %s",
				level:  Info,
				words: []interface{}{
					"test",
					42,
					"args",
				},
			},
			useIcons:    true,
			useTime:     true,
			useCorrelId: true,
			correlId:    "test-123",
			expect:      `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z \[test-123\] test 42 args\x1b\[0m\n`,
		},
		{
			name: "no correlation ID in env",
			args: args{
				format: "%s",
				level:  Info,
				words: []interface{}{
					"no correlation",
				},
			},
			useCorrelId: true,
			expect:      "\x1b[0mno correlation\x1b[0m\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output bytes.Buffer
			l := &CmdLogger{
				writer: &output,
			}

			l.UseIcons(tt.useIcons)
			l.UseTimestamp(tt.useTime)
			l.UseCorrelationId(tt.useCorrelId)
			if tt.correlId != "" {
				os.Setenv("CORRELATION_ID", tt.correlId)
			}

			l.Log(tt.args.format, tt.args.level, tt.args.words...)

			if tt.useTime {
				// Use regex matching for timestamp tests
				assert.Regexp(t, tt.expect, output.String())
			} else {
				assert.Equal(t, tt.expect, output.String())
			}

			if tt.correlId != "" {
				os.Unsetenv("CORRELATION_ID")
			}
		})
	}
}

func TestCmdLogger_Debug(t *testing.T) {
	var output bytes.Buffer
	l := &CmdLogger{writer: &output}

	tests := []struct {
		name     string
		format   string
		args     []interface{}
		expected string
	}{
		{
			name:     "simple debug message",
			format:   "debug: %s",
			args:     []interface{}{"test"},
			expected: "\x1b[36mdebug: test\x1b[0m\n",
		},
		{
			name:     "debug without args",
			format:   "debug message",
			args:     nil,
			expected: "\x1b[36mdebug message\x1b[0m\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output.Reset()
			l.Debug(tt.format, tt.args...)
			assert.Equal(t, tt.expected, output.String())
		})
	}
}

func TestCmdLogger_Error(t *testing.T) {
	var output bytes.Buffer
	l := &CmdLogger{writer: &output}

	tests := []struct {
		name     string
		format   string
		args     []interface{}
		expected string
	}{
		{
			name:     "simple error message",
			format:   "error: %s",
			args:     []interface{}{"test"},
			expected: "\x1b[31merror: test\x1b[0m\n",
		},
		{
			name:     "error without args",
			format:   "error message",
			args:     nil,
			expected: "\x1b[31merror message\x1b[0m\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output.Reset()
			l.Error(tt.format, tt.args...)
			assert.Equal(t, tt.expected, output.String())
		})
	}
}

func TestCmdLogger_Warning(t *testing.T) {
	var output bytes.Buffer
	l := &CmdLogger{writer: &output}

	tests := []struct {
		name     string
		format   string
		args     []interface{}
		expected string
	}{
		{
			name:     "simple warning message",
			format:   "warning: %s",
			args:     []interface{}{"test"},
			expected: "\x1b[33mwarning: test\x1b[0m\n",
		},
		{
			name:     "warning without args",
			format:   "warning message",
			args:     nil,
			expected: "\x1b[33mwarning message\x1b[0m\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output.Reset()
			l.Warn(tt.format, tt.args...)
			assert.Equal(t, tt.expected, output.String())
		})
	}
}

func TestCmdLogger_Info(t *testing.T) {
	var output bytes.Buffer
	l := &CmdLogger{writer: &output}

	tests := []struct {
		name     string
		format   string
		args     []interface{}
		expected string
	}{
		{
			name:     "simple info message",
			format:   "info: %s",
			args:     []interface{}{"test"},
			expected: "\x1b[0minfo: test\x1b[0m\n",
		},
		{
			name:     "info without args",
			format:   "info message",
			args:     nil,
			expected: "\x1b[0minfo message\x1b[0m\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output.Reset()
			l.Info(tt.format, tt.args...)
			assert.Equal(t, tt.expected, output.String())
		})
	}
}

func TestCmdLogger_IsTimestampEnabled(t *testing.T) {
	tests := []struct {
		name     string
		l        *CmdLogger
		expected bool
	}{
		{
			name:     "timestamp enabled",
			l:        &CmdLogger{useTimestamp: true},
			expected: true,
		},
		{
			name:     "timestamp disabled",
			l:        &CmdLogger{useTimestamp: false},
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.l.IsTimestampEnabled()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCmdLogger_LogIcon(t *testing.T) {
	type args struct {
		icon   LoggerIcon
		format string
		level  Level
		words  []interface{}
	}
	tests := []struct {
		name        string
		args        args
		expect      string
		useIcons    bool
		useTime     bool
		useCorrelId bool
		correlId    string
	}{
		{
			name: "info with custom icon",
			args: args{
				icon:   "🔍",
				format: "%s",
				level:  Info,
				words: []interface{}{
					"test message",
				},
			},
			useIcons: true,
			expect:   "\x1b[0m🔍 test message\x1b[0m\n",
		},
		{
			name: "error with custom icon",
			args: args{
				icon:   "❌",
				format: "%s",
				level:  Error,
				words: []interface{}{
					"error message",
				},
			},
			useIcons: true,
			expect:   "\x1b[31m❌ error message\x1b[0m\n",
		},
		{
			name: "warning with custom icon and timestamp",
			args: args{
				icon:   "⚠️",
				format: "%s",
				level:  Warning,
				words: []interface{}{
					"warning message",
				},
			},
			useIcons: true,
			useTime:  true,
			expect:   `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z ⚠️ warning message\x1b\[0m\n`,
		},
		{
			name: "debug with custom icon and correlation id",
			args: args{
				icon:   "🐛",
				format: "%s",
				level:  Debug,
				words: []interface{}{
					"debug message",
				},
			},
			useIcons:    true,
			useCorrelId: true,
			correlId:    "test-123",
			expect:      "\x1b[36m[test-123] 🐛 debug message\x1b[0m\n",
		},
		{
			name: "trace with custom icon and correlation id",
			args: args{
				icon:   "🐛",
				format: "%s",
				level:  Trace,
				words: []interface{}{
					"trace message",
				},
			},
			useIcons:    true,
			useCorrelId: true,
			correlId:    "test-123",
			expect:      "\x1b[37m[test-123] 🐛 trace message\x1b[0m\n",
		},
		{
			name: "icons disabled",
			args: args{
				icon:   "📝",
				format: "%s",
				level:  Info,
				words: []interface{}{
					"no icon",
				},
			},
			useIcons: false,
			expect:   "\x1b[0mno icon\x1b[0m\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output bytes.Buffer
			l := &CmdLogger{
				writer: &output,
			}

			l.UseIcons(tt.useIcons)
			l.UseTimestamp(tt.useTime)
			l.UseCorrelationId(tt.useCorrelId)
			if tt.correlId != "" {
				os.Setenv("CORRELATION_ID", tt.correlId)
				defer os.Unsetenv("CORRELATION_ID")
			}

			l.LogIcon(tt.args.icon, tt.args.format, tt.args.level, tt.args.words...)

			if tt.useTime {
				assert.Regexp(t, tt.expect, output.String())
			} else {
				assert.Equal(t, tt.expect, output.String())
			}
		})
	}
}

func TestCmdLogger_LogHighlight(t *testing.T) {
	type args struct {
		format         string
		level          Level
		highlightColor strcolor.ColorCode
		words          []interface{}
	}
	tests := []struct {
		name        string
		args        args
		expect      string
		useTime     bool
		useCorrelId bool
		correlId    string
	}{
		{
			name: "info with highlighted words",
			args: args{
				format:         "Test message: %s and %s",
				level:          Info,
				highlightColor: strcolor.Red,
				words: []interface{}{
					"highlighted",
					"colored",
				},
			},
			expect: "\x1b[0mTest message: \x1b[31mhighlighted\x1b[0m and \x1b[31mcolored\x1b[0m\x1b[0m\n",
		},
		{
			name: "error with highlighted word",
			args: args{
				format:         "Error: %s",
				level:          Error,
				highlightColor: strcolor.Blue,
				words: []interface{}{
					"critical",
				},
			},
			expect: "\x1b[31mError: \x1b[34mcritical\x1b[0m\x1b[0m\n",
		},
		{
			name: "warning with timestamp and highlight",
			args: args{
				format:         "Warning: %s",
				level:          Warning,
				highlightColor: strcolor.Green,
				words: []interface{}{
					"important",
				},
			},
			useTime: true,
			expect:  `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z Warning: \x1b\[32mimportant\x1b\[0m\x1b\[0m\n`,
		},
		{
			name: "debug with correlation id and highlight",
			args: args{
				format:         "Debug: %s",
				level:          Debug,
				highlightColor: strcolor.Yellow,
				words: []interface{}{
					"testing",
				},
			},
			useCorrelId: true,
			correlId:    "test-123",
			expect:      "\x1b[36m[test-123] Debug: \x1b[33mtesting\x1b[0m\x1b[0m\n",
		},
		{
			name: "trace with correlation id and highlight",
			args: args{
				format:         "Trace: %s",
				level:          Trace,
				highlightColor: strcolor.Yellow,
				words: []interface{}{
					"testing",
				},
			},
			useCorrelId: true,
			correlId:    "test-123",
			expect:      "\x1b[37m[test-123] Trace: \x1b[33mtesting\x1b[0m\x1b[0m\n",
		},
		{
			name: "no words to highlight",
			args: args{
				format:         "Plain message",
				level:          Info,
				highlightColor: strcolor.Red,
				words:          []interface{}{},
			},
			expect: "\x1b[0mPlain message\x1b[0m\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output bytes.Buffer
			l := &CmdLogger{
				writer: &output,
			}

			l.UseTimestamp(tt.useTime)
			l.UseCorrelationId(tt.useCorrelId)
			if tt.correlId != "" {
				os.Setenv("CORRELATION_ID", tt.correlId)
				defer os.Unsetenv("CORRELATION_ID")
			}

			l.LogHighlight(tt.args.format, tt.args.level, tt.args.highlightColor, tt.args.words...)

			if tt.useTime {
				assert.Regexp(t, tt.expect, output.String())
			} else {
				assert.Equal(t, tt.expect, output.String())
			}
		})
	}
}

func TestCmdLogger_Success(t *testing.T) {
	var output bytes.Buffer
	l := &CmdLogger{writer: &output}

	tests := []struct {
		name        string
		format      string
		args        []interface{}
		useIcons    bool
		useTime     bool
		useCorrelId bool
		correlId    string
		expected    string
	}{
		{
			name:     "simple success message",
			format:   "success: %s",
			args:     []interface{}{"test"},
			useIcons: false,
			expected: "\x1b[32msuccess: test\x1b[0m\n",
		},
		{
			name:     "success with icon",
			format:   "success: %s",
			args:     []interface{}{"test"},
			useIcons: true,
			expected: "\x1b[32m👍 success: test\x1b[0m\n",
		},
		{
			name:     "success without args",
			format:   "success message",
			args:     nil,
			useIcons: false,
			expected: "\x1b[32msuccess message\x1b[0m\n",
		},
		{
			name:        "success with correlation id",
			format:      "success: %s",
			args:        []interface{}{"test"},
			useIcons:    true,
			useCorrelId: true,
			correlId:    "test-123",
			expected:    "\x1b[32m[test-123] 👍 success: test\x1b[0m\n",
		},
		{
			name:     "success with timestamp",
			format:   "success: %s",
			args:     []interface{}{"test"},
			useIcons: true,
			useTime:  true,
			expected: `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z 👍 success: test\x1b\[0m\n`,
		},
		{
			name:        "success with all features",
			format:      "success: %s",
			args:        []interface{}{"test"},
			useIcons:    true,
			useTime:     true,
			useCorrelId: true,
			correlId:    "test-123",
			expected:    `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z \[test-123\] 👍 success: test\x1b\[0m\n`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output.Reset()
			l.UseIcons(tt.useIcons)
			l.UseTimestamp(tt.useTime)
			l.UseCorrelationId(tt.useCorrelId)

			if tt.correlId != "" {
				os.Setenv("CORRELATION_ID", tt.correlId)
				defer os.Unsetenv("CORRELATION_ID")
			}

			l.Success(tt.format, tt.args...)

			if tt.useTime {
				assert.Regexp(t, tt.expected, output.String())
			} else {
				assert.Equal(t, tt.expected, output.String())
			}
		})
	}
}

func TestCmdLogger_Command(t *testing.T) {
	var output bytes.Buffer
	l := &CmdLogger{writer: &output}

	tests := []struct {
		name        string
		format      string
		args        []interface{}
		useIcons    bool
		useTime     bool
		useCorrelId bool
		correlId    string
		expected    string
	}{
		{
			name:     "simple command message",
			format:   "command: %s",
			args:     []interface{}{"test"},
			useIcons: false,
			expected: "\x1b[35mcommand: test\x1b[0m\n",
		},
		{
			name:     "command with icon",
			format:   "command: %s",
			args:     []interface{}{"test"},
			useIcons: true,
			expected: "\x1b[35m🔧 command: test\x1b[0m\n",
		},
		{
			name:     "command without args",
			format:   "command message",
			args:     nil,
			useIcons: false,
			expected: "\x1b[35mcommand message\x1b[0m\n",
		},
		{
			name:        "command with correlation id",
			format:      "command: %s",
			args:        []interface{}{"test"},
			useIcons:    true,
			useCorrelId: true,
			correlId:    "test-123",
			expected:    "\x1b[35m[test-123] 🔧 command: test\x1b[0m\n",
		},
		{
			name:     "command with timestamp",
			format:   "command: %s",
			args:     []interface{}{"test"},
			useIcons: true,
			useTime:  true,
			expected: `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z 🔧 command: test\x1b\[0m\n`,
		},
		{
			name:        "command with all features",
			format:      "command: %s",
			args:        []interface{}{"test"},
			useIcons:    true,
			useTime:     true,
			useCorrelId: true,
			correlId:    "test-123",
			expected:    `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z \[test-123\] 🔧 command: test\x1b\[0m\n`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output.Reset()
			l.UseIcons(tt.useIcons)
			l.UseTimestamp(tt.useTime)
			l.UseCorrelationId(tt.useCorrelId)

			if tt.correlId != "" {
				os.Setenv("CORRELATION_ID", tt.correlId)
				defer os.Unsetenv("CORRELATION_ID")
			}

			l.Command(tt.format, tt.args...)

			if tt.useTime {
				assert.Regexp(t, tt.expected, output.String())
			} else {
				assert.Equal(t, tt.expected, output.String())
			}
		})
	}
}

func TestCmdLogger_Disabled(t *testing.T) {
	var output bytes.Buffer
	l := &CmdLogger{writer: &output}

	tests := []struct {
		name        string
		format      string
		args        []interface{}
		useIcons    bool
		useTime     bool
		useCorrelId bool
		correlId    string
		expected    string
	}{
		{
			name:     "simple disabled message",
			format:   "disabled: %s",
			args:     []interface{}{"test"},
			useIcons: false,
			expected: "\x1b[90mdisabled: test\x1b[0m\n",
		},
		{
			name:     "disabled with icon",
			format:   "disabled: %s",
			args:     []interface{}{"test"},
			useIcons: true,
			expected: "\x1b[90m◾ disabled: test\x1b[0m\n",
		},
		{
			name:     "disabled without args",
			format:   "disabled message",
			args:     nil,
			useIcons: false,
			expected: "\x1b[90mdisabled message\x1b[0m\n",
		},
		{
			name:        "disabled with correlation id",
			format:      "disabled: %s",
			args:        []interface{}{"test"},
			useIcons:    true,
			useCorrelId: true,
			correlId:    "test-123",
			expected:    "\x1b[90m[test-123] ◾ disabled: test\x1b[0m\n",
		},
		{
			name:     "disabled with timestamp",
			format:   "disabled: %s",
			args:     []interface{}{"test"},
			useIcons: true,
			useTime:  true,
			expected: `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z ◾ disabled: test\x1b\[0m\n`,
		},
		{
			name:        "disabled with all features",
			format:      "disabled: %s",
			args:        []interface{}{"test"},
			useIcons:    true,
			useTime:     true,
			useCorrelId: true,
			correlId:    "test-123",
			expected:    `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z \[test-123\] ◾ disabled: test\x1b\[0m\n`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output.Reset()
			l.UseIcons(tt.useIcons)
			l.UseTimestamp(tt.useTime)
			l.UseCorrelationId(tt.useCorrelId)

			if tt.correlId != "" {
				os.Setenv("CORRELATION_ID", tt.correlId)
				defer os.Unsetenv("CORRELATION_ID")
			}

			l.Disabled(tt.format, tt.args...)

			if tt.useTime {
				assert.Regexp(t, tt.expected, output.String())
			} else {
				assert.Equal(t, tt.expected, output.String())
			}
		})
	}
}

func TestCmdLogger_Notice(t *testing.T) {
	var output bytes.Buffer
	l := &CmdLogger{writer: &output}

	tests := []struct {
		name        string
		format      string
		args        []interface{}
		useIcons    bool
		useTime     bool
		useCorrelId bool
		correlId    string
		expected    string
	}{
		{
			name:     "simple notice message",
			format:   "notice: %s",
			args:     []interface{}{"test"},
			useIcons: false,
			expected: "\x1b[34mnotice: test\x1b[0m\n",
		},
		{
			name:     "notice with icon",
			format:   "notice: %s",
			args:     []interface{}{"test"},
			useIcons: true,
			expected: "\x1b[34m🚩 notice: test\x1b[0m\n",
		},
		{
			name:     "notice without args",
			format:   "notice message",
			args:     nil,
			useIcons: false,
			expected: "\x1b[34mnotice message\x1b[0m\n",
		},
		{
			name:        "notice with correlation id",
			format:      "notice: %s",
			args:        []interface{}{"test"},
			useIcons:    true,
			useCorrelId: true,
			correlId:    "test-123",
			expected:    "\x1b[34m[test-123] 🚩 notice: test\x1b[0m\n",
		},
		{
			name:     "notice with timestamp",
			format:   "notice: %s",
			args:     []interface{}{"test"},
			useIcons: true,
			useTime:  true,
			expected: `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z 🚩 notice: test\x1b\[0m\n`,
		},
		{
			name:        "notice with all features",
			format:      "notice: %s",
			args:        []interface{}{"test"},
			useIcons:    true,
			useTime:     true,
			useCorrelId: true,
			correlId:    "test-123",
			expected:    `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z \[test-123\] 🚩 notice: test\x1b\[0m\n`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output.Reset()
			l.UseIcons(tt.useIcons)
			l.UseTimestamp(tt.useTime)
			l.UseCorrelationId(tt.useCorrelId)

			if tt.correlId != "" {
				os.Setenv("CORRELATION_ID", tt.correlId)
				defer os.Unsetenv("CORRELATION_ID")
			}

			l.Notice(tt.format, tt.args...)

			if tt.useTime {
				assert.Regexp(t, tt.expected, output.String())
			} else {
				assert.Equal(t, tt.expected, output.String())
			}
		})
	}
}

func TestCmdLogger_Trace(t *testing.T) {
	var output bytes.Buffer
	l := &CmdLogger{writer: &output}

	tests := []struct {
		name        string
		format      string
		args        []interface{}
		useIcons    bool
		useTime     bool
		useCorrelId bool
		correlId    string
		expected    string
	}{
		{
			name:     "simple trace message",
			format:   "trace: %s",
			args:     []interface{}{"test"},
			useIcons: false,
			expected: "\x1b[37mtrace: test\x1b[0m\n",
		},
		{
			name:     "trace with icon",
			format:   "trace: %s",
			args:     []interface{}{"test"},
			useIcons: true,
			expected: "\x1b[37m💡 trace: test\x1b[0m\n",
		},
		{
			name:     "trace without args",
			format:   "trace message",
			args:     nil,
			useIcons: false,
			expected: "\x1b[37mtrace message\x1b[0m\n",
		},
		{
			name:        "trace with correlation id",
			format:      "trace: %s",
			args:        []interface{}{"test"},
			useIcons:    true,
			useCorrelId: true,
			correlId:    "test-123",
			expected:    "\x1b[37m[test-123] 💡 trace: test\x1b[0m\n",
		},
		{
			name:     "trace with timestamp",
			format:   "trace: %s",
			args:     []interface{}{"test"},
			useIcons: true,
			useTime:  true,
			expected: `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z 💡 trace: test\x1b\[0m\n`,
		},
		{
			name:        "trace with all features",
			format:      "trace: %s",
			args:        []interface{}{"test"},
			useIcons:    true,
			useTime:     true,
			useCorrelId: true,
			correlId:    "test-123",
			expected:    `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z \[test-123\] 💡 trace: test\x1b\[0m\n`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output.Reset()
			l.UseIcons(tt.useIcons)
			l.UseTimestamp(tt.useTime)
			l.UseCorrelationId(tt.useCorrelId)

			if tt.correlId != "" {
				os.Setenv("CORRELATION_ID", tt.correlId)
				defer os.Unsetenv("CORRELATION_ID")
			}

			l.Trace(tt.format, tt.args...)

			if tt.useTime {
				assert.Regexp(t, tt.expected, output.String())
			} else {
				assert.Equal(t, tt.expected, output.String())
			}
		})
	}
}

func TestCmdLogger_Exception(t *testing.T) {
	var output bytes.Buffer
	l := &CmdLogger{writer: &output}

	tests := []struct {
		name        string
		err         error
		format      string
		args        []interface{}
		useIcons    bool
		useTime     bool
		useCorrelId bool
		correlId    string
		expected    string
	}{
		{
			name:     "exception with empty format",
			err:      fmt.Errorf("test error"),
			format:   "",
			args:     nil,
			useIcons: false,
			expected: "\x1b[31mtest error\x1b[0m\n",
		},
		{
			name:     "exception with format",
			err:      fmt.Errorf("test error"),
			format:   "Operation failed",
			args:     nil,
			useIcons: false,
			expected: "\x1b[31mOperation failed, err test error\x1b[0m\n",
		},
		{
			name:     "exception with format and args",
			err:      fmt.Errorf("test error"),
			format:   "Operation %s failed",
			args:     []interface{}{"save"},
			useIcons: false,
			expected: "\x1b[31mOperation save failed, err test error\x1b[0m\n",
		},
		{
			name:     "exception with icon",
			err:      fmt.Errorf("test error"),
			format:   "Operation failed",
			args:     nil,
			useIcons: true,
			expected: "\x1b[31m🚨 Operation failed, err test error\x1b[0m\n",
		},
		{
			name:        "exception with correlation id",
			err:         fmt.Errorf("test error"),
			format:      "Operation failed",
			args:        nil,
			useIcons:    true,
			useCorrelId: true,
			correlId:    "test-123",
			expected:    "\x1b[31m[test-123] 🚨 Operation failed, err test error\x1b[0m\n",
		},
		{
			name:     "exception with timestamp",
			err:      fmt.Errorf("test error"),
			format:   "Operation failed",
			args:     nil,
			useIcons: true,
			useTime:  true,
			expected: `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z 🚨 Operation failed, err test error\x1b\[0m\n`,
		},
		{
			name:        "exception with all features",
			err:         fmt.Errorf("test error"),
			format:      "Operation %s failed",
			args:        []interface{}{"save"},
			useIcons:    true,
			useTime:     true,
			useCorrelId: true,
			correlId:    "test-123",
			expected:    `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z \[test-123\] 🚨 Operation save failed, err test error\x1b\[0m\n`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output.Reset()
			l.UseIcons(tt.useIcons)
			l.UseTimestamp(tt.useTime)
			l.UseCorrelationId(tt.useCorrelId)

			if tt.correlId != "" {
				os.Setenv("CORRELATION_ID", tt.correlId)
				defer os.Unsetenv("CORRELATION_ID")
			}

			l.Exception(tt.err, tt.format, tt.args...)

			if tt.useTime {
				assert.Regexp(t, tt.expected, output.String())
			} else {
				assert.Equal(t, tt.expected, output.String())
			}
		})
	}
}

func TestCmdLogger_LogError(t *testing.T) {
	var output bytes.Buffer
	l := &CmdLogger{writer: &output}

	tests := []struct {
		name        string
		err         error
		useIcons    bool
		useTime     bool
		useCorrelId bool
		correlId    string
		expected    string
	}{
		{
			name:     "simple error message",
			err:      fmt.Errorf("test error"),
			useIcons: false,
			expected: "\x1b[31mtest error\x1b[0m\n",
		},
		{
			name:     "error with icon",
			err:      fmt.Errorf("test error"),
			useIcons: true,
			expected: "\x1b[31m🚨 test error\x1b[0m\n",
		},
		{
			name:        "error with correlation id",
			err:         fmt.Errorf("test error"),
			useIcons:    true,
			useCorrelId: true,
			correlId:    "test-123",
			expected:    "\x1b[31m[test-123] 🚨 test error\x1b[0m\n",
		},
		{
			name:     "error with timestamp",
			err:      fmt.Errorf("test error"),
			useIcons: true,
			useTime:  true,
			expected: `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z 🚨 test error\x1b\[0m\n`,
		},
		{
			name:        "error with all features",
			err:         fmt.Errorf("test error"),
			useIcons:    true,
			useTime:     true,
			useCorrelId: true,
			correlId:    "test-123",
			expected:    `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z \[test-123\] 🚨 test error\x1b\[0m\n`,
		},
		{
			name:     "nil error",
			err:      nil,
			useIcons: true,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output.Reset()
			l.UseIcons(tt.useIcons)
			l.UseTimestamp(tt.useTime)
			l.UseCorrelationId(tt.useCorrelId)

			if tt.correlId != "" {
				os.Setenv("CORRELATION_ID", tt.correlId)
				defer os.Unsetenv("CORRELATION_ID")
			}

			l.LogError(tt.err)

			if tt.useTime {
				assert.Regexp(t, tt.expected, output.String())
			} else {
				assert.Equal(t, tt.expected, output.String())
			}
		})
	}
}

func TestCmdLogger_Fatal(t *testing.T) {
	var output bytes.Buffer
	l := &CmdLogger{writer: &output}

	tests := []struct {
		name        string
		format      string
		args        []interface{}
		useIcons    bool
		useTime     bool
		useCorrelId bool
		correlId    string
		expected    string
	}{
		{
			name:     "simple fatal message",
			format:   "fatal: %s",
			args:     []interface{}{"test"},
			useIcons: false,
			expected: "\x1b[31mfatal: test\x1b[0m\n",
		},
		{
			name:     "fatal with icon",
			format:   "fatal: %s",
			args:     []interface{}{"test"},
			useIcons: true,
			expected: "\x1b[31m🚨 fatal: test\x1b[0m\n",
		},
		{
			name:     "fatal without args",
			format:   "fatal message",
			args:     nil,
			useIcons: false,
			expected: "\x1b[31mfatal message\x1b[0m\n",
		},
		{
			name:        "fatal with correlation id",
			format:      "fatal: %s",
			args:        []interface{}{"test"},
			useIcons:    true,
			useCorrelId: true,
			correlId:    "test-123",
			expected:    "\x1b[31m[test-123] 🚨 fatal: test\x1b[0m\n",
		},
		{
			name:     "fatal with timestamp",
			format:   "fatal: %s",
			args:     []interface{}{"test"},
			useIcons: true,
			useTime:  true,
			expected: `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z 🚨 fatal: test\x1b\[0m\n`,
		},
		{
			name:        "fatal with all features",
			format:      "fatal: %s",
			args:        []interface{}{"test"},
			useIcons:    true,
			useTime:     true,
			useCorrelId: true,
			correlId:    "test-123",
			expected:    `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z \[test-123\] 🚨 fatal: test\x1b\[0m\n`,
		},
		{
			name:     "fatal with multiple arguments",
			format:   "fatal: %s %d %s",
			args:     []interface{}{"test", 42, "error"},
			useIcons: true,
			expected: "\x1b[31m🚨 fatal: test 42 error\x1b[0m\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output.Reset()
			l.UseIcons(tt.useIcons)
			l.UseTimestamp(tt.useTime)
			l.UseCorrelationId(tt.useCorrelId)

			if tt.correlId != "" {
				os.Setenv("CORRELATION_ID", tt.correlId)
				defer os.Unsetenv("CORRELATION_ID")
			}

			l.Fatal(tt.format, tt.args...)

			if tt.useTime {
				assert.Regexp(t, tt.expected, output.String())
			} else {
				assert.Equal(t, tt.expected, output.String())
			}
		})
	}
}

func TestCmdLogger_FatalError(t *testing.T) {
	var output bytes.Buffer
	l := &CmdLogger{writer: &output}

	tests := []struct {
		name        string
		err         error
		format      string
		args        []interface{}
		useIcons    bool
		useTime     bool
		useCorrelId bool
		correlId    string
		expected    string
		shouldPanic bool
	}{
		{
			name:        "fatal error with message",
			err:         fmt.Errorf("test error"),
			format:      "Operation failed",
			args:        nil,
			useIcons:    false,
			expected:    "\x1b[31mOperation failed\x1b[0m\n",
			shouldPanic: true,
		},
		{
			name:        "fatal error with format args",
			err:         fmt.Errorf("test error"),
			format:      "Operation %s failed",
			args:        []interface{}{"save"},
			useIcons:    true,
			expected:    "\x1b[31m🚨 Operation save failed\x1b[0m\n",
			shouldPanic: true,
		},
		{
			name:        "fatal error with correlation id",
			err:         fmt.Errorf("test error"),
			format:      "Operation failed",
			useIcons:    true,
			useCorrelId: true,
			correlId:    "test-123",
			expected:    "\x1b[31m[test-123] 🚨 Operation failed\x1b[0m\n",
			shouldPanic: true,
		},
		{
			name:        "fatal error with timestamp",
			err:         fmt.Errorf("test error"),
			format:      "Operation failed",
			useIcons:    true,
			useTime:     true,
			expected:    `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z 🚨 Operation failed\x1b\[0m\n`,
			shouldPanic: true,
		},
		{
			name:        "fatal error with all features",
			err:         fmt.Errorf("test error"),
			format:      "Operation %s failed",
			args:        []interface{}{"save"},
			useIcons:    true,
			useTime:     true,
			useCorrelId: true,
			correlId:    "test-123",
			expected:    `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z \[test-123\] 🚨 Operation save failed\x1b\[0m\n`,
			shouldPanic: true,
		},
		{
			name:        "nil error should not panic",
			err:         nil,
			format:      "Operation failed",
			useIcons:    true,
			expected:    "\x1b[31m🚨 Operation failed\x1b[0m\n",
			shouldPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output.Reset()
			l.UseIcons(tt.useIcons)
			l.UseTimestamp(tt.useTime)
			l.UseCorrelationId(tt.useCorrelId)

			if tt.correlId != "" {
				os.Setenv("CORRELATION_ID", tt.correlId)
				defer os.Unsetenv("CORRELATION_ID")
			}

			if tt.shouldPanic {
				defer func() {
					r := recover()
					if r == nil {
						t.Error("FatalError did not panic")
					}
					if err, ok := r.(error); !ok || err != tt.err {
						t.Errorf("Expected panic with error %v, got %v", tt.err, r)
					}
				}()
			}

			l.FatalError(tt.err, tt.format, tt.args...)

			if tt.useTime {
				assert.Regexp(t, tt.expected, output.String())
			} else {
				assert.Equal(t, tt.expected, output.String())
			}

			if tt.shouldPanic {
				t.Error("Expected panic did not occur")
			}
		})
	}
}
