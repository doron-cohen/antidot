package tui

import "github.com/AlecAivazis/survey/v2"

func Confirm(text string) bool {
	var answer bool
	prompt := &survey.Confirm{
		Message: text,
	}
	survey.AskOne(prompt, &answer)
	return answer
}
