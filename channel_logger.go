package log

import (
	"fmt"
	"sync"
	"time"

	strcolor "github.com/cjlapao/common-go/strcolor"
	"github.com/google/uuid"
)

type LogMessage struct {
	Level     string
	Message   string
	Timestamp time.Time
	Icon      LoggerIcon
	IsTask    bool
}

type Subscriber struct {
	id      string
	filter  func(LogMessage) bool
	channel chan LogMessage
}

// String returns a formatted string representation of the LogMessage
func (m LogMessage) String() string {
	timestamp := m.Timestamp.Format(time.RFC3339)
	if m.Icon != "" {
		return fmt.Sprintf("[%s] %s %s: %s", timestamp, m.Icon, m.Level, m.Message)
	}
	return fmt.Sprintf("[%s] %s: %s", timestamp, m.Level, m.Message)
}

// ChannelLogger Command Line Logger implementation
type ChannelLogger struct {
	useTimestamp      bool
	userCorrelationId bool
	useIcons          bool
	subscribers       []Subscriber
	channelMutex      sync.RWMutex
}

func (l *ChannelLogger) Init() Logger {
	return &ChannelLogger{
		useTimestamp:      false,
		userCorrelationId: false,
		useIcons:          false,
		subscribers:       make([]Subscriber, 0),
		channelMutex:      sync.RWMutex{},
	}
}

func (l *ChannelLogger) IsTimestampEnabled() bool {
	return l.useTimestamp
}

func (l *ChannelLogger) UseTimestamp(value bool) {
	l.useTimestamp = value
}

func (l *ChannelLogger) UseCorrelationId(value bool) {
	l.userCorrelationId = value
}

func (l *ChannelLogger) UseIcons(value bool) {
	l.useIcons = value
}

func (l *ChannelLogger) printMessage(format string, icon LoggerIcon, level string, words ...interface{}) {
	if len(l.subscribers) == 0 {
		return // Do nothing if no subscribers
	}

	if len(words) > 0 {
		format = fmt.Sprintf(format, words...)
	}

	msg := LogMessage{
		Level:     level,
		Message:   format,
		Timestamp: time.Now(),
		Icon:      icon,
	}

	if l.useIcons && icon != "" {
		msg.Message = fmt.Sprintf("%s %s", icon, msg.Message)
	}

	// Send message to all active subscribers
	l.channelMutex.RLock()
	defer l.channelMutex.RUnlock()

	for _, sub := range l.subscribers {
		if sub.filter(msg) { // Use filter instead of id
			select {
			case sub.channel <- msg:
				// Message sent successfully
			default:
				// Channel is full, skip this message for this subscriber
			}
		}
	}
}

func (l *ChannelLogger) Log(format string, level Level, words ...interface{}) {
	switch level {
	case 0:
		l.printMessage(format, "", "error", words...)
	case 1:
		l.printMessage(format, "", "warn", words...)
	case 2:
		l.printMessage(format, "", "info", words...)
	case 3:
		l.printMessage(format, "", "debug", words...)
	case 4:
		l.printMessage(format, "", "trace", words...)
	}
}

// Log Log information message
func (l *ChannelLogger) LogIcon(icon LoggerIcon, format string, level Level, words ...interface{}) {
	switch level {
	case 0:
		l.printMessage(format, icon, "error", words...)
	case 1:
		l.printMessage(format, icon, "warn", words...)
	case 2:
		l.printMessage(format, icon, "info", words...)
	case 3:
		l.printMessage(format, icon, "debug", words...)
	case 4:
		l.printMessage(format, icon, "trace", words...)
	}
}

// LogHighlight Log information message
func (l *ChannelLogger) LogHighlight(format string, level Level, highlightColor strcolor.ColorCode, words ...interface{}) {
	if len(words) > 0 {
		for i := range words {
			word := fmt.Sprintf("%v", words[i])
			words[i] = GetColorString(ColorCode(highlightColor), word)
		}
	}

	switch level {
	case 0:
		l.printMessage(format, "", "error", words...)
	case 1:
		l.printMessage(format, "", "warn", words...)
	case 2:
		l.printMessage(format, "", "info", words...)
	case 3:
		l.printMessage(format, "", "debug", words...)
	case 4:
		l.printMessage(format, "", "trace", words...)
	}
}

