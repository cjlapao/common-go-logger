package log

import "fmt"

type ColorCode int

const (
	Black ColorCode = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
	BrightBlack ColorCode = iota + 82
	BrightRed
	BrightGreen
	BrightYellow
	BrightBlue
	BrightMagenta
	BrightCyan
	BrightWhite
)

func GetColorString(colorCode ColorCode, words ...string) string {
	var builder string
	for _, m := range words {
		if len(builder) > 0 {
			builder += " "
		}
		builder += m
	}

	return fmt.Sprintf("\033[%vm%v\033[0m", fmt.Sprint(colorCode), builder)
}
