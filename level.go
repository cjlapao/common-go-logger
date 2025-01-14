package log

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

func (l Level) String() string {
	return []string{"error", "warning", "info", "debug", "trace"}[l]
}
