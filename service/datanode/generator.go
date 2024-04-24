package datanode

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"github.com/daniel1302/vega-assistant/github"
	"github.com/daniel1302/vega-assistant/network"
	"github.com/daniel1302/vega-assistant/types"
	"github.com/daniel1302/vega-assistant/utils"
	"github.com/daniel1302/vega-assistant/vegaapi"
	"github.com/daniel1302/vega-assistant/vegacmd"
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
	//	defer os.RemoveAll(outputDir)

	logger.Info("Downloading vega binary")
	vegaBinaryPath, err := github.DownloadArtifact(
		gen.networkConfig.Repository,
		gen.userSettings.VegaBinaryVersion,
		outputDir,
		github.ArtifactVega,
	)
	if err != nil {
		return fmt.Errorf("failed to download vega binary: %w", err)
	}
	logger.Infof("Vega downloaded to %s", vegaBinaryPath)

	logger.Info("Downloading visor binary")
	visorBinaryPath, err := github.DownloadArtifact(
		gen.networkConfig.Repository,
		gen.userSettings.VisorBinaryVersion,
		outputDir,
		github.ArtifactVisor,
	)
	if err != nil {
		return fmt.Errorf("failed to download visor binary: %w", err)
	}
	logger.Infof("Visor downloaded to %s", visorBinaryPath)

	logger.Info("Checking binaries versions")
	vegaVersion, err := utils.ExecuteBinary(vegaBinaryPath, []string{"version"}, nil)
	if err != nil {
		return fmt.Errorf("failed to check vega version: %w", err)
	}
	logger.Infof("Vega version is %s", vegaVersion)
	VisorBinaryVersion, err := utils.ExecuteBinary(visorBinaryPath, []string{"version"}, nil)
	if err != nil {
		return fmt.Errorf("failed to check visor version: %w", err)
	}
	logger.Infof("Visor version is %s", VisorBinaryVersion)

	if err := gen.initNode(logger, visorBinaryPath, vegaBinaryPath); err != nil {
		return fmt.Errorf("failed to init vega node: %w", err)
	}

	if err := gen.prepareVisorHome(logger); err != nil {
		return fmt.Errorf("failed to prepare visor home: %w", err)
	}

	if err := gen.copyBinaries(logger, vegaBinaryPath, visorBinaryPath); err != nil {
		return fmt.Errorf("failed to copy binaries to visor home: %w", err)
	}

	restartSnapshot, err := gen.selectSnapshotForRestart(logger)
	if err != nil {
		return fmt.Errorf("failed to select snapshot for restart: %w", err)
	}

	if err := gen.updateConfigs(logger, restartSnapshot); err != nil {
		return fmt.Errorf("failed to update config files for the node: %w", err)
	}

	if err := gen.downloadGenesis(logger); err != nil {
		return fmt.Errorf("failed to download genesis: %w", err)
	}
	return nil
}

func (gen *DataNodeGenerator) downloadGenesis(logger *zap.SugaredLogger) error {
	genesisDestination := filepath.Join(gen.userSettings.TendermintHome, vegacmd.GenesisPath)
	logger.Infof("Downloading genesis.json file from %s", gen.networkConfig.GenesisURL)
	if err := utils.DownloadFile(gen.networkConfig.GenesisURL, genesisDestination); err != nil {
		return fmt.Errorf("failed to download genesis: %w", err)
	}
	logger.Infof("Genesis downloaded to %s", genesisDestination)

	return nil
}

func (gen *DataNodeGenerator) copyBinaries(
	logger *zap.SugaredLogger,
	vegaBinaryPath, visorBinaryPath string,
) error {
	vegavisorDstFilePath := filepath.Join(gen.userSettings.VisorHome, "visor")
	logger.Infof("Copying vegavisor from %s to %s", visorBinaryPath, vegavisorDstFilePath)
	if err := utils.CopyFile(visorBinaryPath, vegavisorDstFilePath); err != nil {
		return fmt.Errorf("failed to copy visor binary: %w", err)
	}
	logger.Info("Visor binary copied")

	version := gen.userSettings.VegaBinaryVersion
	if gen.userSettings.Mode == StartFromBlock0 {
		version = "genesis"
	}

	vegaDstFilePath := filepath.Join(gen.userSettings.VisorHome, version, "vega")
	logger.Infof("Copying vega from %s to %s", vegaBinaryPath, vegaDstFilePath)
	if err := utils.CopyFile(vegaBinaryPath, vegaDstFilePath); err != nil {
		return fmt.Errorf("failed to copy vega binary: %w", err)
	}
	logger.Info("Vega binary copied")

	versionDirectory := filepath.Join(gen.userSettings.VisorHome, version)
	currentDirectory := filepath.Join(gen.userSettings.VisorHome, "current")
	logger.Infof("Creating symlink from %s to %s", versionDirectory, currentDirectory)
	if err := os.Symlink(versionDirectory, currentDirectory); err != nil {
		return fmt.Errorf(
			"failed to create symlink from %s to %s: %w",
			versionDirectory,
			currentDirectory,
			err,
		)
	}
	logger.Info("Symlink created")

	return nil
}

