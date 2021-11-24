package main

import (
	"embed"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/terran-stakers/terranseed/internal/http"
	"github.com/terran-stakers/terranseed/internal/tendermint"
	"os"
)

var (
	logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "main")

	//go:embed web
	res   embed.FS
	files = map[string]string{
		"/":                "web/templates/index.html",
		"/assets/gmaps.js": "assets/js/gmaps.js",
		"/assets/main.css": "assets/css/main.css",
	}
)

// DefaultConfig returns a seed config initialized with default values
func DefaultConfig() *tendermint.Config {
	return &tendermint.Config{
		ListenAddress:       "tcp://0.0.0.0:6969",
		HttpPort:            "8888",
		ChainID:             "osmosis-1",
		NodeKeyFile:         "node_key.json",
		AddrBookFile:        "addrbook.json",
		AddrBookStrict:      true,
		MaxNumInboundPeers:  3000,
		MaxNumOutboundPeers: 1000,
		Seeds:               "1b077d96ceeba7ef503fb048f343a538b2dcdf1b@136.243.218.244:26656,2308bed9e096a8b96d2aa343acc1147813c59ed2@3.225.38.25:26656,085f62d67bbf9c501e8ac84d4533440a1eef6c45@95.217.196.54:26656,f515a8599b40f0e84dfad935ba414674ab11a668@osmosis.blockpane.com:26656",
	}
}

func WebResources() *http.WebResources {
	return &http.WebResources{
		res,
		files,
	}
}

func main() {
	seedConfig := DefaultConfig()
	webResources := WebResources();

	tendermint.InitConfig(seedConfig)

	logger.Info("Starting Seed Node...")
//	tendermint.StartSeedNode(*seedConfig)

	logger.Info("Starting Web Server...")
	http.StartWebServer(*seedConfig, *webResources)
}
