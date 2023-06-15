package main

import (
	"fmt"
	"os"

	"github.com/tcnksm/go-input"
	"go.uber.org/zap"

	"github.com/daniel1302/vega-asistant/generator"
	"github.com/daniel1302/vega-asistant/network"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	ui := &input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}

	state := generator.NewStateMachine()
	state.Run(ui)

	generator, _ := generator.NewDataNodeGenerator(state.Settings, network.MainnetConfig())
	if err := generator.Run(sugar); err != nil {
		panic(err)
	}

	fmt.Printf("%s", state.Dump())
}
