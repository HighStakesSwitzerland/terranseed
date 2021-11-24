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
	logger.Info("test")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		page, ok := webResources.Files[r.URL.Path]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		tpl, err := template.ParseFS(webResources.Res, page)
		if err != nil {
			logger.Info("page {} not found in pages cache...", r.RequestURI)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		data := map[string]interface{}{
			"userAgent": r.UserAgent(),
		}
		if err := tpl.Execute(w, data); err != nil {
			return
		}
	})
	http.FileServer(http.FS(webResources.Res))
	logger.Info("HTTP Server started", "port", seedConfig.HttpPort)
	err := http.ListenAndServe(":"+seedConfig.HttpPort, nil)
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
