package http

import (
	"embed"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/terran-stakers/terranseed/internal/tendermint"
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

func StartWebServer(seedConfig tendermint.Config, webResources WebResources) {
	// serve static assets
	fs := http.FileServer(http.Dir("./web/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	// serve html files
	http.HandleFunc("/", serveTemplate)

	// start web server in non-blocking
	err := http.ListenAndServe(":"+seedConfig.HttpPort, nil)
	logger.Info("HTTP Server started", "port", seedConfig.HttpPort)
	if err != nil {
		panic(err)
	}
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
		logger.Error(err.Error())
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
