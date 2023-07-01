package datanode

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	pg "github.com/go-pg/pg/v11"
	"github.com/tcnksm/go-input"
	"golang.org/x/mod/semver"

	"github.com/daniel1302/vega-assistant/network"
	"github.com/daniel1302/vega-assistant/types"
	"github.com/daniel1302/vega-assistant/uilib"
	"github.com/daniel1302/vega-assistant/utils"
	"github.com/daniel1302/vega-assistant/vegaapi"
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
	StateExistingVisorHome
	StateSelectVegaHome
	StateExistingVegaHome
	StateSelectTendermintHome
	StateExistingTendermintHome
	StateGetSQLCredentials
	StateCheckLatestVersion
	StateSummary
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
	MainnetVersion string
	MainnetChainId string
	SQLCredentials types.SQLCredentials
}

func DefaultGenerateSettings() GenerateSettings {
	return GenerateSettings{
		Mode:           StartFromNetworkHistory,
		VisorHome:      filepath.Join(utils.CurrentUserHomePath(), "vegavisor_home"),
		VegaHome:       filepath.Join(utils.CurrentUserHomePath(), "vega_home"),
		TendermintHome: filepath.Join(utils.CurrentUserHomePath(), "tendermint_home"),

		SQLCredentials: types.SQLCredentials{
			Host:         "localhost",
			User:         "vega",
			Pass:         "vega",
			Port:         5432,
			DatabaseName: "vega",
		},
	}
}

func NewStateMachine() StateMachine {
	return StateMachine{
		CurrentState: StateSelectStartupMode,
		Settings:     DefaultGenerateSettings(),
	}
}

func (state StateMachine) Dump() string {
	result, err := json.MarshalIndent(state, "", "    ")
	if err != nil {
		return ""
	}

	return string(result)
}

