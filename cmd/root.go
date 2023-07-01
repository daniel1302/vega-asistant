package cmd

import (
	"encoding/json"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type RootArgs struct {
	Logger *zap.SugaredLogger
}

var Args RootArgs

var RootCmd = &cobra.Command{
	Use:   "vega-assistant",
	Short: "Helps manage vega manual way",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		rawJSON := []byte(`{
		"level": "info",
		"outputPaths": ["stdout"],
		"errorOutputPaths": ["stderr"],
		"encoding": "console",
		"encoderConfig": {
			"messageKey": "message",
			"levelEncoder": "lowercase"
		}
	}`)
		var cfg zap.Config
		if err := json.Unmarshal(rawJSON, &cfg); err != nil {
			panic(err)
		}
		logger := zap.Must(cfg.Build())

		Args.Logger = logger.Sugar()
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if Args.Logger != nil {
			defer Args.Logger.Sync()
		}
	},
}
