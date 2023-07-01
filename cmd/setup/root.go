package setup

import (
	"github.com/spf13/cobra"

	"github.com/daniel1302/vega-assistant/cmd"
)

type SetupArgs struct {
	*cmd.RootArgs
}

var setupArgs SetupArgs

// Root Command for OPS
var RootCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup a node commands",
}

func init() {
	setupArgs.RootArgs = &cmd.Args

	RootCmd.AddCommand(dataNodeCmd)
	RootCmd.AddCommand(postgresqlDockerComposeCmd)
	RootCmd.AddCommand(systemdCmd)
	RootCmd.AddCommand(postStartCmd)
}
