package log

import "github.com/fatih/color"

const (
	LOG_LEVEL string = "LOG_LEVEL"
)

// Logger Ansi Colors
const (
	SuccessColor  = color.FgGreen
	InfoColor     = color.FgHiWhite
	NoticeColor   = color.FgHiCyan
	WarningColor  = color.FgYellow
	ErrorColor    = color.FgRed
	DebugColor    = color.FgMagenta
	TraceColor    = color.FgHiMagenta
	CommandColor  = color.FgBlue
	DisabledColor = color.FgHiBlack
)
