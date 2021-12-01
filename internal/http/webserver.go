package http

import (
	"embed"
	"encoding/json"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/terran-stakers/terranseed/internal/geoloc"
	"github.com/terran-stakers/terranseed/internal/seednode"
	"io/fs"
	"net/http"
	"os"
)

var (
	logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "web")
)

type WebResources struct {
	Res   embed.FS
	Files map[string]string
}

func StartWebServer(seedConfig seednode.TSConfig, webResources fs.FS, ips *[]geoloc.GeolocalizedPeers) {
	// serve static assets
	http.Handle("/", http.FileServer(http.FS(webResources)))

	// serve endpoint
	http.HandleFunc("/api/peers", writePeers)

	// start web server in non-blocking
	go func() {
		err := http.ListenAndServe(":"+seedConfig.HttpPort, nil)
		logger.Info("HTTP Server started", "port", seedConfig.HttpPort)
		if err != nil {
			panic(err)
		}
	}()
}

func writePeers(w http.ResponseWriter, r *http.Request) {
	marshal, err := json.Marshal(&geoloc.ResolvedPeers)
	if err != nil {
		logger.Info("Failed to marshal peers list")
		return
	}
	_, err = w.Write(marshal)
	if err != nil {
		return
	}
}
