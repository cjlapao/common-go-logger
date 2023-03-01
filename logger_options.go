package log

// LogOptions Definition
type LoggerOptions int64

const (
	WithTimestamp LoggerOptions = iota
	WithCorrelationId
)
