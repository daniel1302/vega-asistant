package vegacmd

import (
	"fmt"

	"github.com/daniel1302/vega-asistant/utils"
)

func InitTendermint(binaryPath, tendermintHome string) error {
	_, err := utils.ExecuteBinary(binaryPath, []string{"tm", "init", "--home", tendermintHome}, nil)
	if err != nil {
		return fmt.Errorf("failed to init tendermint: %w", err)
	}

	return nil
}
