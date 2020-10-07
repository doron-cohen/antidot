package tui

import (
	"fmt"

	"github.com/wzshiming/ctc"
)

type Style ctc.Color

const (
	reset ctc.Color = ctc.Reset

	Red     ctc.Color = ctc.ForegroundBrightRed
	Cyan    ctc.Color = ctc.ForegroundBrightCyan
	Blue    ctc.Color = ctc.ForegroundBrightBlue
	Magenta ctc.Color = ctc.ForegroundBrightMagenta
	Green   ctc.Color = ctc.ForegroundBrightGreen
	Gray    ctc.Color = ctc.ForegroundBrightBlack
)

func ApplyStyle(style ctc.Color, text string) string {
	// TODO: respect NO_COLOR
	return fmt.Sprintf("%s%s%s", style, text, reset)
}

func ApplyStylef(style ctc.Color, template string, components ...interface{}) string {
	text := fmt.Sprintf(template, components...)
	return ApplyStyle(style, text)
}
