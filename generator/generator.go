package generator

import (
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/daniel1302/vega-asistant/network"
)

type DataNodeGenerator struct {
	userSettings  GenerateSettings
	networkConfig network.NetworkConfig
}

func NewDataNodeGenerator(
	settings GenerateSettings,
	networkConfig network.NetworkConfig,
) (*DataNodeGenerator, error) {
	return &DataNodeGenerator{
		userSettings:  settings,
		networkConfig: networkConfig,
	}, nil
}

func (gen *DataNodeGenerator) Run(logger *zap.SugaredLogger) error {
	outputDir, err := os.MkdirTemp("", "vega-assistant")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(outputDir)

	logger.Info("Downloading vega binary")
	vegaBinaryPath, err := DownloadArtifact(
		gen.networkConfig.Repository,
		gen.userSettings.MainnetVersion,
		outputDir,
		ArtifactVega,
	)
	if err != nil {
		return fmt.Errorf("failed to download vega binary: %w", err)
	}
	logger.Info("Vega downloaded to %s", vegaBinaryPath)

	logger.Info("Downloading visor binary")
	visorBinaryPath, err := DownloadArtifact(
		gen.networkConfig.Repository,
		gen.userSettings.MainnetVersion,
		outputDir,
		ArtifactVisor,
	)
	if err != nil {
		return fmt.Errorf("failed to download visor binary: %w", err)
	}
	logger.Info("Visor downloaded to %s", visorBinaryPath)

	return nil
}
