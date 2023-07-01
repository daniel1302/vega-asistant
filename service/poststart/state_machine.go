package poststart

import (
	"fmt"
	"path/filepath"

	"github.com/tcnksm/go-input"

	"github.com/daniel1302/vega-asistant/uilib"
	"github.com/daniel1302/vega-asistant/utils"
)

type State int

const (
	StateGetVegaHome State = iota
	StateGetTendermintHome
	StateSummary
)

type ServiceSettings struct {
	VegaHome       string
	TendermintHome string
}

type StateMachine struct {
	Settings     ServiceSettings
	CurrentState State
}

func DefaultServiceSettings() ServiceSettings {
	return ServiceSettings{
		VegaHome:       filepath.Join(utils.CurrentUserHomePath(), "vega_home"),
		TendermintHome: filepath.Join(utils.CurrentUserHomePath(), "tendermint_home"),
	}
}

func NewStateMachine() StateMachine {
	return StateMachine{
		CurrentState: StateGetVegaHome,
		Settings:     DefaultServiceSettings(),
	}
}

func (state *StateMachine) Run(ui *input.UI) error {
STATE_RUN:
	for {
		switch state.CurrentState {
		case StateGetVegaHome:
			answer, err := uilib.AskString(ui, "What is your vega home?", state.Settings.VegaHome, checkIfExists("Vega Home"))
			if err != nil {
				return fmt.Errorf("failed to ask for vega home: %w", err)
			}
			state.Settings.VegaHome = answer
			state.CurrentState = StateGetTendermintHome

		case StateGetTendermintHome:
			answer, err := uilib.AskString(ui, "What is your tendermint home?", state.Settings.TendermintHome, checkIfExists("Tendermint Home"))
			if err != nil {
				return fmt.Errorf("failed to ask for tendermint home: %w", err)
			}
			state.Settings.TendermintHome = answer
			state.CurrentState = StateSummary

		case StateSummary:
			printSummary(state.Settings)
			answer, err := uilib.AskYesNo(ui, "Is it correct?", uilib.AnswerYes)
			if err != nil {
				return fmt.Errorf("failed to ask if summary correct: %w", err)
			}

			if answer == uilib.AnswerNo {
				state.CurrentState = StateGetVegaHome
			} else {
				break STATE_RUN
			}
		}
	}

	return nil
}

func checkIfExists(name string) func(string) error {
	return func(filePath string) error {
		if !utils.FileExists(filePath) {
			return fmt.Errorf(
				"%s(%s) does not exists: did you initialize your node?",
				name,
				filePath,
			)
		}
		return nil
	}
}
