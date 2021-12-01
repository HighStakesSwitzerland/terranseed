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

func main() {
	seedConfig, nodeKey := seednode.InitConfig()
	embeddedFS, _ := fs.Sub(res, "dist/terranseed")

	logger.Info("Starting Seed Node...")

	sw := seednode.StartSeedNode(seedConfig, nodeKey)

	logger.Info("Starting Web Server...")
	http.StartWebServer(seedConfig, embeddedFS, &geolocalizedIps)

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
