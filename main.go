package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/tcnksm/go-input"
	"go.uber.org/zap"

	"github.com/daniel1302/vega-asistant/generator"
	"github.com/daniel1302/vega-asistant/network"
)

func main() {
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
	defer logger.Sync()

	sugar := logger.Sugar()
	ui := &input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}

	state := generator.NewStateMachine()
	err := state.Run(ui)
	if err != nil {
		panic(fmt.Errorf("failed to generate data-node: %w", err))
	}

	generator, _ := generator.NewDataNodeGenerator(state.Settings, network.MainnetConfig())
	if err := generator.Run(sugar); err != nil {
		panic(err)
	}

	fmt.Printf("%s", state.Dump())
}
