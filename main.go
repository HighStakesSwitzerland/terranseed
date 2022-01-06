package main

import (
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/p2p"
	"github.com/terran-stakers/terranseed/internal/geoloc"
	"github.com/terran-stakers/terranseed/internal/http"
	"github.com/terran-stakers/terranseed/internal/seednode"
	"os"
	"time"
)

var (
	logger          = log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "main")
	geolocalizedIps = new([]geoloc.GeolocalizedPeers)
)

func main() {
	seedConfig, nodeKey := seednode.InitConfig()

	logger.Info("Starting Web Server...")
	http.StartWebServer(seedConfig, geolocalizedIps)

	logger.Info("Starting Seed Node...")
	sw := seednode.StartSeedNode(seedConfig, nodeKey)

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
