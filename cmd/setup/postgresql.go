package setup

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tcnksm/go-input"
	"go.uber.org/zap"

	generator "github.com/daniel1302/vega-asistant/generator/postgresql"
)

type PostgresqlDockerComposeArgs struct {
	*SetupArgs
}

var postgresqlDockerComposeArgs PostgresqlDockerComposeArgs

var postgresqlDockerComposeCmd = &cobra.Command{
	Use:   "postgresql",
	Short: "Prepares docker-compose.yaml file to start the postgresql server with TimescaleDB extension enabled",
	RunE: func(cmd *cobra.Command, args []string) error {
		return setupPostgresqlDockerCompose(postgresqlDockerComposeArgs.Logger)
	},
}

func init() {
	postgresqlDockerComposeArgs.SetupArgs = &setupArgs
}

func setupPostgresqlDockerCompose(logger *zap.SugaredLogger) error {
	ui := &input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}
	state := generator.NewStateMachine()
	err := state.Run(ui)
	if err != nil {
		return fmt.Errorf("failed to run state machine: %w", err)
	}

	if err := generator.PrepareDockerComposeFile(logger, state.Settings); err != nil {
		return fmt.Errorf("failed to prepare docker-compose.yaml: %w", err)
	}

	generator.PrintInstructions(state.Settings.Home)

	return nil
}
