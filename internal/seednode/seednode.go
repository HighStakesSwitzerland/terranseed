package seednode

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	tmstrings "github.com/tendermint/tendermint/libs/strings"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/p2p/pex"
	"github.com/tendermint/tendermint/version"
	"os"
	"path/filepath"
	"time"
)

var (
	logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "config")
)

func StartSeedNode(seedConfig TSConfig, nodeKey p2p.NodeKey) *p2p.Switch {
	cfg := config.DefaultP2PConfig()
	cfg.AllowDuplicateIP = true

	protocolVersion :=
		p2p.NewProtocolVersion(
			version.P2PProtocol,
			version.BlockProtocol,
			0,
		)

	// NodeInfo gets info on your node
	nodeInfo := p2p.DefaultNodeInfo{
		ProtocolVersion: protocolVersion,
		DefaultNodeID:   nodeKey.ID(),
		ListenAddr:      seedConfig.ListenAddress,
		Network:         seedConfig.ChainId,
		Version:         "1.0.0",
		Channels:        []byte{pex.PexChannel},
		Moniker:         fmt.Sprintf("%s-seed", seedConfig.ChainId),
	}

	addr, err := p2p.NewNetAddressString(p2p.IDAddressString(nodeInfo.DefaultNodeID, nodeInfo.ListenAddr))
	if err != nil {
		panic(err)
	}

	transport := p2p.NewMultiplexTransport(nodeInfo, nodeKey, p2p.MConnConfig(cfg))
	if err := transport.Listen(*addr); err != nil {
		panic(err)
	}

	userHomeDir, _ := homedir.Dir()
	addrBookFilePath := filepath.Join(userHomeDir, ".terranseed", "config", "addrbook.json")
	addrBook := pex.NewAddrBook(addrBookFilePath, seedConfig.AddrBookStrict)
	addrBook.SetLogger(logger.With("module", "addrbook"))

	pexReactor := pex.NewReactor(addrBook, &pex.ReactorConfig{
		SeedMode:                     true,
		Seeds:                        tmstrings.SplitAndTrim(seedConfig.Seeds, ",", " "),
		SeedDisconnectWaitPeriod:     1 * time.Second, // default is 28 hours, we just want to harvest as many addresses as possible
		PersistentPeersMaxDialPeriod: 5 * time.Minute,               // use exponential back-off
	})

	sw := p2p.NewSwitch(cfg, transport)
	sw.SetNodeKey(&nodeKey)
	sw.SetAddrBook(addrBook)
	sw.AddReactor("pex", pexReactor)

	// Set loggers. Uncomment to enable
	if seedConfig.LogLevel == "debug" {
    // Switch module logs a lot, and it is not very useful, set as debug
    sw.SetLogger(logger.With("module", "switch"))
    // Same for pex module
    pexReactor.SetLogger(logger.With("module", "pex"))
  }

	// last
	sw.SetNodeInfo(nodeInfo)

	tmos.TrapSignal(logger, func() {
		logger.Info("shutting down...")
		addrBook.Save()
		err := sw.Stop()
		if err != nil {
			panic(err)
		}
	})

	err = sw.Start()
	if err != nil {
		panic(err)
	}

	return sw

}
