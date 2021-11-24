package gmaps

import (
	"github.com/tendermint/tendermint/libs/log"
	"os"
)

var (
	logger    = log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "maps")
)

//https://socketloop.com/tutorials/golang-find-location-by-ip-address-and-display-with-google-map
