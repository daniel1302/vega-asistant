package setup

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type PostgresqlDockerComposeArgs struct {
	*SetupArgs
}

var postgresqlDockerComposeArgs PostgresqlDockerComposeArgs

var postgresqlDockerComposeCmd = cobra.Command{
	Use:   "postgresql-docker-compose",
	Short: "Prepares docker-compose.yaml file to start the postgresql server with TimescaleDB extension enabled",
	RunE: func(cmd *cobra.Command, args []string) error {
		return setupPostgresqlDockerCompose(postgresqlDockerComposeArgs.Logger)
	},
}

func init() {
	postgresqlDockerComposeArgs.SetupArgs = &setupArgs
}

func setupPostgresqlDockerCompose(logger *zap.SugaredLogger) error {
	return nil
}