func (gen *DataNodeGenerator) prepareVisorHome(logger *zap.SugaredLogger) error {
	runConfigDirPath := filepath.Join(gen.userSettings.VisorHome, gen.userSettings.VegaBinaryVersion)
	version := gen.userSettings.VegaBinaryVersion

	if gen.userSettings.Mode == StartFromBlock0 {
		runConfigDirPath = filepath.Join(gen.userSettings.VisorHome, "genesis")
		version = "genesis"
	}

	logger.Infof("Preparing %s folder for vega", runConfigDirPath)
	if err := os.MkdirAll(runConfigDirPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to make directory: %w", err)
	}
	logger.Infof("Folder %s created", runConfigDirPath)

	runConfigPath := filepath.Join(runConfigDirPath, "run-config.toml")
	logger.Infof("Preparing run-config toml file in %s", runConfigPath)
	runConfigContent, err := vegacmd.TemplateVisorRunConfig(
		version,
		gen.userSettings.VegaHome,
		gen.userSettings.TendermintHome,
	)
	if err != nil {
		return fmt.Errorf("failed to generate run-config.toml from template: %w", err)
	}
	if err := os.WriteFile(runConfigPath, []byte(runConfigContent), os.ModePerm); err != nil {
		return fmt.Errorf("failed to write run-config.toml in %s: %w", runConfigContent, err)
	}
	logger.Infof("The run-config.toml file saved in %s", runConfigPath)

	return nil
}

