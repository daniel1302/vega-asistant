package generator

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/tcnksm/go-input"

	"github.com/daniel1302/vega-asistant/utils"
)

type (
	State       int
	StartupMode string
)

const (
	StartFromBlock0         StartupMode = "start-from-block-0"
	StartFromNetworkHistory StartupMode = "startup-from-network-history"
)

const (
	StateSelectStartupMode State = iota
	StateSelectVisorHome
	StateSelectVegaHome
	StateSelectTendermintHome
)

type StateMachine struct {
	CurrentState State

	Settings GenerateSettings
}

type GenerateSettings struct {
	Mode StartupMode

	VisorHome      string
	VegaHome       string
	TendermintHome string
	DataNodeHome   string
}

func NewStateMachine() StateMachine {
	return StateMachine{
		CurrentState: StateSelectStartupMode,
	}
}

func (state StateMachine) Dump() string {
	result, err := json.MarshalIndent(state, "", "    ")
	if err != nil {
		return ""
	}

	return string(result)
}

func (state *StateMachine) Run(ui *input.UI) error {
STATE_RUN:
	for {
		switch state.CurrentState {
		case StateSelectStartupMode:
			mode, err := SelectStartupMode(ui)
			if err != nil {
				return fmt.Errorf("failed selecting startup mode: %w", err)
			}
			state.Settings.Mode = *mode
			state.CurrentState = StateSelectVisorHome

		case StateSelectVisorHome:
			defaultValue := filepath.Join(utils.CurrentUserHomePath(), "vegavisor_home")
			visorHome, err := AskPath(ui, "vegavisor home", defaultValue)
			if err != nil {
				return fmt.Errorf("failed getting vegavisor home: %w", err)
			}

			state.Settings.VisorHome = visorHome
			state.CurrentState = StateSelectVegaHome

		case StateSelectVegaHome:
			defaultValue := filepath.Join(utils.CurrentUserHomePath(), "vega_home")
			vegaHome, err := AskPath(ui, "vega home", defaultValue)
			if err != nil {
				return fmt.Errorf("failed getting vega home: %w", err)
			}
			state.Settings.VegaHome = vegaHome
			state.Settings.DataNodeHome = vegaHome
			state.CurrentState = StateSelectTendermintHome

		case StateSelectTendermintHome:
			defaultValue := filepath.Join(utils.CurrentUserHomePath(), "tendermint_home")
			tendermintHome, err := AskPath(ui, "tendermint home", defaultValue)
			if err != nil {
				return fmt.Errorf("failed getting tendermint home: %w", err)
			}
			state.Settings.TendermintHome = tendermintHome
			break STATE_RUN

		}
	}
	return nil
}
