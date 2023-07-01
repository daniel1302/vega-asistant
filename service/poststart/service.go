package poststart

import (
	"fmt"
	"path/filepath"

	"go.uber.org/zap"

	"github.com/daniel1302/vega-asistant/utils"
	"github.com/daniel1302/vega-asistant/vegacmd"
)

func UpdateConfig(logger *zap.SugaredLogger, vegaHome, tendermintHome string) error {
	dataNodeConfig := map[string]interface{}{
		"SQLStore.WipeOnStartup":           false,
		"AutoInitialiseFromNetworkHistory": false,
	}
	tendermintConfig := map[string]interface{}{
		"statesync.enable": false,
	}

	dataNodeConfigPath := filepath.Join(vegaHome, vegacmd.DataNodeConfigPath)
	if !utils.FileExists(dataNodeConfigPath) {
		return fmt.Errorf("data node config(%s) does not exists", dataNodeConfigPath)
	}

	tendermintConfigPath := filepath.Join(tendermintHome, vegacmd.TenderminConfigPath)
	if !utils.FileExists(tendermintConfigPath) {
		return fmt.Errorf("tendermint config(%s) does not exists", tendermintConfigPath)
	}

	logger.Infof("Updating core config(%s). New values: %v", dataNodeConfigPath, dataNodeConfig)
	if err := utils.UpdateConfig(dataNodeConfigPath, "toml", dataNodeConfig); err != nil {
		return fmt.Errorf("failed to update data node config(%s): %w", dataNodeConfigPath, err)
	}
	logger.Info("Data node config updated")

	logger.Infof(
		"Updating tendermint config(%s). New values: %v",
		tendermintConfigPath,
		dataNodeConfig,
	)
	if err := utils.UpdateConfig(tendermintConfigPath, "toml", tendermintConfig); err != nil {
		return fmt.Errorf("failed to update tendermitn config(%s): %w", tendermintConfigPath, err)
	}
	logger.Info("Tendermint config updated")

	return nil
}
