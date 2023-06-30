package generator

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	pg "github.com/go-pg/pg/v11"
	"github.com/tcnksm/go-input"
	"golang.org/x/mod/semver"

	"github.com/daniel1302/vega-asistant/network"
	"github.com/daniel1302/vega-asistant/types"
	"github.com/daniel1302/vega-asistant/utils"
	"github.com/daniel1302/vega-asistant/vegaapi"
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

func (state *StateMachine) Run(ui *input.UI) error {
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
			visorHome, err := AskPath(ui, "vegavisor home", state.Settings.VisorHome)
			if err != nil {
				return fmt.Errorf("failed getting vegavisor home: %w", err)
			}

			state.Settings.VisorHome = visorHome
			state.CurrentState = StateSelectVegaHome

		case StateSelectVegaHome:
			vegaHome, err := AskPath(ui, "vega home", state.Settings.VegaHome)
			if err != nil {
				return fmt.Errorf("failed getting vega home: %w", err)
			}
			state.Settings.VegaHome = vegaHome
			state.Settings.DataNodeHome = vegaHome
			state.CurrentState = StateSelectTendermintHome

		case StateSelectTendermintHome:
			tendermintHome, err := AskPath(ui, "tendermint home", state.Settings.TendermintHome)
			if err != nil {
				return fmt.Errorf("failed getting tendermint home: %w", err)
			}
			state.Settings.TendermintHome = tendermintHome
			state.CurrentState = StateGetSQLCredentials

		case StateGetSQLCredentials:
			sqlCredentials, err := AskSQLCredentials(ui, state.Settings.SQLCredentials, checkSQLCredentials)
			if err != nil {
				return fmt.Errorf("failed getting sql credentials: %w", err)
			}
			state.Settings.SQLCredentials = *sqlCredentials
			state.CurrentState = StateCheckLatestVersion

		case StateCheckLatestVersion:
			statisticsResponse, err := vegaapi.Statistics(network.MainnetConfig().DataNodesRESTUrls)
			if err != nil {
				return fmt.Errorf("failed to get response for the /statistics endpoint from the mainnet servers: %w", err)
			}
			if state.Settings.Mode == StartFromBlock0 {
				state.Settings.MainnetVersion = network.MainnetConfig().GenesisVersion
			} else {
				state.Settings.MainnetVersion = statisticsResponse.Statistics.AppVersion
			}

			state.Settings.MainnetChainId = statisticsResponse.Statistics.ChainID
			state.CurrentState = StateSummary

		case StateSummary:
			printSummary(state.Settings)

			correctResponse, err := AskYesNo(ui, "Is it correct?", AnswerYes)
			if err != nil {
				return fmt.Errorf("failed asking for correct summary: %w", err)
			}

			if correctResponse == AnswerNo {
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
