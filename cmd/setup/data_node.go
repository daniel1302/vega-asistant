package setup

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tcnksm/go-input"
	"go.uber.org/zap"

	"github.com/daniel1302/vega-assistant/network"
	service "github.com/daniel1302/vega-assistant/service/datanode"
)

type SetupDataNodeArgs struct {
	*SetupArgs
}

var setupDataNodeArgs SetupDataNodeArgs

var dataNodeCmd = &cobra.Command{
	Use:   "data-node",
	Short: "Prepare data-node on your computer",
	RunE: func(cmd *cobra.Command, args []string) error {
		return dataNodeSetup(setupDataNodeArgs.Logger)
	},
}

func init() {
	setupDataNodeArgs.SetupArgs = &setupArgs
}

func dataNodeSetup(logger *zap.SugaredLogger) error {
	ui := &input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}
	state := service.NewStateMachine()
	err := state.Run(ui, network.MainnetConfig())
	if err != nil {
		return fmt.Errorf("failed to generate data-node: %w", err)
	}

	svc, err := service.NewDataNodeGenerator(state.Settings, network.MainnetConfig())
	if err != nil {
		return fmt.Errorf("failed to start generator service: %w", err)
	}
	if err := svc.Run(logger); err != nil {
		return fmt.Errorf("failed to setup data-node: %w", err)
	}

	service.PrintInstructions(state.Settings.VisorHome)

	return nil
}
