package tendermint

import (
	"encoding/json"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/p2p/pex"
	"net/http"
	"time"
)

var (
	peerList []*peer
)

type peer struct {
	Moniker string `json:"Moniker"`
	IP     string `json:"ip"`
	Geoloc string `json:"Geoloc"`
}

type peers struct {
	Peers []*peer `json:"Peers"`
}

func InitPeers(sw *p2p.Switch) {
	go func() {
		// Fire periodically
		ticker := time.NewTicker(5 * time.Second)

		for {
			select {
			case <-ticker.C:
				logger.Info("Peers list", "peers", sw.Peers().List())
				peerList = p2pPeersToPeerList(sw.Peers().List())
			}
		}
	}()
}

func p2pPeersToPeerList(list []p2p.Peer) []*peer {
	var _peers []*peer
	for _, p := range list {
		_peers = append(_peers, &peer{
			Moniker: p.NodeInfo().(p2p.DefaultNodeInfo).Moniker,
			IP:      p.(pex.Peer).RemoteIP().String(),
			Geoloc:  ""})
	}
	return _peers
}

func WritePeers(w http.ResponseWriter, r *http.Request) {
	marshal, err := json.Marshal(peers{Peers: peerList})
	if err != nil {
		logger.Info("Failed to marshal peers list")
		return
	}
	_, err = w.Write(marshal)
	if err != nil {
		return
	}
}
