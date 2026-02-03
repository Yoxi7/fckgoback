package archarchive

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

func Ask(items []string, message string) (string, error) {
	if len(items) == 0 {
		return "", fmt.Errorf("no items to choose from")
	}

	var result string
	prompt := &survey.Select{
		Message: message,
		Options: items,
	}

	err := survey.AskOne(prompt, &result)
	if err != nil {
		return "", err
	}

	return result, nil
}

