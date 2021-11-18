package main

import (
	"embed"
	"fmt"
	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	tmstrings "github.com/tendermint/tendermint/libs/strings"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/p2p/pex"
	"github.com/tendermint/tendermint/version"
	"html/template"
	oslog "log"
	"net/http"
	"os"
	"path/filepath"
	"googlemaps.github.io/maps"
)

var (
	//https://socketloop.com/tutorials/golang-find-location-by-ip-address-and-display-with-google-map
	//go:embed assets web
	res   embed.FS
	files = map[string]string{
		"/": "web/index.html",
		"/assets/gmaps.js": "assets/gmaps.js",
		"/assets/main.css": "assets/main.css",
	}
)

// Config defines the configuration format
type Config struct {
	ListenAddress       string `toml:"laddr" comment:"Address to listen for incoming connections"`
	ChainID             string `toml:"chain_id" comment:"network identifier (todo move to cli flag argument? keeps the config network agnostic)"`
	NodeKeyFile         string `toml:"node_key_file" comment:"path to node_key (relative to tendermint-seed home directory or an absolute path)"`
	AddrBookFile        string `toml:"addr_book_file" comment:"path to address book (relative to tendermint-seed home directory or an absolute path)"`
	AddrBookStrict      bool   `toml:"addr_book_strict" comment:"Set true for strict routability rules\n Set false for private or local networks"`
	MaxNumInboundPeers  int    `toml:"max_num_inbound_peers" comment:"maximum number of inbound connections"`
	MaxNumOutboundPeers int    `toml:"max_num_outbound_peers" comment:"maximum number of outbound connections"`
	Seeds               string `toml:"seeds" comment:"seed nodes we can use to discover peers"`
}

// DefaultConfig returns a seed config initialized with default values
func DefaultConfig() *Config {
	return &Config{
		ListenAddress:       "tcp://0.0.0.0:6969",
		ChainID:             "osmosis-1",
		NodeKeyFile:         "config/node_key.json",
		AddrBookFile:        "data/addrbook.json",
		AddrBookStrict:      true,
		MaxNumInboundPeers:  1000,
		MaxNumOutboundPeers: 1000,
		Seeds:               "1b077d96ceeba7ef503fb048f343a538b2dcdf1b@136.243.218.244:26656,2308bed9e096a8b96d2aa343acc1147813c59ed2@3.225.38.25:26656,085f62d67bbf9c501e8ac84d4533440a1eef6c45@95.217.196.54:26656,f515a8599b40f0e84dfad935ba414674ab11a668@osmosis.blockpane.com:26656",
	}
}

func startWebServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		page, ok := files[r.URL.Path]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		tpl, err := template.ParseFS(res, page)
		if err != nil {
			oslog.Printf("page %s not found in pages cache...", r.RequestURI)
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
	http.FileServer(http.FS(res))
	oslog.Println("server started...")
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		panic(err)
	}
}

func main() {

	startWebServer()

/*	idOverride := os.Getenv("ID")
	seedOverride := os.Getenv("SEEDS")
	userHomeDir, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	homeDir := filepath.Join(userHomeDir, ".terranseed")
	configFile := "config/config.toml"
	configFilePath := filepath.Join(homeDir, configFile)
	MkdirAllPanic(filepath.Dir(configFilePath), os.ModePerm)
	SeedConfig := DefaultConfig()
	if idOverride != "" {
		SeedConfig.ChainID = idOverride
	}
	if seedOverride != "" {
		SeedConfig.Seeds = seedOverride
	}
	Start(*SeedConfig)*/
}

// MkdirAllPanic invokes os.MkdirAll but panics if there is an error
func MkdirAllPanic(path string, perm os.FileMode) {
	err := os.MkdirAll(path, perm)
	if err != nil {
		panic(err)
	}
}

// Start starts a Tenderseed
func Start(SeedConfig Config) {
	logger := log.NewTMLogger(
		log.NewSyncWriter(os.Stdout),
	)

	chainID := SeedConfig.ChainID
	nodeKeyFilePath := SeedConfig.NodeKeyFile
	addrBookFilePath := SeedConfig.AddrBookFile

	MkdirAllPanic(filepath.Dir(nodeKeyFilePath), os.ModePerm)
	MkdirAllPanic(filepath.Dir(addrBookFilePath), os.ModePerm)

	cfg := config.DefaultP2PConfig()
	cfg.AllowDuplicateIP = true

	// allow a lot of inbound peers since we disconnect from them quickly in seed mode
	cfg.MaxNumInboundPeers = 3000

	// keep trying to make outbound connections to exchange peering info
	cfg.MaxNumOutboundPeers = 400

	nodeKey, err := p2p.LoadOrGenNodeKey(nodeKeyFilePath)
	if err != nil {
		panic(err)
	}

	logger.Info("terranseed",
		"key", nodeKey.ID(),
		"listen", SeedConfig.ListenAddress,
		"chain", chainID,
		"strict-routing", SeedConfig.AddrBookStrict,
		"max-inbound", SeedConfig.MaxNumInboundPeers,
		"max-outbound", SeedConfig.MaxNumOutboundPeers,
	)

	// TODO(roman) expose per-module log levels in the config
	filteredLogger := log.NewFilter(logger, log.AllowInfo())

	protocolVersion :=
		p2p.NewProtocolVersion(
			version.P2PProtocol,
			version.BlockProtocol,
			0,
		)

	// NodeInfo gets info on yhour node
	nodeInfo := p2p.DefaultNodeInfo{
		ProtocolVersion: protocolVersion,
		DefaultNodeID:   nodeKey.ID(),
		ListenAddr:      SeedConfig.ListenAddress,
		Network:         chainID,
		Version:         "0.6.9",
		Channels:        []byte{pex.PexChannel},
		Moniker:         fmt.Sprintf("%s-seed", chainID),
	}

	addr, err := p2p.NewNetAddressString(p2p.IDAddressString(nodeInfo.DefaultNodeID, nodeInfo.ListenAddr))
	if err != nil {
		panic(err)
	}

	transport := p2p.NewMultiplexTransport(nodeInfo, *nodeKey, p2p.MConnConfig(cfg))
	if err := transport.Listen(*addr); err != nil {
		panic(err)
	}

	book := pex.NewAddrBook(addrBookFilePath, SeedConfig.AddrBookStrict)
	book.SetLogger(filteredLogger.With("module", "book"))

	pexReactor := pex.NewReactor(book, &pex.ReactorConfig{
		SeedMode: true,
		Seeds:    tmstrings.SplitAndTrim(SeedConfig.Seeds, ",", " "),
	})
	pexReactor.SetLogger(filteredLogger.With("module", "pex"))

	sw := p2p.NewSwitch(cfg, transport)
	sw.SetLogger(filteredLogger.With("module", "switch"))
	sw.SetNodeKey(nodeKey)
	sw.SetAddrBook(book)
	sw.AddReactor("pex", pexReactor)

	// last
	sw.SetNodeInfo(nodeInfo)

	tmos.TrapSignal(logger, func() {
		logger.Info("shutting down...")
		book.Save()
		err := sw.Stop()
		if err != nil {
			panic(err)
		}
	})

	err = sw.Start()
	if err != nil {
		panic(err)
	}

	sw.Wait()
}
