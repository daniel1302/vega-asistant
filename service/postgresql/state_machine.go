package postgresql

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/tcnksm/go-input"

	"github.com/daniel1302/vega-asistant/uilib"
	"github.com/daniel1302/vega-asistant/utils"
)

type State int

const (
	StateGetHome State = iota
	StateExistingHome
	StateGetPostgresqlUsername
	StateGetPostgresqlPassword
	StateGetPostgresqlDatabase
	StateGetPostgresqlPort
	StateSummary
)

type GeneratorSettings struct {
	Home               string
	PostgresqlUsername string
	PostgresqlPassword string
	PostgresqlDatabase string
	PostgresqlPort     int
}

type StateMachine struct {
	Settings     GeneratorSettings
	CurrentState State
}

func DefaultGeneratorSettings() GeneratorSettings {
	return GeneratorSettings{
		Home:               filepath.Join(utils.CurrentUserHomePath(), "vega_postgresql"),
		PostgresqlUsername: "vega",
		PostgresqlPassword: "vega",
		PostgresqlDatabase: "vega",
		PostgresqlPort:     5432,
	}
}

func NewStateMachine() StateMachine {
	return StateMachine{
		CurrentState: StateGetHome,
		Settings:     DefaultGeneratorSettings(),
	}
}

func (state *StateMachine) Run(ui *input.UI) error {
STATE_RUN:
	for {
		switch state.CurrentState {
		case StateGetHome:
			answer, err := uilib.AskPath(ui, "Home for docker-compose", state.Settings.Home)
			if err != nil {
				return fmt.Errorf("failed to ask for home: %w", err)
			}

			state.Settings.Home = answer
			if utils.FileExists(answer) {
				state.CurrentState = StateExistingHome
			} else {
				state.CurrentState = StateGetPostgresqlUsername
			}

		case StateExistingHome:
			removeAnswer, err := uilib.AskRemoveExistingFile(ui, state.Settings.Home, uilib.AnswerYes)
			if err != nil {
				return fmt.Errorf("failed to get answer for remove existing home: %w", err)
			}

			if removeAnswer == uilib.AnswerNo {
				return fmt.Errorf("the home dir exists. You must provide different home or remove it")
			}

			if err := os.RemoveAll(state.Settings.Home); err != nil {
				return fmt.Errorf("failed to remove home: %w", err)
			}

			state.CurrentState = StateGetPostgresqlUsername

		case StateGetPostgresqlUsername:
			username, err := uilib.AskString(ui, "PostgreSQL username", state.Settings.PostgresqlUsername, validatePostgreSQLCredentialsString)
			if err != nil {
				return fmt.Errorf("failed to ask for PostgreSQL username: %w", err)
			}
			state.Settings.PostgresqlUsername = username
			state.CurrentState = StateGetPostgresqlPassword

		case StateGetPostgresqlPassword:
			password, err := uilib.AskString(ui, "PostgreSQL user password", state.Settings.PostgresqlPassword, validatePostgreSQLCredentialsString)
			if err != nil {
				return fmt.Errorf("failed to ask for PostgreSQL password: %w", err)
			}
			state.Settings.PostgresqlPassword = password
			state.CurrentState = StateGetPostgresqlDatabase

		case StateGetPostgresqlDatabase:
			databaseName, err := uilib.AskString(ui, "PostgreSQL database name", state.Settings.PostgresqlDatabase, validatePostgreSQLCredentialsString)
			if err != nil {
				return fmt.Errorf("failed to ask for PostgreSQL database name: %w", err)
			}
			state.Settings.PostgresqlDatabase = databaseName
			state.CurrentState = StateGetPostgresqlPort

		case StateGetPostgresqlPort:
			port, err := uilib.AskInt(ui, "PostgreSQL port", state.Settings.PostgresqlPort)
			if err != nil {
				return fmt.Errorf("failed to ask for PostgreSQL port: %w", err)
			}
			state.Settings.PostgresqlPort = port
			state.CurrentState = StateSummary

		case StateSummary:
			printSummary(state.Settings)

			answer, err := uilib.AskYesNo(ui, "Is it correct?", uilib.AnswerYes)
			if err != nil {
				return fmt.Errorf("failed to ask if summary correct: %w", err)
			}

			if answer == uilib.AnswerNo {
				state.CurrentState = StateGetHome
			} else {
				break STATE_RUN
			}
		}
	}

	return nil
}

func validatePostgreSQLCredentialsString(s string) error {
	strRegex := regexp.MustCompile(`^[A-Za-z0-9_\.-]{3,}$`)

	if !strRegex.MatchString(s) {
		return fmt.Errorf(
			"string '%s' must contains ony digits, characters and the following chars: _.-",
		)
	}
	return nil
}
