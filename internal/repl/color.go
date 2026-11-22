package repl

import (
	"fmt"
	"strings"
)

const (
	Reset = "\033[0m"
	Bold = "\033[1m"
	Dim = "\033[2m"
	Italic = "\033[3m"
	Underline = "\033[4m"

	Red = "\033[31m"
	Green = "\033[32m"
	Yellow = "\033[33m"
	Blue = "\033[34m"
	Magenta = "\033[35m"
	Cyan = "\033[36m"
	White = "\033[37m"

	BoldRed = "\033[1;31m"
	BoldGreen = "\033[1;32m"
	BoldYellow = "\033[1;33m"
	BoldBlue = "\033[1;34m"
	BoldMagenta = "\033[1;35m"
	BoldCyan = "\033[1;36m"
	BoldWhite = "\033[1;37m"
)

func Colorize(color, text string) string {
	return color + text + Reset
}

func FormatHPBar(current, max, width int) string {
	if max <= 0 {
		max = 1
	}

	if current < 0 {
		current = 0
	}

	if current > max {
		current = max
	}

	filledLen := (current * width) / max
	emptyLen := width - filledLen

	var barColor string
	pct := float64(current) / float64(max)
	switch {
	case pct > 0.50:
		barColor = Green
	case pct > 0.25:
		barColor = Yellow
	default:
		barColor = Red
	}

	bar := strings.Repeat("█", filledLen) + strings.Repeat("░", emptyLen)
	return fmt.Sprintf("%s[%s]%s %d/%d HP", barColor, bar, Reset, current, max)
}

func FormatDanger(danger int) string {
	switch {
	case danger <= 0:
		return Colorize(Green, "0 (Peaceful)")
	case danger == 1:
		return Colorize(Yellow, "1 (Unsettling)")
	case danger == 2:
		return Colorize(Yellow, "2 (Dangerous)")
	case danger == 3:
		return Colorize(Red, "3 (Deadly)")
	default:
		return Colorize(BoldRed, fmt.Sprintf("%d (Perilous)", danger))
	}
}
