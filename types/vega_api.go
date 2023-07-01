package types

type VegaStatistics struct {
	Statistics struct {
		ChainID    string `json:"chainId"`
		AppVersion string `json:"appVersion"`
	} `json:"statistics"`
}

type CoreSnapshot struct {
	CoreVersion string `json:"coreVersion"`
	BlockHeight string `json:"blockHeight"`
	BlockHash   string `json:"blockHash"`
}

type CoreSnapshots struct {
	CoreSnapshots struct {
		Edges []struct {
			Node CoreSnapshot `json:"node"`
		} `json:"edges"`
	} `json:"coreSnapshots"`
}

type NetworkHistorySegment struct {
	FromHeight       string `json:"fromHeight"`
	ToHeight         string `json:"toHeight"`
	HistorySegmentId string `json:"historySegmentId"`
}

type NetworkHistorySegments struct {
	Segments []NetworkHistorySegment `json:"segments"`
}
