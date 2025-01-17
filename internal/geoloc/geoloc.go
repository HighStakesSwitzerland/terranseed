package geoloc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/transport/http/jsonrpc"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/terran-stakers/terranseed/internal/seednode"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	ResolvedPeers []GeolocalizedPeers
	logger        = log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "geoloc")
	ipApiUrl      = "http://ip-api.com/batch"
)

type GeolocalizedPeers struct {
	seednode.Peer
	Country string  `json:"country"`
	Region  string  `json:"region"`
	City    string  `json:"city"`
	Lat     float32 `json:"lat"`
	Lon     float32 `json:"lon"`
	Isp     string  `json:"isp"`
	Org     string  `json:"org"`
	As      string  `json:"as"`
	NodeId  string  `json:"nodeId"`
}

type ipServiceResponse struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"country_code"`
	Region      string  `json:"region"`
	RegionName  string  `json:"region_name"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float32 `json:"lat"`
	Lon         float32 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
	Query       string  `json:"Query"`
}

/*
	Resolve ips using https://ip-api.com/ geolocation free service
	Appends the new resolved peers to the ResolvedPeers slice, so we keep the full list since the startup
*/
func ResolveIps(peerList []*seednode.Peer) {
  if len(peerList) > 0 {
    ResolvedPeers = append(ResolvedPeers, resolve(getUnresolvedPeers(peerList))...)
    logger.Info(fmt.Sprintf("We have %d total resolved peers", len(ResolvedPeers)))
  }
}

func resolve(peers []*seednode.Peer) []GeolocalizedPeers {
	chunkSize := 10
	var geolocalizedPeers []GeolocalizedPeers
	unresolvedPeers := getUnresolvedPeers(peers)
	peersLength := len(unresolvedPeers)
  if (peersLength > 0) {
    logger.Info(fmt.Sprintf("There is %d new peers that need resolution", peersLength))
  }

	for i := 0; i < peersLength; i += chunkSize {
		end := i + chunkSize
		if end > peersLength {
			end = peersLength
		}
		var chunk []*seednode.Peer
		chunk = append(chunk, unresolvedPeers[i:end]...)
		if len(chunk) > 0 {
			ipServiceResponses := fillGeolocData(chunk)
			var peer *seednode.Peer
			var newGeolocalizedPeer GeolocalizedPeers
			for _, elt := range ipServiceResponses {
				peer = findPeerInList(elt, unresolvedPeers)
				if peer == nil {
					logger.Error("Could not find peer in existing list! It may have not been resolved by the service")
					continue
				}
				newGeolocalizedPeer = GeolocalizedPeers{
					Peer:    *peer,
					Country: elt.Country,
					Region:  elt.Region,
					City:    elt.City,
					Lat:     elt.Lat,
					Lon:     elt.Lon,
					Isp:     elt.Isp,
					Org:     elt.Org,
					As:      elt.As,
					NodeId:  string(peer.NodeId),
				}
				geolocalizedPeers = append(geolocalizedPeers, newGeolocalizedPeer)
			}
		}
	}
	return geolocalizedPeers
}

func fillGeolocData(chunk []*seednode.Peer) []ipServiceResponse {
	logger.Info(fmt.Sprintf("Calling ip-api service with %d IPs", len(chunk)))
	var list []string

	for _, peer := range chunk {
		list = append(list, peer.IP)
	}

	payload, err := json.Marshal(list)
	if err != nil {
		logger.Error("Failed to marshal peers list for geoloc service", err)
		return nil
	}

	post, err := http.Post(ipApiUrl, jsonrpc.ContentType, bytes.NewBuffer(payload))
	if err != nil {
		logger.Error("IP geoloc service returned an error", err)
		return nil
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error("Error while waiting for response", err)
		}
	}(post.Body)

	body, err := ioutil.ReadAll(post.Body)
	if err != nil {
		logger.Error("Error reading reponse", err)
	}

	//Decode the data
	response := make([]ipServiceResponse, 0)
	if err := json.Unmarshal(body, &response); err != nil {
		logger.Error("Error while unmarshalling response", err)
	}

	return response
}

func getUnresolvedPeers(peers []*seednode.Peer) []*seednode.Peer {
	var peersToResolve []*seednode.Peer

	for _, peer := range peers {
		if !isResolved(peer) {
			peersToResolve = append(peersToResolve, peer)
		}
	}
	return peersToResolve
}

func isResolved(peer *seednode.Peer) bool {
	if ResolvedPeers == nil {
		return false
	}
	for _, elt := range ResolvedPeers {
		if elt.Peer.IP == peer.IP {
			return true
		}
	}
	return false
}

func findPeerInList(ipServiceResponse ipServiceResponse, peer []*seednode.Peer) *seednode.Peer {
	for _, elt := range peer {
		if elt.IP == ipServiceResponse.Query {
			return elt
		}
	}
	return nil
}
