package seednode

import (
	"bytes"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/config"
	tmos "github.com/tendermint/tendermint/libs/os"
	"github.com/tendermint/tendermint/p2p"

	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// TSConfig extends tendermint P2PConfig with the things we need
type TSConfig struct {
	config.P2PConfig `mapstructure:",squash"`
	HttpPort         string `mapstructure:"http_port"`
	ChainId          string `mapstructure:"chain_id"`
	LogLevel         string `mapstructure:"log_level"`
}

var configTemplate *template.Template

func init() {
	var err error
	tmpl := template.New("configFileTemplate").Funcs(template.FuncMap{
		"StringsJoin": strings.Join,
	})
	if configTemplate, err = tmpl.Parse(defaultConfigTemplate); err != nil {
		panic(err)
	}
}

func InitConfig() (TSConfig, p2p.NodeKey) {
  var tsConfig = new(TSConfig)
	userHomeDir, err := homedir.Dir()
  if err != nil {
		panic(err)
	}

  // init config directory & files if they don't exists yet
  homeDir := filepath.Join(userHomeDir, ".terranseed", "config")
  nodeKeyFilePath := filepath.Join(homeDir, "node_key.json")
  configFilePath := filepath.Join(homeDir, "config.toml")

  if err = os.MkdirAll(homeDir, os.ModePerm); err != nil {
    panic(err)
  }

	viper.SetConfigName("config")
	viper.AddConfigPath(filepath.Join(userHomeDir, ".terranseed", "config"))

	if err := viper.ReadInConfig(); err == nil {
		logger.Info(fmt.Sprintf("Using config file: %s", viper.ConfigFileUsed()))
		err := viper.Unmarshal(tsConfig)
		if err != nil {
			panic("Invalid config file!")
		}
	} else if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		// ignore not found error, return other errors
		logger.Info("No existing configuration found, generating one")
    tsConfig = initDefaultConfig()
    writeConfigFile(configFilePath, tsConfig)
  } else {
    panic(err)
  }

	nodeKey, err := p2p.LoadOrGenNodeKey(nodeKeyFilePath)
	if err != nil {
		panic(err)
	}

	logger.Info("Configuration",
		"path", configFilePath,
		"node listen", tsConfig.ListenAddress,
		"http server port", tsConfig.HttpPort,
		"chain", tsConfig.ChainId,
	)
  if tsConfig.Seeds == "" || tsConfig.ChainId == "" {
    panic("Don't forget to fill ChainId and Seeds in config.toml file and personalize settings")
  }
	return *tsConfig, *nodeKey
}

func initDefaultConfig() *TSConfig {
	p2PConfig := config.DefaultP2PConfig()
  tsConfig := &TSConfig{
    P2PConfig: *p2PConfig,
    HttpPort:  "8090",
    LogLevel:  "info",
  }
  tsConfig.AllowDuplicateIP = true
  tsConfig.MaxNumInboundPeers = 3000
  tsConfig.MaxNumOutboundPeers = 3000
  return tsConfig
}

// WriteConfigFile renders config using the template and writes it to configFilePath.
func writeConfigFile(configFilePath string, config *TSConfig) {
	var buffer bytes.Buffer

	if err := configTemplate.Execute(&buffer, config); err != nil {
		panic(err)
	}

	tmos.MustWriteFile(configFilePath, buffer.Bytes(), 0644)
}

const defaultConfigTemplate = `# This is a TOML config file.
# For more information, see https://github.com/toml-lang/toml

# NOTE: Any path below can be absolute (e.g. "/var/myawesomeapp/data") or
# relative to the home directory (e.g. "data"). The home directory is
# "$HOME/.tendermint" by default, but could be changed via $TMHOME env variable
# or --home cmd flag.

#######################################################
###    TerranSeed Server Configuration Options      ###
#######################################################
[chain]
# the ChainId of the network to crawl
chain_id = "{{ .ChainId }}"

# Comma separated list of seed nodes to connect to
seeds = "{{ .Seeds }}"

[web]
http_port = "{{ .HttpPort }}"

[p2p]
# Output level for logging: "info" or "debug". debug will enable pex and addrbook verbose logs
log_level = "{{ .LogLevel }}"

# TCP or UNIX socket address for the RPC server to listen on
laddr = "{{ .ListenAddress }}"

# Address to advertise to peers for them to dial
# If empty, will use the same port as the laddr,
# and will introspect on the listener or use UPnP
# to figure out the address. ip and port are required
# example: 159.89.10.97:26656
external_address = "{{ .ExternalAddress }}"

# Must be true, otherwise Tendermint's p2p library will
# deny making connections to peers with the same IP address.
addr_book_strict = {{ .AddrBookStrict }}

# Maximum number of inbound peers. This value can be huge as we don't keep connections opened
max_num_inbound_peers = {{ .MaxNumInboundPeers }}

# Maximum number of outbound peers to connect to, excluding persistent peers
max_num_outbound_peers = {{ .MaxNumOutboundPeers }}

# Maximum pause when redialing a persistent peer (if zero, exponential backoff is used)
persistent_peers_max_dial_period = "{{ .PersistentPeersMaxDialPeriod }}"

# Rate at which packets can be sent, in bytes/second
send_rate = {{ .SendRate }}

# Rate at which packets can be received, in bytes/second
recv_rate = {{ .RecvRate }}

# Toggle to disable guard against peers connecting from the same ip. Similar to addr_book_strict
allow_duplicate_ip = {{ .AllowDuplicateIP }}
`
