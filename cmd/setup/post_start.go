package setup

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tcnksm/go-input"
	"go.uber.org/zap"

	service "github.com/daniel1302/vega-asistant/service/poststart"
)

type PostStartArgs struct {
	*SetupArgs
}

var postStartArgs PostStartArgs

var postStartCmd = &cobra.Command{
	Use:   "post-start",
	Short: "Put configuration adjustments required after node has been started",
	Run: func(cmd *cobra.Command, args []string) {
		if err := setupPostStart(postStartArgs.Logger); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	postStartArgs.SetupArgs = &setupArgs
}

func setupPostStart(logger *zap.SugaredLogger) error {
	ui := &input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}
	state := service.NewStateMachine()
	err := state.Run(ui)
	if err != nil {
		return fmt.Errorf("failed to run state machine: %w", err)
	}

	if err := service.UpdateConfig(logger, state.Settings.VegaHome, state.Settings.TendermintHome); err != nil {
		return fmt.Errorf("failed to update configs: %s", err)
	}

	return nil
}
