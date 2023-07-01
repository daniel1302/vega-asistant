package uilib

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tcnksm/go-input"

	"github.com/daniel1302/vega-asistant/types"
)

type YesNoAnswer string

const (
	AnswerYes YesNoAnswer = "Yes"
	AnswerNo  YesNoAnswer = "No"
)

func AskPath(ui *input.UI, name, defaultValue string) (string, error) {
	response, err := ui.Ask(fmt.Sprintf("What is your %s", name), &input.Options{
		Default:  defaultValue,
		Required: true,
		Loop:     true,
		ValidateFunc: func(s string) error {
			return nil
		},
	})
	if err != nil {
		return "", types.NewInputError(err)
	}

	return response, nil
}

func AskRemoveExistingFile(
	ui *input.UI,
	filePath string,
	defaultAnswer YesNoAnswer,
) (YesNoAnswer, error) {
	return AskYesNo(
		ui,
		fmt.Sprintf("File %s exists. Do you want to remove it?", filePath),
		defaultAnswer,
	)
}

func AskString(
	ui *input.UI,
	question string,
	defaultAnswer string,
	validateFunc input.ValidateFunc,
) (string, error) {
	answer, err := ui.Ask(question, &input.Options{
		Default:      defaultAnswer,
		Required:     true,
		Loop:         true,
		ValidateFunc: validateFunc,
	})
	if err != nil {
		return "", fmt.Errorf("failed to ask for '%s': %w", question, err)
	}

	return answer, nil
}

func AskInt(ui *input.UI, question string, defaultAnswer int) (int, error) {
	answer, err := ui.Ask(question, &input.Options{
		Default:  fmt.Sprintf("%d", defaultAnswer),
		Required: true,
		Loop:     true,
		ValidateFunc: func(s string) error {
			if _, err := strconv.Atoi(s); err != nil {
				return fmt.Errorf("invalid int provided(%s): %w", s, err)
			}

			return nil
		},
	})
	if err != nil {
		return 0, fmt.Errorf("failed to ask for '%s': %w", question, err)
	}

	answerInt, err := strconv.Atoi(answer)
	if err != nil {
		return 0, fmt.Errorf("failed to convert answer into string for '%s': %w", question, err)
	}

	return answerInt, nil
}

func AskYesNo(ui *input.UI, question string, defaultAnswer YesNoAnswer) (YesNoAnswer, error) {
	answer, err := ui.Ask(question,
		&input.Options{
			Default:  string(defaultAnswer),
			Required: true,
			Loop:     true,
			ValidateFunc: func(s string) error {
				normalizedResponse := strings.ToLower(s)
				if normalizedResponse != "yes" && normalizedResponse != "no" {
					return fmt.Errorf("invalid response; got %s, expected Yes or No", s)
				}
				return nil
			},
		},
	)
	if err != nil {
		return defaultAnswer, fmt.Errorf("failed to ask for yes/no: %w", err)
	}

	if strings.ToLower(answer) == "yes" {
		return AnswerYes, nil
	}

	return AnswerNo, nil
}
