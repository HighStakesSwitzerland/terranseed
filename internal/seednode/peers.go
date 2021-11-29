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
	IP      string `json:"-"` // IPs should not be sent to the frontend
	NodeId  p2p.ID `json:"-"`
}

/*
	Returns the current reactor peers. As in seed mode the pex module disconnects quickly, this list can grow and shrink
	according to the current connexions
*/
func GetPeers(peers []p2p.Peer) []*Peer {
  if len(peers) > 0 {
    logger.Info(fmt.Sprintf("Address book contains %d new peers", len(peers)), "peers", peers)
    peerList = p2pPeersToPeerList(peers)
    return peerList
  }
  return nil
}

func p2pPeersToPeerList(list []p2p.Peer) []*Peer {
	var _peers []*Peer
	for _, p := range list {
		_peers = append(_peers, &Peer{
			Moniker: p.NodeInfo().(p2p.DefaultNodeInfo).Moniker,
			IP:      p.(pex.Peer).RemoteIP().String(),
			NodeId:  p.NodeInfo().ID(),
		})
	}
	return _peers
}
