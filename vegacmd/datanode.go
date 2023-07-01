package vegacmd

import (
	"fmt"

	"github.com/daniel1302/vega-asistant/utils"
)

func InitDataNode(binaryPath, vegaHome string, chainId string) error {
	_, err := utils.ExecuteBinary(
		binaryPath,
		[]string{"datanode", "init", "--home", vegaHome, chainId},
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize data-node: %w", err)
	}

	return nil
}
