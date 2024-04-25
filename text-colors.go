package main

import (
	"fmt"
	"strings"
)

type ColorCode int

const (
	NoColor ColorCode = iota // - 1
	Black
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

type TextColor struct {
	Color   ColorCode
	BGColor ColorCode
	Bold    bool
	Bright  bool
	Blink   bool
}

func coloredText0(text string, fg, bg ColorCode, bold, bright, blink bool) string {
	if fg == NoColor && bg == NoColor && !bold {
		return text
	}
	var shift ColorCode = 0
	if bright {
		shift = 60
	}
	var codes []string
	if fg != NoColor {
		codes = append(codes, fmt.Sprintf("%d", (shift+30)+fg-1))
	}
	if bg != NoColor {
		codes = append(codes, fmt.Sprintf("%d", (shift+40)+bg-1))
	}
	if bold {
		codes = append(codes, "1")
	}
	if blink {
		codes = append(codes, "5")
	}
	codesStr := strings.Join(codes, ";")
	// fmt.Println("codes:", codesStr)
	return fmt.Sprintf("\033[%sm%s\033[0m", codesStr, text)
}

func coloredText(text string, color TextColor) string {
	if color.Color == NoColor && color.BGColor == NoColor && !color.Bold && !color.Bright && !color.Blink {
		return text
	}
	var shift ColorCode = 0
	if color.Bright {
		shift = 60
	}
	var codes []string
	if color.Color != NoColor {
		codes = append(codes, fmt.Sprintf("%d", (shift+30)+color.Color-1))
	}
	if color.BGColor != NoColor {
		codes = append(codes, fmt.Sprintf("%d", (shift+40)+color.BGColor-1))
	}
	if color.Bold {
		codes = append(codes, "1")
	}
	if color.Blink {
		codes = append(codes, "5")
	}
	codesStr := strings.Join(codes, ";")
	// fmt.Println("codes:", codesStr)
	return fmt.Sprintf("\033[%sm%s\033[0m", codesStr, text)
}
