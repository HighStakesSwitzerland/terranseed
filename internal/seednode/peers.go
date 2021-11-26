package seednode

import (
	"fmt"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/p2p/pex"
)

var (
	peerList []*Peer
)

type Peer struct {
	Moniker string `json:"moniker"`
	IP      string `json:"ip"`
}

/*
	Returns the current reactor peers. As in seed mode the pex module disconnects quickly, this list can grow and shrink
	according to the current connexions
*/
func GetPeers(sw *p2p.Switch) []*Peer {
	logger.Info(fmt.Sprintf("Address book contains %d new peers", len(sw.Peers().List())), "peers", sw.Peers().List())
	peerList = p2pPeersToPeerList(sw.Peers().List())
	return peerList
}

func p2pPeersToPeerList(list []p2p.Peer) []*Peer {
	var _peers []*Peer
	for _, p := range list {
		_peers = append(_peers, &Peer{
			Moniker: p.NodeInfo().(p2p.DefaultNodeInfo).Moniker,
			IP:      p.(pex.Peer).RemoteIP().String(),
		})
	}
	return _peers
}
