package main

import (
	"fmt"
	"os"

	"github.com/daniel1302/vega-asistant/cmd"
	"github.com/daniel1302/vega-asistant/cmd/setup"
)

func init() {
	cmd.RootCmd.AddCommand(setup.RootCmd)
}

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
