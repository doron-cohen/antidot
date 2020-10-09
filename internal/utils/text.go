package utils

import (
	"fmt"
	"strings"
)

func IndentLines(lines string) string {
	var builder strings.Builder
	for _, line := range strings.Split(lines, "\n") {
		builder.WriteString(fmt.Sprintf("  %s\n", line))
	}
	return builder.String()
}
