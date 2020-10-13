package tui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Confirm(text string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s - enter 'yes' to proceed: ", text)
	answer, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	sanitized := strings.TrimSpace(strings.ToLower(answer))
	return sanitized == "yes", nil
}
