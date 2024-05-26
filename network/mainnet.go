package network

import "github.com/daniel1302/vega-assistant/types"

type BinaryOverride struct {
	OldVersion string
	NewVersion string
	Block      uint64
}

type NetworkConfig struct {
	GenesisVersion            string
	Repository                string
	GenesisURL                string
	LowestVisorVersion        string
	DataNodesRESTUrls         []string
	TendermintSeeds           []string
	BootstrapPeers            []types.EndpointWithVegaREST
	TendermintRPCServers      []types.EndpointWithVegaREST
	TendermintPersistentPeers []string
	BinariesOverride          []BinaryOverride
}

func MainnetConfig() NetworkConfig {
	return NetworkConfig{
		GenesisVersion:     "v0.71.4",
		LowestVisorVersion: "v0.73.6",
		Repository:         "vegaprotocol/vega",
		GenesisURL:         "https://raw.githubusercontent.com/vegaprotocol/networks/master/mainnet1/genesis.json",
		DataNodesRESTUrls: []string{
			// "https://api0.vega.community",
			"https://api1.vega.community",
			"https://api2.vega.community",
			"https://api3.vega.community",
		},
		TendermintSeeds: []string{
			"b0db58f5651c85385f588bd5238b42bedbe57073@13.125.55.240:26656",
			"abe207dae9367995526812d42207aeab73fd6418@18.158.4.175:26656",
			"198ecd046ebb9da0fc5a3270ee9a1aeef57a76ff@144.76.105.240:26656",
			"211e435c2162aedb6d687409d5d7f67399d198a9@65.21.60.252:26656",
			"c5b11e1d819115c4f3974d14f76269e802f3417b@34.88.191.54:26656",
			"61051c21f083ee30c835a34a0c17c5d1ceef3c62@51.178.75.45:26656",
			"b0db58f5651c85385f588bd5238b42bedbe57073@18.192.52.234:26656",
			"36a2ca7bb6a50427be2181c8ebb7f62ac62ebaf5@m2.vega.community:26656",
			"9903c02a0ff881dc369fc7daccb22c1f9680d2dd@api0.vega.community:26656",
			"9903c02a0ff881dc369fc7daccb22c1f9680d2dd@api0.vega.community:26656",
			"32d7380b195c088c0605c5d24bcf15ff1dade05f@api1.vega.community:26656",
			"4f26ec99d3cf6f0e9e973c0a5f3da87d89ec6677@api2.vega.community:26656",
			"eafacd11af53cd9fb2a14eada53485779cbee4ab@api3.vega.community:26656",
			"9de3ca2bbeb62d165d39acbbcf174e7ac3a6b7c9@be3.vega.community:26656",
		},
		TendermintRPCServers: []types.EndpointWithVegaREST{
			{REST: "https://api0.vega.community", Endpoint: "api0.vega.community:26657"},
			{REST: "https://api1.vega.community", Endpoint: "api1.vega.community:26657"},
			{REST: "https://api2.vega.community", Endpoint: "api2.vega.community:26657"},
			{REST: "https://api3.vega.community", Endpoint: "api3.vega.community:26657"},
		},
		BootstrapPeers: []types.EndpointWithVegaREST{
			{REST: "https://api0.vega.community", Endpoint: "/dns/api0.vega.community/tcp/4001/ipfs/12D3KooWAHkKJfX7rt1pAuGebP9g2BGTT5w7peFGyWd2QbpyZwaw"},
			{REST: "https://api1.vega.community", Endpoint: "/dns/api1.vega.community/tcp/4001/ipfs/12D3KooWDZrusS1p2XyJDbCaWkVDCk2wJaKi6tNb4bjgSHo9yi5Q"},
			{REST: "https://api2.vega.community", Endpoint: "/dns/api2.vega.community/tcp/4001/ipfs/12D3KooWEH9pQd6P7RgNEpwbRyavWcwrAdiy9etivXqQZzd7Jkrh"},
			{REST: "https://api3.vega.community", Endpoint: "/dns/api3.vega.community/tcp/4001/ipfs/12D3KooWHSoYzEqSfUWEXfFbSnmRhWcP2WgZG2GRT8fzZzio5BTY"},
		},
		TendermintPersistentPeers: []string{
			// "55b8ac477ddd6c0c9bae411dfa6ebfb46e7b4022@veganodeoperator.com:26656",
			// "2d1bbf1229bd7f8e57e89c61346ab6928d61881b@127.0.0.1:26656",
			// "7f7735c30a6cbc70daab5bdf7f9ebe77b662e4aa@be0.vega.community:26656",
			// "e1bbd644b509aacbcc5d5b47692c15297fc7fb50@be1.vega.community:26656",
			// "9de3ca2bbeb62d165d39acbbcf174e7ac3a6b7c9@be3.vega.community:26656",
			// "9903c02a0ff881dc369fc7daccb22c1f9680d2dd@api0.vega.community:26656",
			// "32d7380b195c088c0605c5d24bcf15ff1dade05f@api1.vega.community:26656",
			// "4f26ec99d3cf6f0e9e973c0a5f3da87d89ec6677@api2.vega.community:26656",
			// "eafacd11af53cd9fb2a14eada53485779cbee4ab@api3.vega.community:26656",
			// "9de3ca2bbeb62d165d39acbbcf174e7ac3a6b7c9@be3.vega.community:26656",
		},
		BinariesOverride: []BinaryOverride{
			{
				OldVersion: "v0.75.8",
				NewVersion: "v0.75.8-fix.2",
				Block:      47865000,
			},
		},
	}
}
