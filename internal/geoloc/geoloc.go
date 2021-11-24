package geoloc

import (
	"github.com/tendermint/tendermint/libs/log"
	"os"
)

var (
	logger    = log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "geoloc")
)