// Info log information message
func (l *ChannelLogger) Info(format string, words ...interface{}) {
	l.printMessage(format, IconInfo, "info", words...)
}

// Success log message
func (l *ChannelLogger) Success(format string, words ...interface{}) {
	l.printMessage(format, IconThumbsUp, "success", words...)
}

// Warn log message
func (l *ChannelLogger) Warn(format string, words ...interface{}) {
	l.printMessage(format, IconWarning, "warn", words...)
}

// Command log message
func (l *ChannelLogger) Command(format string, words ...interface{}) {
	l.printMessage(format, IconWrench, "command", words...)
}

// Disabled log message
func (l *ChannelLogger) Disabled(format string, words ...interface{}) {
	l.printMessage(format, IconBlackSquare, "disabled", words...)
}

// Notice log message
func (l *ChannelLogger) Notice(format string, words ...interface{}) {
	l.printMessage(format, IconFlag, "notice", words...)
}

// Debug log message
func (l *ChannelLogger) Debug(format string, words ...interface{}) {
	l.printMessage(format, IconFire, "debug", words...)
}

// Trace log message
func (l *ChannelLogger) Trace(format string, words ...interface{}) {
	l.printMessage(format, IconBulb, "trace", words...)
}

// Error log message
func (l *ChannelLogger) Error(format string, words ...interface{}) {
	l.printMessage(format, IconRevolvingLight, "error", words...)
}

// Error log message
func (l *ChannelLogger) Exception(err error, format string, words ...interface{}) {
	if format == "" {
		format = err.Error()
	} else {
		format = format + ", err " + err.Error()
	}
	l.printMessage(format, IconRevolvingLight, "error", words...)
}

// LogError log message
func (l *ChannelLogger) LogError(message error) {
	if message != nil {
		l.printMessage(message.Error(), IconRevolvingLight, "error")
	}
}

// Fatal log message
func (l *ChannelLogger) Fatal(format string, words ...interface{}) {
	l.printMessage(format, IconRevolvingLight, "error", words...)
}

// FatalError log message
func (l *ChannelLogger) FatalError(e error, format string, words ...interface{}) {
	l.Error(format, words...)
	if e != nil {
		panic(e)
	}
}

// Add Subscribe method to ChannelLogger
func (l *ChannelLogger) Subscribe(id string, callback func(LogMessage) bool) (string, chan LogMessage) {
	l.channelMutex.Lock()
	defer l.channelMutex.Unlock()

	if id == "" {
		id = uuid.New().String()
	}

	// Generate unique ID for this subscription
	subID := fmt.Sprintf("sub_%s", id)
	ch := make(chan LogMessage, 100)

	// Check if subscription ID already exists
	for _, sub := range l.subscribers {
		if sub.id == subID {
			return subID, sub.channel
		}
	}

	// Each subscription will get its own channel
	l.subscribers = append(l.subscribers, Subscriber{
		id:      subID,
		filter:  callback,
		channel: ch,
	})
	return subID, ch
}

// Unsubscribe removes a subscription and closes its channel
func (l *ChannelLogger) Unsubscribe(subscriptionID string) bool {
	l.channelMutex.Lock()
	defer l.channelMutex.Unlock()

	// Find and remove the subscription
	for i, sub := range l.subscribers {
		if sub.id == subscriptionID {
			// Close the channel
			close(sub.channel)

			// Remove the subscriber from the slice
			l.subscribers = append(l.subscribers[:i], l.subscribers[i+1:]...)
			return true
		}
	}
	return false
}

// Update Channel method to handle the new return signature
func (l *ChannelLogger) Channel() (string, chan LogMessage) {
	return l.Subscribe("", func(LogMessage) bool { return true })
}

// Update Close method to handle local subscribers
func (l *ChannelLogger) Close() {
	l.channelMutex.Lock()
	defer l.channelMutex.Unlock()

	for _, sub := range l.subscribers {
		close(sub.channel)
	}
	l.subscribers = nil
}