func (state *StateMachine) Run(ui *input.UI, networkConfig network.NetworkConfig) error {
STATE_RUN:
	for {
		switch state.CurrentState {
		case StateSelectStartupMode:
			mode, err := SelectStartupMode(ui, state.Settings.Mode)
			if err != nil {
				return fmt.Errorf("failed selecting startup mode: %w", err)
			}
			state.Settings.Mode = *mode
			state.CurrentState = StateSelectVisorHome

		case StateSelectVisorHome:
			visorHome, err := uilib.AskPath(ui, "vegavisor home", state.Settings.VisorHome)
			if err != nil {
				return fmt.Errorf("failed getting vegavisor home: %w", err)
			}

			state.Settings.VisorHome = visorHome
			if utils.FileExists(visorHome) {
				state.CurrentState = StateExistingVisorHome
			} else {
				state.CurrentState = StateSelectVegaHome
			}

		case StateExistingVisorHome:
			removeAnswer, err := uilib.AskRemoveExistingFile(ui, state.Settings.VisorHome, uilib.AnswerYes)
			if err != nil {
				return fmt.Errorf("failed to get answer for remove existing visor home: %w", err)
			}

			if removeAnswer == uilib.AnswerNo {
				return fmt.Errorf("visor home exists. You must provide different visor home or remove it")
			}

			if err := os.RemoveAll(state.Settings.VisorHome); err != nil {
				return fmt.Errorf("failed to remove vegavisor home: %w", err)
			}

			state.CurrentState = StateSelectVegaHome

		case StateSelectVegaHome:
			vegaHome, err := uilib.AskPath(ui, "vega home", state.Settings.VegaHome)
			if err != nil {
				return fmt.Errorf("failed getting vega home: %w", err)
			}
			state.Settings.VegaHome = vegaHome
			state.Settings.DataNodeHome = vegaHome

			if utils.FileExists(vegaHome) {
				state.CurrentState = StateExistingVegaHome
			} else {
				state.CurrentState = StateSelectTendermintHome
			}

		case StateExistingVegaHome:
			removeAnswer, err := uilib.AskRemoveExistingFile(ui, state.Settings.VegaHome, uilib.AnswerYes)
			if err != nil {
				return fmt.Errorf("failed to get answer for remove existing vega home: %w", err)
			}

			if removeAnswer == uilib.AnswerNo {
				return fmt.Errorf("vega home exists. You must provide different vega home or remove it")
			}

			if err := os.RemoveAll(state.Settings.VegaHome); err != nil {
				return fmt.Errorf("failed to remove vega home: %w", err)
			}

			state.CurrentState = StateSelectTendermintHome

		case StateSelectTendermintHome:
			tendermintHome, err := uilib.AskPath(ui, "tendermint home", state.Settings.TendermintHome)
			if err != nil {
				return fmt.Errorf("failed getting tendermint home: %w", err)
			}
			state.Settings.TendermintHome = tendermintHome

			if utils.FileExists(tendermintHome) {
				state.CurrentState = StateExistingTendermintHome
			} else {
				state.CurrentState = StateGetSQLCredentials
			}

		case StateExistingTendermintHome:
			removeAnswer, err := uilib.AskRemoveExistingFile(ui, state.Settings.TendermintHome, uilib.AnswerYes)
			if err != nil {
				return fmt.Errorf("failed to get answer for remove existing tendermint home: %w", err)
			}

			if removeAnswer == uilib.AnswerNo {
				return fmt.Errorf("tendermint home exists. You must provide different tendermint home or remove it")
			}

			if err := os.RemoveAll(state.Settings.TendermintHome); err != nil {
				return fmt.Errorf("failed to remove tendermint home: %w", err)
			}

			state.CurrentState = StateGetSQLCredentials

		case StateGetSQLCredentials:
			sqlCredentials, err := AskSQLCredentials(ui, state.Settings.SQLCredentials, checkSQLCredentials)
			if err != nil {
				return fmt.Errorf("failed getting sql credentials: %w", err)
			}
			state.Settings.SQLCredentials = *sqlCredentials
			state.CurrentState = StateCheckLatestVersion

		case StateCheckLatestVersion:
			statisticsResponse, err := vegaapi.Statistics(networkConfig.DataNodesRESTUrls)
			if err != nil {
				return fmt.Errorf("failed to get response for the /statistics endpoint from the mainnet servers: %w", err)
			}
			if state.Settings.Mode == StartFromBlock0 {
				state.Settings.MainnetVersion = networkConfig.GenesisVersion
			} else {
				state.Settings.MainnetVersion = statisticsResponse.Statistics.AppVersion
			}

			state.Settings.MainnetChainId = statisticsResponse.Statistics.ChainID
			state.CurrentState = StateSummary

		case StateSummary:
			printSummary(state.Settings)

			correctResponse, err := uilib.AskYesNo(ui, "Is it correct?", uilib.AnswerYes)
			if err != nil {
				return fmt.Errorf("failed asking for correct summary: %w", err)
			}

			if correctResponse == uilib.AnswerNo {
				state.CurrentState = StateSelectStartupMode
				break
			}

			break STATE_RUN
		}
	}
	return nil
}

func checkSQLCredentials(creds types.SQLCredentials) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%d", creds.Host, creds.Port),
		User:     creds.User,
		Password: creds.Pass,
		Database: creds.DatabaseName,
	})
	defer db.Close(ctx)

	var n int
	_, err := db.QueryOne(ctx, pg.Scan(&n), "SELECT 1")
	if err != nil {
		return err
	}

	var timescaleVersion string
	_, err = db.QueryOne(
		ctx,
		pg.Scan(&timescaleVersion),
		"SELECT extversion FROM pg_extension WHERE extname = 'timescaledb'",
	)
	if err != nil {
		return fmt.Errorf("failed to check timescale extension version: %w", err)
	}

	if !strings.HasPrefix(timescaleVersion, "v") {
		timescaleVersion = fmt.Sprintf("v%s", timescaleVersion)
	}

	if semver.Compare(timescaleVersion, "v2.8.0") != 0 {
		return fmt.Errorf(
			"Vega support only timescale v2.8.0. Installed version is %s",
			timescaleVersion,
		)
	}

	return nil
}