func (gen *DataNodeGenerator) updateConfigs(
	logger *zap.SugaredLogger,
	restartSnapshot *types.CoreSnapshot,
) error {
	dataNodeConfig := map[string]interface{}{
		"SQLStore.RetentionPeriod":                    gen.userSettings.DataRetention,
		"SQLStore.ConnectionConfig.Host":              gen.userSettings.SQLCredentials.Host,
		"SQLStore.ConnectionConfig.Port":              gen.userSettings.SQLCredentials.Port,
		"SQLStore.ConnectionConfig.Username":          gen.userSettings.SQLCredentials.User,
		"SQLStore.ConnectionConfig.Password":          gen.userSettings.SQLCredentials.Pass,
		"SQLStore.ConnectionConfig.Database":          gen.userSettings.SQLCredentials.DatabaseName,
		"SQLStore.WipeOnStartup":                      true,
		"NetworkHistory.Store.BootstrapPeers":         gen.networkConfig.BootstrapPeers,
		"NetworkHistory.Initialise.MinimumBlockCount": gen.userSettings.NetworkHistoryMinBlockCount,
		"NetworkHistory.Initialise.Timeout":           "4h",
		"NetworkHistory.RetryTimeout":                 "15s",
		"API.RateLimit.Rate":                          300.0,
		"API.RateLimit.Burst":                         1000,
		// This is controversial for vega but most of the people does not care about network history
		"NetworkHistory.Publish": false,
	}

	vegaConfig := map[string]interface{}{
		"Snapshot.StartHeight":      -1,
		"Broker.Socket.Enabled":     true,
		"Broker.Socket.DialTimeout": "4h",
	}

	tendermintConfig := map[string]interface{}{
		"p2p.seeds":              strings.Join(gen.networkConfig.TendermintSeeds, ","),
		"p2p.persistent_peers":   strings.Join(gen.networkConfig.TendermintPersistentPeers, ","),
		"pex":                    true,
		"statesync.enable":       false,
		"statesync.rpc_servers":  strings.Join(gen.networkConfig.TendermintRPCServers, ","),
		"statesync.trust_period": "672h0m0s",
	}

	vegavisorConfig := map[string]interface{}{
		"maxNumberOfFirstConnectionRetries": 43200,
		"autoInstall.enabled":               true,
		"autoInstall.repositoryOwner":       strings.Split(gen.networkConfig.Repository, "/")[0],
		"autoInstall.repository":            strings.Split(gen.networkConfig.Repository, "/")[1],
		"autoInstall.asset.name": fmt.Sprintf(
			"vega-%s-%s.zip",
			runtime.GOOS,
			runtime.GOARCH,
		),
		"autoInstall.asset.binaryName": "vega",
	}

	if gen.userSettings.Mode == StartFromNetworkHistory {
		if restartSnapshot == nil {
			return fmt.Errorf(
				"failed to start node from network history: no selected snapshot for restart",
			)
		}

		if restartSnapshot.BlockHash == "" {
			return fmt.Errorf(
				"cannot start vega from the network-history when latest snapshot is empty",
			)
		}

		trustHeight, err := strconv.Atoi(restartSnapshot.BlockHeight)
		if err != nil {
			return fmt.Errorf("failed to convert trust block height from string to int: %w", err)
		}

		// We cannot use statis StartHeight value because it is not working when we are syncing more blocks from the data-node
		// Tendermint does not offer more than 10 snapshots.
		// vegaConfig["Snapshot.StartHeight"] = trustHeight
		dataNodeConfig["AutoInitialiseFromNetworkHistory"] = true
		tendermintConfig["statesync.enable"] = true
		tendermintConfig["statesync.trust_height"] = trustHeight
		tendermintConfig["statesync.trust_hash"] = restartSnapshot.BlockHash
	}

	dataNodeConfigPath := filepath.Join(gen.userSettings.DataNodeHome, vegacmd.DataNodeConfigPath)
	logger.Infof(
		"Updating data-node config(%s). New parameters: %v",
		dataNodeConfigPath,
		dataNodeConfig,
	)
	if err := utils.UpdateConfig(dataNodeConfigPath, "toml", dataNodeConfig); err != nil {
		return fmt.Errorf("failed to update the data-node config; %w", err)
	}
	logger.Info("Data-node config updated")

	vegaConfigPath := filepath.Join(gen.userSettings.VegaHome, vegacmd.CoreConfigPath)
	logger.Infof("Updating vega-core config(%s). New parameters: %v", vegaConfigPath, vegaConfig)
	if err := utils.UpdateConfig(vegaConfigPath, "toml", vegaConfig); err != nil {
		return fmt.Errorf("failed to update the vega config; %w", err)
	}
	logger.Info("Vega-core config updated")

	tendermintConfigPath := filepath.Join(
		gen.userSettings.TendermintHome,
		vegacmd.TenderminConfigPath,
	)
	logger.Infof(
		"Updating tendermint config(%s). New parameters: %v",
		tendermintConfigPath,
		tendermintConfig,
	)
	if err := utils.UpdateConfig(tendermintConfigPath, "toml", tendermintConfig); err != nil {
		return fmt.Errorf("failed to update the tendermint config; %w", err)
	}
	logger.Info("Tendermint config updated")

	vegavisorConfigPath := filepath.Join(gen.userSettings.VisorHome, vegacmd.VegavisorConfigPath)
	logger.Infof(
		"Updating vegavisor config(%s). New parameters: %v",
		vegavisorConfigPath,
		vegavisorConfig,
	)
	if err := utils.UpdateConfig(vegavisorConfigPath, "toml", vegavisorConfig); err != nil {
		return fmt.Errorf("failed to update vegavisor config: %w", err)
	}
	logger.Info("Vegavisor config updated")

	return nil
}

