package generator

import (
	"fmt"

	"github.com/tcnksm/go-input"

	"github.com/daniel1302/vega-asistant/types"
	"github.com/daniel1302/vega-asistant/utils"
)

func SelectStartupMode(ui *input.UI) (*StartupMode, error) {
	const msg = `How do you want to start your data-node?

  - Starting from block 0 - Starts the node from the genesis binary, replays all the blocks and 
                            does all of the protocol upgrades automatically.
        * Depending on network age it can takes up to several days to catch your node up.
        * Full network history is available on your node.
  
  - Starting from network history - Start the node from latest binary, and download all required 
                                    informations from the running network. 
        * It takes up to several minutes.
        * No historical data is available on your node.`
	response, err := ui.Select(
		msg,
		[]string{string(StartFromBlock0), string(StartFromNetworkHistory)},
		&input.Options{
			Default:  string(StartFromNetworkHistory),
			Loop:     true,
			Required: true,
		},
	)
	if err != nil {
		return nil, types.NewInputError(err)
	}

	result := StartFromNetworkHistory
	if response == string(StartFromBlock0) {
		result = StartFromBlock0
	}

	return &result, nil
}

func AskPath(ui *input.UI, name, defaultValue string) (string, error) {
	response, err := ui.Ask(fmt.Sprintf("What is your %s", name), &input.Options{
		Default:  defaultValue,
		Required: true,
		Loop:     true,
		ValidateFunc: func(s string) error {
			if utils.FileExists(s) {
				return fmt.Errorf(
					"given path exists on your fs, remove this file or provide another directory",
				)
			}

			return nil
		},
	})
	if err != nil {
		return "", types.NewInputError(err)
	}

	return response, nil
}
