package poststart

import (
	"fmt"
	"path/filepath"

	"go.uber.org/zap"

	"github.com/daniel1302/vega-assistant/utils"
	"github.com/daniel1302/vega-assistant/vegacmd"
)

func UpdateConfig(logger *zap.SugaredLogger, vegaHome, tendermintHome string) error {
	dataNodeConfig := map[string]interface{}{
		"SQLStore.WipeOnStartup":           false,
		"AutoInitialiseFromNetworkHistory": false,
	}
	tendermintConfig := map[string]interface{}{
		"statesync.enable": false,
	}
	coreConfig := map[string]interface{}{
		"Snapshot.StartHeight": -1,
	}

	dataNodeConfigPath := filepath.Join(vegaHome, vegacmd.DataNodeConfigPath)
	if !utils.FileExists(dataNodeConfigPath) {
		return fmt.Errorf("data node config(%s) does not exists", dataNodeConfigPath)
	}

	tendermintConfigPath := filepath.Join(tendermintHome, vegacmd.TenderminConfigPath)
	if !utils.FileExists(tendermintConfigPath) {
		return fmt.Errorf("tendermint config(%s) does not exists", tendermintConfigPath)
	}

	coreConfigPath := filepath.Join(vegaHome, vegacmd.CoreConfigPath)
	if !utils.FileExists(coreConfigPath) {
		return fmt.Errorf("vega core config(%s) does not exists", coreConfigPath)
	}

	logger.Infof(
		"Updating data-node config(%s). New values: %v",
		dataNodeConfigPath,
		dataNodeConfig,
	)
	if err := utils.UpdateConfig(dataNodeConfigPath, "toml", dataNodeConfig); err != nil {
		return fmt.Errorf("failed to update data node config(%s): %w", dataNodeConfigPath, err)
	}
	logger.Info("Data node config updated")

	logger.Infof("Updating core config(%s). New values: %v", coreConfigPath, coreConfig)
	if err := utils.UpdateConfig(coreConfigPath, "toml", coreConfig); err != nil {
		return fmt.Errorf("failed to update core config(%s): %w", coreConfigPath)
	}
	logger.Info("Core config updated")

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
