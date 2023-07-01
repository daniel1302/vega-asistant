package setup

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
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

func setupPostStart(logger *zap.SugaredLogger) error {
	return nil
}
