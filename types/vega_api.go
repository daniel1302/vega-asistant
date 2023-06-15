package types

type VegaStatistics struct {
	Statistics struct {
		ChainID    string `json:"chainId"`
		AppVersion string `json:"appVersion"`
	} `json:"statistics"`
}
