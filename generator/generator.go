package generator

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/daniel1302/vega-asistant/network"
	"github.com/daniel1302/vega-asistant/vegaapi"
)

type DataNodeGenerator struct {
	VegaVersion string
	ChainID     string

	VisorBinaryPath string
	VegaBinaryPath  string

	UserSettings  GenerateSettings
	NetworkConfig network.NetworkConfig
}

func NewDataNodeGenerator(
	settings GenerateSettings,
	networkConfig network.NetworkConfig,
) (*DataNodeGenerator, error) {
	return &DataNodeGenerator{
		UserSettings:  settings,
		NetworkConfig: networkConfig,
	}, nil
}

func (gen *DataNodeGenerator) Run(logger *zap.SugaredLogger) error {
	// TODO: add validation for network config and user settings

	logger.Info("DDD")
	statistics, err := vegaapi.Statistics(gen.NetworkConfig.DataNodesRESTUrls)
	if err != nil {
		return fmt.Errorf("failed to get statistics from the network: %w", err)
	}

	fmt.Println("version is: %s", statistics.Statistics.AppVersion)
	return nil
}
