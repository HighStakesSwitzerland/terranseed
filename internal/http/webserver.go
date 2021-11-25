package http

import (
	"embed"
	"encoding/json"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/terran-stakers/terranseed/internal/geoloc"
	"github.com/terran-stakers/terranseed/internal/seednode"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

var (
	logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "web")
)

type WebResources struct {
	Res   embed.FS
	Files map[string]string
}

func StartWebServer(seedConfig seednode.Config, webResources WebResources, ips *[]geoloc.GeolocalizedPeers) {
	// serve static assets
	fs := http.FileServer(http.Dir("./web/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	// serve html files
	http.HandleFunc("/", serveTemplate)

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

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	index := filepath.Join("./web/templates", "index.html")
	templates := filepath.Join("./web/templates", filepath.Clean(r.URL.Path))
	logger.Info("index", "i", index, "t", templates)

	// Return a 404 if the template doesn't exist
	fileInfo, err := os.Stat(templates)

	if err != nil || fileInfo.IsDir() {
		http.Redirect(w, r, "/index.html", 302)
		return
	}

	tmpl, err := template.ParseFiles(index, templates)
	if err != nil {
		// Log the detailed error
		logger.Info(err.Error())
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
	}

	err = tmpl.ExecuteTemplate(w, "index", nil)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
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
