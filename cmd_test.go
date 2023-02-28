package log

import (
	"bytes"
	"testing"

	"github.com/cjlapao/common-go-logger/entities"
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
		level  entities.Level
		words  []interface{}
	}
	tests := []struct {
		name   string
		args   args
		expect string
	}{
		{
			name: "info",
			args: args{
				format: "%s",
				level:  entities.Info,
				words: []interface{}{
					"I am",
				},
			},
			expect: "I am\x1b[0m\n",
		},
		{
			name: "info, no words",
			args: args{
				format: "just me",
				level:  entities.Info,
			},
			expect: "just me\x1b[0m\n",
		},
		{
			name: "info, empty words",
			args: args{
				format: "just me",
				level:  entities.Info,
				words:  []interface{}{},
			},
			expect: "just me\x1b[0m\n",
		},
		{
			name: "error",
			args: args{
				format: "%s",
				level:  entities.Error,
				words: []interface{}{
					"some error",
				},
			},
			expect: "some error\x1b[0m\n",
		},
		{
			name: "warn",
			args: args{
				format: "warning: %s",
				level:  entities.Error,
				words: []interface{}{
					"some warning",
				},
			},
			expect: "warning: some warning\x1b[0m\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output bytes.Buffer
			l := &CmdLogger{
				writer: &output,
			}

			l.Log(tt.args.format, tt.args.level, tt.args.words...)
			if tt.expect != output.String() {
				t.Errorf("got %s but expected %s", output.String(), tt.expect)
			}
		})
	}
}
