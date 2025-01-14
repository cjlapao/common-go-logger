# Common Logger

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![PR Builder](https://github.com/cjlapao/common-go-logger/actions/workflows/pr.yml/badge.svg)](https://github.com/cjlapao/common-go-logger/actions/workflows/pr.yml)
[![CI Release](https://github.com/cjlapao/common-go-logger/actions/workflows/ci.yml/badge.svg)](https://github.com/cjlapao/common-go-logger/actions/workflows/ci.yml)
  
This is a common logger for Go projects. It provides flexible logging capabilities with support for multiple outputs, colors, icons, and various logging levels. It can log to files, console, or both, with different logging levels for each output.

## Features

- Multiple output targets (Console, File, Channel)
- Multiple log levels (Error, Warning, Info, Debug, Trace)
- Colorized console output
- Unicode icon support
- Timestamp support
- Correlation ID support
- ADO pipeline integration
- Thread-safe logging
- Customizable message highlighting

## Installation

```bash
go get github.com/cjlapao/common-go-logger
```

## Usage

```go
logger := log.NewCmdLogger()
// Basic logging
logger.Info("Hello, %s!", "World")
logger.Error("Something went wrong: %s", err)
logger.Debug("Debug message")
logger.Warn("Warning message")
logger.Trace("Trace message")
```

### Enable Features

```go
logger := log.NewCmdLogger()
// Enable timestamps
logger.UseTimestamp(true)
// Enable icons
logger.UseIcons(true)
// Enable correlation ID
logger.UseCorrelationId(true)
```

### Log Levels and Special Functions

```go
// Different logging levels
logger.Info("Information message")
logger.Success("Operation succeeded")
logger.Warn("Warning message")
logger.Error("Error occurred")
logger.Debug("Debug information")
logger.Trace("Trace message")
logger.Fatal("Fatal error")
// Special logging functions
logger.Command("Executing command")
logger.Disabled("Feature is disabled")
logger.Notice("System notification")
logger.Exception(err, "Operation failed")
logger.FatalError(err, "Critical failure")
```

### Channel Logger and Subscribers

```go
// Create a logger with channel support
service := log.New()
service.AddChannelLogger()

// Define subscribers
sub1 := func(msg LogMessage) {
    fmt.Printf("Subscriber 1: %s\n", msg.Message)
}
sub2 := func(msg LogMessage) {
    fmt.Printf("Subscriber 2: %s\n", msg.Message)
}

// Add subscribers and get their IDs
sub1ID := service.OnMessage(sub1)
sub2ID := service.OnMessage(sub2)

// Use the logger
service.Info("This message goes to all subscribers")

// Remove a specific subscriber
service.RemoveSubscriber(sub1ID)

// This message only goes to sub2
service.Info("This message only goes to subscriber 2")

// Clean up when done
service.CloseSubscribers()
```

## Environment Variables

- `CORRELATION_ID`: Sets the correlation ID for log tracking

## Output Colors

The logger uses different colors for various log levels in console output:

- Error: Red
- Warning: Yellow
- Success: Green
- Info: Default (White)
- Debug: Cyan
- Trace: Light Gray
- Notice: Blue
- Command: Magenta
- Disabled: Dark Gray

## Icons

List of available icons:

| Icon | Description |
|:----|:------------|
|IconHammer |           :hammer: |
|IconFire |             :fire: |
|IconWrench |           :wrench: |
|IconKey |              :key: |
|IconLock |             :lock: |
|IconOpenLock |         :unlock: |
|IconBell |             :bell: |
|IconMagnifyingGlass |  :mag: |
|IconBook |             :book: |
|IconBulb |             :bulb: |
|IconBomb |             :bomb: |
|IconLargeWhiteSquare | :white_large_square: |
|IconCircle |           :o: |
|IconWarning |          :warning: |
|IconRightArrow |       :arrow_right: |
|IconHourGlass |        :hourglass: |
|IconInfo |             :information_source: |
|IconFlag |             :triangular_flag_on_post: |
|IconRocket |           :rocket: |
|IconCheckMark |        :white_check_mark: |
|IconCrossMark |        :x: |
|IconRevolvingLight |   :rotating_light: |
|IconBlackSquare |      :black_square: |
|IconFolder |           :file_folder: |
|IconClipboard |        :clipboard: |
|IconRightwardsArrow |  :arrow_right: |
|IconExclamationMark |  :exclamation: |
|IconAsterisk |         :asterisk: |
|IconRightHand |        :point_right: |
|IconCheckbox |         :ballot_box_with_check: |
|IconToilet |           :toilet: |
|IconThumbsUp |         :thumbsup: |
|IconThumbDown |        :thumbsdown: |
|IconPage |             :page_facing_up: |

## Default Icon Usage

When icons are enabled, these icons are used by default for different message types:

- Info: ℹ️ (IconInfo)
- Success: 👍 (IconThumbsUp)
- Warning: ⚠️ (IconWarning)
- Error/Fatal: 🚨 (IconRevolvingLight)
- Debug: 🔥 (IconFire)
- Trace: 💡 (IconBulb)
- Command: 🔧 (IconWrench)
- Disabled: ◾ (IconBlackSquare)
- Notice: 🚩 (IconFlag)

## Thread Safety

The logger is designed to be thread-safe and can be safely used from multiple goroutines.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
