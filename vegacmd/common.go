package vegacmd

import "path/filepath"

type VegaNodeMode string

const (
	VegaNodeFull      VegaNodeMode = "full"
	VegaNodeValidator VegaNodeMode = "validator"
	VegaNodeSeed      VegaNodeMode = "seed"
)

var (
	CoreConfigPath      = filepath.Join("config", "node", "config.toml")
	DataNodeConfigPath  = filepath.Join("config", "data-node", "config.toml")
	VegavisorConfigPath = filepath.Join("config.toml")
	TenderminConfigPath = filepath.Join("config", "config.toml")
	GenesisPath         = filepath.Join("config", "genesis.json")
)

func BinaryVersion(binary string) (string, error) {
	return "", nil
}
