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

	ConfigFile string
}

var setupDataNodeArgs SetupDataNodeArgs

var dataNodeCmd = &cobra.Command{
	Use:   "data-node",
	Short: "Prepare data-node on your computer",
	RunE: func(cmd *cobra.Command, args []string) error {
		return dataNodeSetup(setupDataNodeArgs.Logger, setupDataNodeArgs.ConfigFile)
	},
}

func init() {
	setupDataNodeArgs.SetupArgs = &setupArgs
	dataNodeCmd.PersistentFlags().StringVar(
		&setupDataNodeArgs.ConfigFile,
		"config-file",
		"config.toml",
		"Config file to read values from. If there is an error in config file, default values are used",
	)
}

func dataNodeSetup(logger *zap.SugaredLogger, configFile string) error {
	ui := &input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}
	config, err := service.ReadGeneratorSettingsFromFile(configFile)
	if err != nil {
		logger.Info("Could not load config file. Using default values", zap.String("reason", err.Error()))

		config = service.DefaultGenerateSettings()
	}

	state := service.NewStateMachine(logger, *config)
	if err := state.Run(ui, network.MainnetConfig()); err != nil {
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
