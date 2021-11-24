package tendermint

import (
	"github.com/mitchellh/go-homedir"
	"github.com/tendermint/tendermint/p2p"
	"os"
	"path/filepath"
)

var (
	ConfigDir = ".tinyseed"
)

// Config defines the configuration format
type Config struct {
	ListenAddress       string `toml:"laddr" comment:"Address to listen for incoming connections"`
	HttpPort            string `toml:"http_port" comment:"Port for the http server"`
	ChainID             string `toml:"chain_id" comment:"network identifier (todo move to cli flag argument? keeps the config network agnostic)"`
	NodeKeyFile         string `toml:"node_key_file" comment:"path to node_key (relative to tendermint-seed home directory or an absolute path)"`
	AddrBookFile        string `toml:"addr_book_file" comment:"path to address book (relative to tendermint-seed home directory or an absolute path)"`
	AddrBookStrict      bool   `toml:"addr_book_strict" comment:"Set true for strict routability rules\n Set false for private or local networks"`
	MaxNumInboundPeers  int    `toml:"max_num_inbound_peers" comment:"maximum number of inbound connections"`
	MaxNumOutboundPeers int    `toml:"max_num_outbound_peers" comment:"maximum number of outbound connections"`
	Seeds   			string `toml:"seeds" comment:"seed nodes we can use to discover peers"`
	NodeKey 			*p2p.NodeKey
}

func InitConfig(seedConfig *Config) {
	userHomeDir, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	// init config directory & files
	homeDir := filepath.Join(userHomeDir, ConfigDir, "config")
	configFilePath := filepath.Join(homeDir, "config.toml")
	nodeKeyFilePath := filepath.Join(homeDir, seedConfig.NodeKeyFile)
	addrBookFilePath := filepath.Join(homeDir, seedConfig.AddrBookFile)

	mkdirAllPanic(filepath.Dir(nodeKeyFilePath), os.ModePerm)
	mkdirAllPanic(filepath.Dir(addrBookFilePath), os.ModePerm)
	mkdirAllPanic(filepath.Dir(configFilePath), os.ModePerm)

	nodeKey, err := p2p.LoadOrGenNodeKey(nodeKeyFilePath)
	if err != nil {
		panic(err)
	}
	seedConfig.NodeKey = nodeKey


	logger.Info("Configuration",
		"key", seedConfig.NodeKey.ID(),
		"node listen", seedConfig.ListenAddress,
		"http server port", seedConfig.HttpPort,
		"chain", seedConfig.ChainID,
		"strict-routing", seedConfig.AddrBookStrict,
		"max-inbound", seedConfig.MaxNumInboundPeers,
		"max-outbound", seedConfig.MaxNumOutboundPeers,
	)

}

// mkdirAllPanic invokes os.MkdirAll but panics if there is an error
func mkdirAllPanic(path string, perm os.FileMode) {
	err := os.MkdirAll(path, perm)
	if err != nil {
		panic(err)
	}
}
