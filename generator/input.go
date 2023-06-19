package generator

import (
	"fmt"
	"strconv"

	input "github.com/tcnksm/go-input"

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

func AskSQLCredentials(ui *input.UI) (*types.SQLCredentials, error) {
	var (
		dbHost string
		dbUser string
		dbPort int
		dbPass string
		dbName string

		err error
	)
	for {
		dbHost, err = ui.Ask("PostgreSQL host for the data-node", &input.Options{
			Default:  "localhost",
			Required: true,
			Loop:     true,
		})
		if err != nil {
			return nil, types.NewInputError(fmt.Errorf("failed to get postgresql host: %w", err))
		}

		dbPortStr, err := ui.Ask("PostgreSQL port for the data-node", &input.Options{
			Default:  "5432",
			Required: true,
			Loop:     true,
			ValidateFunc: func(s string) error {
				if _, err := strconv.Atoi(s); err != nil {
					return fmt.Errorf("port must be numeric: %w", err)
				}

				return nil
			},
		})
		if err != nil {
			return nil, types.NewInputError(fmt.Errorf("failed to get postgresql port: %w", err))
		}

		dbPort, err = strconv.Atoi(dbPortStr)
		if err != nil {
			return nil, types.NewInputError(fmt.Errorf("port must be numeric: %w", err))
		}

		dbUser, err = ui.Ask("PostgreSQL user name for the data-node", &input.Options{
			Default:  "vega",
			Required: true,
			Loop:     true,
		})

		dbPass, err = ui.Ask("PostgreSQL password for the given username", &input.Options{
			Default:  "vega",
			Required: true,
			Loop:     true,
		})

		dbName, err = ui.Ask("PostgreSQL database name for the data-node", &input.Options{
			Default:  "vega",
			Required: true,
			Loop:     true,
		})
		break
	}

	return &types.SQLCredentials{
		Host:         dbHost,
		User:         dbUser,
		Port:         dbPort,
		Pass:         dbPass,
		DatabaseName: dbName,
	}, nil
}
