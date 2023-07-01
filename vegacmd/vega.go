package vegacmd

import (
	"fmt"

	"github.com/daniel1302/vega-assistant/utils"
)

func InitVega(binaryPath, vegaHome string, nodeMode VegaNodeMode) error {
	_, err := utils.ExecuteBinary(
		binaryPath,
		[]string{"init", "--output", "json", "--home", vegaHome, string(nodeMode)},
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to init vega: %w", err)
	}

	return nil
}