func (gen *DataNodeGenerator) selectSnapshotForRestart(
	logger *zap.SugaredLogger,
) (*types.CoreSnapshot, error) {
	if gen.userSettings.Mode != StartFromNetworkHistory {
		return &types.CoreSnapshot{}, nil
	}

	logger.Info("Fetching network snapshots")
	snapshots, err := vegaapi.Snapshots(gen.networkConfig.DataNodesRESTUrls)
	if err != nil {
		return nil, fmt.Errorf("failed to get core snapshot for trusted block: %w", err)
	}

	logger.Infof("Found %d snapshots", len(snapshots.CoreSnapshots.Edges))
	if len(snapshots.CoreSnapshots.Edges) < 3 {
		return nil, fmt.Errorf(
			"not enough snapshots for restart: required at least 3 snapshots, %d got",
			len(snapshots.CoreSnapshots.Edges),
		)
	}

	logger.Info("Fetching network history segments")
	segments, err := vegaapi.NetworkHistorySegments(gen.networkConfig.DataNodesRESTUrls)
	if err != nil {
		return nil, fmt.Errorf("failed to get network-history segments: %w", err)
	}

	logger.Infof("Found %d network-history segments", len(segments.Segments))
	if len(segments.Segments) < 3 {
		return nil, fmt.Errorf(
			"not enough network history segments for restart: required at least 3 segments, %d got",
			len(segments.Segments),
		)
	}

	logger.Info("Finding snapshot for restart")
	snapshotList := []types.CoreSnapshot{}
	for _, snapshot := range snapshots.CoreSnapshots.Edges {
		// cut the invalid snapshots out
		if snapshot.Node.BlockHash == "" || snapshot.Node.BlockHeight == "" {
			continue
		}

		snapshotList = append(snapshotList, snapshot.Node)
	}

	segmentList := []types.NetworkHistorySegment{}
	for _, segment := range segments.Segments {
		if segment.ToHeight == "" {
			continue
		}

		segmentList = append(segmentList, segment)
	}

	// sort lists from the highest to the lowest
	sort.Slice(snapshotList, func(i, j int) bool {
		iHeight, _ := strconv.Atoi(snapshotList[i].BlockHeight)
		jHeight, _ := strconv.Atoi(snapshotList[j].BlockHeight)

		return iHeight > jHeight
	})

	sort.Slice(segmentList, func(i, j int) bool {
		iHeight, _ := strconv.Atoi(segmentList[i].ToHeight)
		jHeight, _ := strconv.Atoi(segmentList[j].ToHeight)

		return iHeight > jHeight
	})

	if len(snapshotList) < 3 {
		return nil, fmt.Errorf("not enough snapshots for restart after filtering")
	}

	if len(segmentList) < 3 {
		return nil, fmt.Errorf("not enough segments for restart after filtering")
	}

	// select 3-rd highest segment for restart(latest segments may noy be published to the IPFS yet)
	selectedSegment := segmentList[2]
	selectedSegmentHeight, err := strconv.Atoi(selectedSegment.ToHeight)
	if err != nil {
		return nil, fmt.Errorf("failed to convert height for selected segment to int: %w", err)
	}

	var selectedSnapshot *types.CoreSnapshot
	// select first snapshot with lower or equal block than 3-rd highest segment
	for idx, snapshot := range snapshotList {
		snapshotHeight, err := strconv.Atoi(snapshot.BlockHeight)
		if err != nil {
			continue // TODO: Maybe we should handle it???
		}
		if snapshotHeight <= selectedSegmentHeight {
			selectedSnapshot = &snapshotList[idx]
			break
		}
	}

	if selectedSnapshot == nil {
		return nil, fmt.Errorf(
			"failed to find snapshot lower than block %s (3-rd highest segment)",
			selectedSegment.ToHeight,
		)
	}

	logger.Infof("Selected snapshot for restart at block %s", selectedSnapshot.BlockHeight)

	return selectedSnapshot, nil
}

func (gen *DataNodeGenerator) initNode(
	logger *zap.SugaredLogger,
	visorBinary, vegaBinary string,
) error {
	logger.Infof("Initializing vegavisor in the %s", gen.userSettings.VisorHome)
	if err := vegacmd.InitVisor(visorBinary, gen.userSettings.VisorHome); err != nil {
		return fmt.Errorf(
			"failed to initialize vegavisor in %s: %w",
			gen.userSettings.VisorHome,
			err,
		)
	}
	logger.Info("Visor successfully initialized")

	logger.Infof("Initializing tendermint in the %s", gen.userSettings.TendermintHome)
	if err := vegacmd.InitTendermint(vegaBinary, gen.userSettings.TendermintHome); err != nil {
		return fmt.Errorf(
			"failed to initialize tendermint in %s: %w",
			gen.userSettings.TendermintHome,
			err,
		)
	}
	logger.Info("Tendermint successfully initialized")

	logger.Infof("Initializing vega in the %s", gen.userSettings.VegaHome)
	if err := vegacmd.InitVega(vegaBinary, gen.userSettings.VegaHome, vegacmd.VegaNodeFull); err != nil {
		return fmt.Errorf(
			"failed to initialize vega in %s: %w",
			gen.userSettings.VegaHome,
			err,
		)
	}
	logger.Info("Visor successfully initialized")

	logger.Infof("Initializing data-node n the %s", gen.userSettings.DataNodeHome)
	if err := vegacmd.InitDataNode(vegaBinary, gen.userSettings.DataNodeHome, gen.userSettings.VegaChainId); err != nil {
		return fmt.Errorf(
			"failed to initialize data-node in %s: %w",
			gen.userSettings.DataNodeHome,
			err,
		)
	}
	logger.Info("Data-node successfully initialized")

	return nil
}
