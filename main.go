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
	res   embed.FS

	geolocalizedIps = make([]geoloc.GeolocalizedPeers, 0)
)

// DefaultConfig returns a seed config initialized with default values
func DefaultConfig() *seednode.Config {
	return &seednode.Config{
		ListenAddress:       "tcp://0.0.0.0:6969",
		HttpPort:            "8080",
		ChainID:             "osmosis-1",
		NodeKeyFile:         "node_key.json",
		AddrBookFile:        "addrbook.json",
		AddrBookStrict:      true,
		MaxNumInboundPeers:  3000,
		MaxNumOutboundPeers: 1000,
		Seeds:               "1b077d96ceeba7ef503fb048f343a538b2dcdf1b@136.243.218.244:26656,2308bed9e096a8b96d2aa343acc1147813c59ed2@3.225.38.25:26656,085f62d67bbf9c501e8ac84d4533440a1eef6c45@95.217.196.54:26656,f515a8599b40f0e84dfad935ba414674ab11a668@osmosis.blockpane.com:26656",
	}
}

func main() {
	seedConfig := DefaultConfig()
	seednode.InitConfig(seedConfig)
  embeddedFS, _ := fs.Sub(res, "dist/terranseed")

  logger.Info("Starting Seed Node...")

	// sw := seednode.StartSeedNode(*seedConfig)

	logger.Info("Starting Web Server...")
	http.StartWebServer(*seedConfig, embeddedFS, &geolocalizedIps)

	// StartGeolocServiceAndBlock(sw)
}

func StartGeolocServiceAndBlock(sw *p2p.Switch) {
	// Fire periodically
	ticker := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-ticker.C:
			peers := seednode.GetPeers(sw)
			geoloc.ResolveIps(peers)
		}
	}
}
