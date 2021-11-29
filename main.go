package main

import (
	"embed"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/p2p"
	"github.com/terran-stakers/terranseed/internal/geoloc"
	"github.com/terran-stakers/terranseed/internal/http"
	"github.com/terran-stakers/terranseed/internal/seednode"
	"io/fs"
	"os"
	"time"
)

var (
	logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "main")

	//go:embed dist/terranseed
	res embed.FS

	geolocalizedIps = make([]geoloc.GeolocalizedPeers, 0)
)

// DefaultConfig returns a seed config initialized with default values
func DefaultConfig() *seednode.Config {
	return &seednode.Config{
		ConfigDir:           ".terranseed",
		ListenAddress:       "tcp://0.0.0.0:6969",
		HttpPort:            "8080",
		ChainID:             "columbus-5",
		NodeKeyFile:         "node_key.json",
		AddrBookFile:        "addrbook.json",
		AddrBookStrict:      true,
		MaxNumInboundPeers:  3000,
		MaxNumOutboundPeers: 1000,
		Seeds:               "04673893a2cf451604b5fd68e064c4e9fe7c410c@13.125.144.21:26656,6d8e943c049a80c161a889cb5fcf3d184215023e@101.101.208.201:26656,87048bf71526fb92d73733ba3ddb79b7a83ca11e@public-seed.terra.dev:26656,6d8e943c049a80c161a889cb5fcf3d184215023e@public-seed2.terra.dev:26656,e999fc20aa5b87c1acef8677cf495ad85061cfb9@seed.terra.delightlabs.io:26656",
	}
}

func main() {
	seedConfig := DefaultConfig()
	seednode.InitConfig(seedConfig)
	embeddedFS, _ := fs.Sub(res, "dist/terranseed")

	logger.Info("Starting Seed Node...")

	sw := seednode.StartSeedNode(*seedConfig)

	logger.Info("Starting Web Server...")
	http.StartWebServer(*seedConfig, embeddedFS, &geolocalizedIps)

	StartGeolocServiceAndBlock(sw)
}

func StartGeolocServiceAndBlock(sw *p2p.Switch) {
	// Fire periodically
	ticker := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-ticker.C:
			peers := seednode.GetPeers(sw.Peers().List())
			geoloc.ResolveIps(peers)
		}
	}
}
