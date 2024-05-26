package types

import "time"

type VegaRawStatistics struct {
	Statistics struct {
		ChainID     string `json:"chainId"`
		AppVersion  string `json:"appVersion"`
		CurrentTime string
		VegaTime    string
		BlockHeight string
	} `json:"statistics"`
}

type VegaStatistics struct {
	BlockHeight    uint64
	DataNodeHeight uint64

	CurrentTime time.Time
	VegaTime    time.Time

	ChainID    string
	AppVersion string
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
