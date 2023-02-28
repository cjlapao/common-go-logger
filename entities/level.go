package entities

// Level Entity
type Level int

// LogLevel Enum Definition
const (
	Error Level = iota
	Warning
	Info
	Debug
	Trace
)
