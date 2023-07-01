package setup

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	generator "github.com/daniel1302/vega-asistant/generator/systemd"
	"github.com/daniel1302/vega-asistant/utils"
)

type SystemdArgs struct {
	*SetupArgs
	VisorHome string
}

var systemdArgs SystemdArgs

var systemdCmd = &cobra.Command{
	Use:   "systemd",
	Short: "Prepares systemd configuration for the data-node",
	Run: func(cmd *cobra.Command, args []string) {
		if err := setupSystemd(systemdArgs.Logger, systemdArgs.VisorHome); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	systemdArgs.SetupArgs = &setupArgs

	systemdCmd.PersistentFlags().
		StringVar(&systemdArgs.VisorHome, "visor-home", filepath.Join(utils.CurrentUserHomePath(), "vegavisor_home"), "The vegavisor home path")
}

func setupSystemd(logger *zap.SugaredLogger, visorHome string) error {
	if err := generator.PrepareSystemd(logger, visorHome); err != nil {
		return fmt.Errorf("failed to prepare systemd service: %w", err)
	}

	generator.PrintInstructions()
	return nil
}
